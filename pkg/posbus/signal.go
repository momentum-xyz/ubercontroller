package posbus

type SignalType uint32

const (
	// SignalNone is not used. If received, the message was not initialized properly.
	SignalNone SignalType = iota

	// A user is only allowed to be in a world 'once' (have a single connection to a world).
	// SignalDualConnection is send went a second (teleport) attempt is made.
	SignalDualConnection

	// The client (3d engine) is ready to receive the world data.
	// @deprecated When the client sends a teleport message, it is assumed to be ready.
	SignalReady

	// The authentication token is invalid.
	// TODO: actually reimplement this? or drop it :)
	SignalInvalidToken

	// Send to client when user is added to the world.
	// @deprecated Client receives [SetWorld] message instead.
	SignalSpawn

	// Send by client when user should be removed from the world.
	// @deprecated Use teleport to go to a different world instead or else just disconnect.
	SignalLeaveWorld

	// @deprecated If connection failed, how would this signal be received then?
	SignalConnectionFailed

	// Send as the first message after a connection, indicate it is now ready to be used.
	// This is actually send by the client library, not the server side application.
	SignalConnected

	// Send when connection is closed. Either by some error or explict action by the client.
	// This is actually send by the client library, not the server side application.
	SignalConnectionClosed

	// Send when user tries to teleport to a world that can't be found.
	SignalWorldDoesNotExist
)

// A Signal is a predefined (small) message to notify the other side of some state or event.
// Used to asynchronously respond to certain other messages or state changes.
type Signal struct {
	// The predefined type of the signal.
	Value SignalType `json:"value"`
}

func init() {
	registerMessage(Signal{})
	addExtraType(SignalType(0))
}

func (g *Signal) GetType() MsgType {
	return 0xADC1964D
}
