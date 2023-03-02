package posbus

import "github.com/google/uuid"

type HandShake struct {
	HandshakeVersion int
	ProtocolVersion  int
	Token            string
	UserId           uuid.UUID
	SessionId        uuid.UUID
	Url              string
	WorldId          uuid.UUID
}
