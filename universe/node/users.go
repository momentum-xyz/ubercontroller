package node

import "github.com/pkg/errors"

func (n *Node) UserConnectionsProcessor() {
	n.log.Info("Started UserConnectionsProcessor")
	for {
		n.log.Info("UserConnectionsProcessor waiting for successful handshake")
		handshake := <-n.handshakeChan
		n.log.Info("UserConnectionsProcessor received successful handshake")
		go func() {
			if err := n.LoadUser(handshake); err != nil {
				n.log.Error(errors.WithMessage(err, "UserConnectionsProcessor failed to spawn user"))
			}
		}()
	}
}

func (n *Node) LoadUser(handshake *HandshakeData) error {
	return nil

}
