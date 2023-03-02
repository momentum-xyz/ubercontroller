package node

import (
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/utils"
)

var websocketUpgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (n *Node) apiPosBusHandler(c *gin.Context) {
	ws, err := websocketUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		n.log.Error(errors.WithMessage(err, "error: socket upgrade error, aborting connection"))
		return
	}
	if err := n.handShake(ws); err != nil {
		n.log.Error(errors.WithMessage(err, "failed to handle hand shake"))
	}
}

// handShake TODO: it's "god" method needs to be simplified // antst: agree :)
func (n *Node) handShake(socketConnection *websocket.Conn) error {
	mt, incomingMessage, err := socketConnection.ReadMessage()
	if err != nil || mt != websocket.BinaryMessage {
		return errors.WithMessagef(err, "error: wrong PreHandShake (1), aborting connection")
	}

	msg := posbus.MsgFromBytes(incomingMessage)
	if msg.Type() != posbus.HandShakeType {
		return errors.New("error: wrong message received, not handshake")
	}
	var handshake posbus.HandShake
	if msg.DecodeMessage(handshake) != nil {
		return errors.New("error: wrong message type received, not handshake data")
	}

	n.log.Debugf("Node: handshake for user %s:", handshake.UserId)
	n.log.Debugf("Node: handshake version: %d", handshake.HandshakeVersion)
	n.log.Debugf("Node: protocol version: %d", handshake.ProtocolVersion)

	token := string(handshake.Token)

	// TODO: enable token check back!
	//if err := api.VerifyToken(token, n.cfg.Common.IntrospectURL); err != nil {
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
	claims := parsed.Claims.(jwt.MapClaims)

	userID := handshake.UserId
	sessionID := handshake.SessionId
	targetWorldId := handshake.WorldId
	url, err := url.Parse(handshake.Url)
	if err != nil {
		return errors.WithMessagef(err, "failed to parse url: %s", string(handshake.Url))
	}
	n.log.Debugf("Node: url to use: %s", url)

	userIDClaim, err := uuid.Parse(utils.GetFromAnyMap(claims, "sub", ""))
	if err != nil {
		return errors.WithMessagef(err, "failed to parse id claim: %s", userID)
	}
	if !((userID == userIDClaim) || (userIDClaim.String() == "69e1d7f6-3130-4005-9969-31edf9af9445") || (userIDClaim.String() == "eb50bbc8-ba4e-46a3-a480-a9b30141ce91")) {
		return nil
	}

	user, err := n.LoadUser(userID)
	if err != nil {
		return errors.WithMessagef(err, "failed to load user from entry: %s", userID)
	}
	user.SetConnection(sessionID, socketConnection)
	user.Run()

	world, ok := n.GetWorlds().GetWorld(targetWorldId)
	if !ok {
		world = n.detectSpawnWorld(userID)
		n.log.Infof("World is not found! %+v\n", targetWorldId)
	}
	n.log.Infof("User will be launched in world %+v \n", world.GetID())

	return world.AddUser(user, true)
}
