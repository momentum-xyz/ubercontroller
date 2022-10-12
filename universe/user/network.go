package user

import (
	"container/list"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/pkg/errors"
	"time"
)

const (
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second
	// send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	inMessageSizeLimit = 1024
	// maximal size of buffer in messages, after which we drop connection as not-working
	maxBufferSize = 10000
	// Negative Number to indicate closed chan, large enough to be less than any number of outstanding
	chanIsClosed = -0x3FFFFFFFFFFFFFFF
)

func (u *User) Close() error {
	return nil
}

func (u *User) StartIOPumps() {
	u.lastPositionUpdateTimestamp = int64(0)

	go u.writePump()
	go u.writePump()
}

func (u *User) readPump() {
	u.conn.SetReadLimit(inMessageSizeLimit)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		messageType, message, err := u.conn.ReadMessage()
		if err != nil {
			closedByClient := false
			if ce, ok := err.(*websocket.CloseError); ok {
				switch ce.Code {
				case websocket.CloseNormalClosure,
					websocket.CloseGoingAway,
					websocket.CloseNoStatusReceived:
					closedByClient = true

				}
			}
			if closedByClient {
				u.log.Info(errors.WithMessagef(err, "Connection: read pump: websocket closed by client"))
			} else {
				u.log.Debug(errors.WithMessage(err, "Connection: read pump: failed to read message from connection"))
			}
			break
		}
		if messageType != websocket.BinaryMessage {
			u.log.Errorf("Connection: read pump: wrong incoming message type: %d", messageType)
		} else {
			u.OnMessage(posbus.MsgFromBytes(message))
		}
	}
	// this close will trigger write defer function anyway
	u.conn.Close()
	u.log.Info("Connection: end of read pump")
}

func (u *User) initiateShutDown(needToRemoveFromWorld bool) {
	//drain send channel
	ns := u.numSendsQueued.Swap(chanIsClosed)
	for i := int64(0); i < ns; i++ {
		<-u.send
	}
	close(u.send)
	u.conn.Close()

	// then remove from world is necessary
	if needToRemoveFromWorld {
		u.world.RemoveUser(u, true)
	}
	return
}

func (u *User) writePump() {
	needToRemoveFromWorld := true
	defer func() {
		u.initiateShutDown(needToRemoveFromWorld)
	}()
	ticker := time.NewTicker(pingPeriod)
	pingMessage, _ := websocket.NewPreparedMessage(websocket.PingMessage, nil)
	buffer := list.New()
	for {
		select {
		case message := <-u.send:
			// we took message from queue
			u.numSendsQueued.Add(-1)

			// sending nil to send chan will stop this evil loop
			if message == nil {
				// if this break was initiated via user.Shutdown() we don't remove from world, assuming something else is taking care of that
				needToRemoveFromWorld = false
				return
			}

			if u.bufferSends.Load() == true {
				//  if we should buffer messages instead of sending
				buffer.PushBack(message)
			} else {
				//  drain buffer and send current message
				if buffer.Len() > 0 {
					var next *list.Element
					for e := buffer.Front(); e != nil; e = next {
						next = e.Next()
						m := buffer.Remove(e).(*websocket.PreparedMessage)
						if u.SendDirectly(m) != nil {
							return
						}
					}
				}
				if u.SendDirectly(message) != nil {
					return
				}
			}

		case <-ticker.C:
			if u.SendDirectly(pingMessage) != nil {
				return
			}
		}
	}
}

func (u *User) Shutdown() {
	ns := u.numSendsQueued.Add(1)
	if ns >= 0 {
		u.send <- nil
	}
}

func (u *User) SendDirectly(message *websocket.PreparedMessage) error {
	// not concurrent, to be used in single particular location
	u.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return u.conn.WritePreparedMessage(message)
}

func (u *User) Send(m *websocket.PreparedMessage) {
	if m == nil {
		return
	}
	// ns acts simultaneously as number of clients in send process and as blocker if negative
	// we increment, and we decrement when leave this method
	ns := u.numSendsQueued.Add(1)
	if ns >= 0 {
		u.send <- m
	}
}

func (u *User) SetConnection(id uuid.UUID, socketConnection *websocket.Conn) error {
	u.sessionID = id
	u.conn = socketConnection
	u.send = make(chan *websocket.PreparedMessage, 10)
	u.bufferSends.Store(true)
	return nil
}

func (u *User) GetSessionID() uuid.UUID {
	return u.sessionID
}
