package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

type HandShake struct {
	HandshakeVersion int       `json:"handshake_version"`
	ProtocolVersion  int       `json:"protocol_version"`
	Token            string    `json:"token"`
	UserId           umid.UMID `json:"user_id"`
	SessionId        umid.UMID `json:"session_id"`
}

func init() {
	registerMessage(HandShake{})
}

func (g *HandShake) GetType() MsgType {
	return 0x7C41941A
}
