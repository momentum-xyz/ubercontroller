package posbus

import "github.com/momentum-xyz/ubercontroller/utils/umid"

// A HandShake is the first message a client sends after connecting.
type HandShake struct {
	// Versioning for this message, for compatibility handling.
	HandshakeVersion int `json:"handshake_version"`

	// Versioning for the protocol to use after the handshake.
	ProtocolVersion int `json:"protocol_version"`

	// Authentication token (JWT).
	Token string `json:"token"`

	// User identifier (should match the token).
	UserId umid.UMID `json:"user_id"`

	// Unique session identifier, for state/reconnection handling.
	SessionId umid.UMID `json:"session_id"`
}

func init() {
	registerMessage(HandShake{})
}

func (g *HandShake) GetType() MsgType {
	return 0x7C41941A
}
