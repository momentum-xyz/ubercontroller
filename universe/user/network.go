package user

import (
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

func (u *User) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		//u.hub.unregister <- u
		u.log.Info("Connection: end of write pump")
	}()

	for {
		select {
		case message, ok := <-u.send:
			if !ok {
				u.log.Debug("Connection: write pump: chan closed")
				return
			}
			if u.SendDirectly(message, ok) != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(u.send)
			for i := 0; i < n; i++ {
				message, ok := <-u.send
				if u.SendDirectly(message, ok) != nil {
					return
				}
			}

		case <-ticker.C:
			u.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := u.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (u *User) SendDirectly(message *websocket.PreparedMessage, ok bool) error {
	u.conn.SetWriteDeadline(time.Now().Add(writeWait))
	if !ok {
		// The hub closed the channel.
		u.conn.WriteMessage(websocket.CloseMessage, []byte{})
		return errors.New("socket closed")
	}
	err := u.conn.WritePreparedMessage(message)
	if err != nil {
		return errors.New("error pushing message")
	}
	return nil
}

func (u *User) Send(m *websocket.PreparedMessage) {
	if m == nil {
		return
	}
	// protect by simple atomic, when acceptance to send chan is closed nCurrentSends is negative
	if u.nCurrentSends.Load() >= 0 {
		u.nCurrentSends.Add(1)
		u.send <- m
		u.nCurrentSends.Add(-1)
	}
}

func (u *User) CloseSendChannel() {
	// cet current number of open sends
	sends := u.nCurrentSends.Swap(chanIsClosed)
	// drain send channel
	ticker := time.NewTicker(time.Microsecond)
	for u.nCurrentSends.Load() > chanIsClosed-sends {
		select {
		case _, _ = <-u.send:

		case <-ticker.C:
		}
	}
	u.log.Info("Channel is drained for user %v", u.id.String())
	//at this point there are no sends are coming anymore GRANTED
	close(u.send)
}
