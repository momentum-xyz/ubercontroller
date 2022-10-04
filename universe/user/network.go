package user

import (
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
	chanIsClosed = -0x1000000000
)

func (u *User) Close() error {
	return nil
}

func (u *User) StartIOPumps() {
	u.lastPositionUpdateTimestamp = int64(0)

	u.conn.SetReadLimit(inMessageSizeLimit)
	u.conn.SetReadDeadline(time.Now().Add(pongWait))
	u.conn.SetPongHandler(func(string) error { u.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	go u.PingTicker()
}

func (u *User) ReadPump() {
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

func (u *User) PingTicker() {
	ticker := time.NewTicker(pingPeriod)
	pingMessage, _ := websocket.NewPreparedMessage(websocket.PingMessage, nil)
	for {
		select {
		case <-ticker.C:
			//if u.quit.Load() {
			//	return
			//}
			if err := u.SendDirectly(pingMessage); err != nil {
				return
			}
		}
	}
}

func (u *User) CloseHandler() {

}

func (u *User) SendDirectly(message *websocket.PreparedMessage) error {
	u.sendMutex.Lock()
	defer u.sendMutex.Unlock()
	u.conn.SetWriteDeadline(time.Now().Add(writeWait))
	err := u.conn.WritePreparedMessage(message)
	if err != nil {
		return errors.New("error pushing message")
		u.CloseHandler()
	}
	return err
}

func (u *User) Send(m *websocket.PreparedMessage) {
	if m == nil {
		return
	}
	if u.readyToSend.Load() {
		if false {
			// TODO: push content of ring buffer

		}
		err := u.SendDirectly(m)
		if err != nil {
			// TODO: shutdown sequence
		}
	} else {
		// TODO: put to ring buffer
	}
}

func (u *User) SetConnection(id uuid.UUID, socketConnection *websocket.Conn) error {
	u.sessionID = id
	u.conn = socketConnection
	u.send = make(chan *websocket.PreparedMessage, 10)
	u.quit.Store(false)
	u.readyToSend.Store(false)
	return nil
}

func (u *User) GetSessionId() uuid.UUID {
	return u.sessionID
}
