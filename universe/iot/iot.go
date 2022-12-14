package iot

import (
	"context"
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	inMessageSizeLimit = 160
	// outPosMessageSize = 48
	outPosMessageSize = 28
)

type IOTMessage struct {
	Type string      `json:"type"`
	What string      `json:"what"`
	Data interface{} `json:"data"`
}

type IOTWorker struct {
	ws   *websocket.Conn
	ctx  context.Context
	log  *zap.SugaredLogger
	send chan *websocket.PreparedMessage
}

func NewIOTWorker(ws *websocket.Conn, ctx context.Context) *IOTWorker {

	iw := IOTWorker{ws: ws}
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return nil
	}

	iw.ctx = ctx
	iw.log = log
	iw.send = make(chan *websocket.PreparedMessage, 10)

	return &iw
}

func (iot *IOTWorker) Run() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		iot.log.Infoln("End of IOTWorker")
	}()

	iot.log.Infoln("Start of IOTWorker")
	iot.ws.SetReadLimit(inMessageSizeLimit)
	iot.ws.SetReadDeadline(time.Now().Add(pongWait))
	iot.ws.SetPongHandler(func(string) error { iot.ws.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	go func() {
		for {
			mt, message, err := iot.ws.ReadMessage()
			if err != nil {
				break
			}
			if mt != websocket.BinaryMessage {
				iot.log.Infoln("error: wrong incoming message type")
			} else {
				err := iot.AcceptMessage(message)
				if err != nil {
					iot.log.Error(err)
				}
			}
		}
		iot.ws.Close()
		iot.log.Infoln("End of read")
	}()

	for {
		select {
		case message, ok := <-iot.send:
			if iot.PushMessage(message, ok) != nil {
				return
			}

			// Add queued chat messages to the current websocket message.
			n := len(iot.send)
			for i := 0; i < n; i++ {
				message, ok := <-iot.send
				if iot.PushMessage(message, ok) != nil {
					return
				}
			}

		case <-ticker.C:
			iot.ws.SetWriteDeadline(time.Now().Add(writeWait))
			if err := iot.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}

	}

}

func (iot *IOTWorker) PushMessage(message *websocket.PreparedMessage, ok bool) error {
	iot.ws.SetWriteDeadline(time.Now().Add(writeWait))
	if !ok {
		// The hub closed the channel.
		iot.ws.WriteMessage(websocket.CloseMessage, []byte{})
		return errors.New("socket closed")
	}

	err := iot.ws.WritePreparedMessage((*websocket.PreparedMessage)(message))
	if err != nil {
		return errors.New("error pushing message")
	}
	// if err := w.Close(); err != nil {
	// 	return
	// }
	return nil
}

func (iot *IOTWorker) AcceptMessage(message []byte) error {
	var msg IOTMessage
	err := json.Unmarshal(message, &msg)
	if err != nil {
		return err
	}
	iot.log.Infof("received message %+v\n", msg)
	return nil
}
