package user

import (
	"container/list"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"
	"github.com/zakaria-chahboun/cute"
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

func (u *User) StartIOPumps() {
	u.send = make(chan *websocket.PreparedMessage, 10)
	u.bufferSends.Store(true)
	u.lastPositionUpdateTimestamp = int64(0)
	u.numSendsQueued.Store(0)
	go u.writePump()
	go u.readPump()
}

func (u *User) ReleaseSendBuffer() {
	u.bufferSends.Store(false)
	u.log.Infof("User: ReleaseSendBuffer: messages waterfall opened: %s", u.GetID())
}

func (u *User) LockSendBuffer() {
	u.bufferSends.Store(true)
}

func (u *User) readPump() {
	u.log.Infof("User: start of read pump: %s", u.GetID())

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
				u.log.Info(
					errors.WithMessagef(err, "User: read pump: websocket closed by client: %s", u.GetID()),
				)
			} else {
				u.log.Debug(
					errors.WithMessagef(err, "User: read pump: failed to read message from connection: %s", u.GetID()),
				)
			}
			break
		}
		if messageType != websocket.BinaryMessage {
			u.log.Errorf("User: read pump: wrong incoming message type: %d: %s", messageType, u.GetID())
		} else {
			if err := u.OnMessage(message); err != nil {
				u.log.Warn(errors.WithMessagef(err, "User: read pump: failed to handle message: %s", u.GetID()))
			}
		}
	}
	// this close will cascade writePump defer function anyway on next send
	u.conn.Close()
	u.log.Infof("User: end of read pump: %s", u.GetID())
}

func (u *User) writePump() {
	u.log.Infof("User: start of write pump: %s", u.GetID())

	needToRemoveFromWorld := true
	defer func() {
		u.log.Infof("User: end of write pump: %s", u.GetID())
		if err := u.close(needToRemoveFromWorld); err != nil {
			u.log.Warnf("User: writePump: failed to close user: %s", u.GetID())
		}
	}()

	ticker := time.NewTicker(pingPeriod)
	pingMessage, _ := websocket.NewPreparedMessage(websocket.PingMessage, nil)
	buffer := list.New()
	for {
		select {
		case message := <-u.send:
			//fmt.Println("Send message")
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
			if u.bufferSends.Load() == false {
				if u.SendDirectly(pingMessage) != nil {
					return
				}
			}
		}
	}
}

func (u *User) SendDirectly(message *websocket.PreparedMessage) error {
	// not concurrent, to be used in single particular location
	//u.directLock.Lock()
	//defer u.directLock.Unlock()
	if message == nil {
		cute.SetTitleColor(cute.BrightRed)
		cute.SetMessageColor(cute.Red)
		cute.Printf("User: SendDirectly", "%+v", errors.WithStack(errors.Errorf("empty message received")))
		return nil
	}

	u.conn.SetWriteDeadline(time.Now().Add(writeWait))
	return u.conn.WritePreparedMessage(message)
}

func (u *User) Send(m *websocket.PreparedMessage) error {
	if m == nil {
		cute.SetTitleColor(cute.BrightRed)
		cute.SetMessageColor(cute.Red)
		cute.Printf("User: Send", "%+v", errors.WithStack(errors.Errorf("empty message received")))
		return nil
	}
	// ns acts simultaneously as number of clients in send process and as blocker if negative
	// we increment, and we decrement when leave this method
	ns := u.numSendsQueued.Add(1)
	if ns >= 0 {
		u.send <- m
	}
	return nil
}

func (u *User) SetConnection(id umid.UMID, socketConnection *websocket.Conn) error {
	u.sessionID = id
	u.conn = socketConnection
	return nil
}

func (u *User) GetSessionID() umid.UMID {
	return u.sessionID
}

func (u *User) close(needToRemoveFromWorld bool) error {
	//drain send channel
	ns := u.numSendsQueued.Swap(chanIsClosed)
	for i := int64(0); i < ns; i++ {
		<-u.send
	}

	//close(u.send)
	u.conn.Close()

	// then remove from world is necessary
	if needToRemoveFromWorld {
		world := u.GetWorld()
		if world != nil {
			if _, err := world.RemoveUser(u, true); err != nil {
				return errors.WithMessagef(err, "failed to remove user from world: %s", world.GetID())
			}
		}
	}

	return nil
}
