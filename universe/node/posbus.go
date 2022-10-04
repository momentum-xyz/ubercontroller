package node

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	flatbuffers "github.com/google/flatbuffers/go"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/posbus-protocol/flatbuff/go/api"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"net/http"
	"net/url"
)

var WebsocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (n *Node) PosBusConnectionHandler(ctx *gin.Context) {
	ws, err := WebsocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "error: socket upgrade error, aborting connection"))
		return
	}

	n.HandShake(ws)
}

// HandShake TODO: it's "god" method needs to be simplified // antst: agree :)
func (n *Node) HandShake(socketConnection *websocket.Conn) {

	mt, incomingMessage, err := socketConnection.ReadMessage()
	if err != nil || mt != websocket.BinaryMessage {
		n.log.Error(errors.WithMessagef(err, "error: wrong PreHandShake (1), aborting connection"))
		return
	}

	msg := posbus.MsgFromBytes(incomingMessage)
	if msg.Type() != posbus.MsgTypeFlatBufferMessage {
		n.log.Error("error: wrong message received, not Handshake.")
		return
	}
	msgObj := posbus.MsgFromBytes(incomingMessage).AsFlatBufferMessage()
	msgType := msgObj.MsgType()
	if msgType != api.MsgHandshake {
		n.log.Error("error: wrong message type received, not Handshake.")
		return
	}

	var handshake *api.Handshake
	unionTable := &flatbuffers.Table{}
	if msgObj.Msg(unionTable) {
		handshake = &api.Handshake{}
		handshake.Init(unionTable.Bytes, unionTable.Pos)
	}

	n.log.Info("handshake for user:", message.DeserializeGUID(handshake.UserId(nil)))
	n.log.Info("handshake version:", handshake.HandshakeVersion())
	n.log.Info("protocol version:", handshake.ProtocolVersion())

	token := string(handshake.UserToken())

	// TODO: enable token check back!
	//if err := auth.VerifyToken(token, n.cfg.Common.IntrospectURL); err != nil {
	//	userID := message.DeserializeGUID(handshake.UserId(nil))
	//	n.log.Errorf("error: wrong PreHandShake (invalid token: %s), aborting connection: %s", userID, err)
	//	socketConnection.SetWriteDeadline(time.Now().Add(10 * time.Second))
	//	socketConnection.WritePreparedMessage(posbus.NewSignalMsg(posbus.SignalInvalidToken).WebsocketMessage())
	//	return nil, false
	//}

	parsed, _ := jwt.Parse(
		token, func(token *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		},
	)

	userID := message.DeserializeGUID(handshake.UserId(nil))
	sessionID := message.DeserializeGUID(handshake.SessionId(nil))
	URL, _ := url.Parse(string(handshake.Url()))
	n.log.Info("URL to use:", URL)

	claims := parsed.Claims.(jwt.MapClaims)
	userIDclaim, _ := uuid.Parse(utils.GetFromAnyMap(claims, "sub", ""))

	if !((userID == userIDclaim) || (userIDclaim.String() == "69e1d7f6-3130-4005-9969-31edf9af9445") || (userIDclaim.String() == "eb50bbc8-ba4e-46a3-a480-a9b30141ce91")) {
		return
	}

	user, err := n.LoadUser(userID)
	if err != nil {
		return
	}
	user.SetConnection(sessionID, socketConnection)

	n.DetectSpawnWorld(userID).AddUser(user, false)
}
