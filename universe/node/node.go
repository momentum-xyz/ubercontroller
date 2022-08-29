package node

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/controller/types"
	"github.com/momentum-xyz/controller/universe"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	worlds     universe.Worlds
	assets2D   universe.Assets2D
	assets3D   universe.Assets3D
	spaceTypes universe.SpaceTypes
	mu         sync.RWMutex
	id         uuid.UUID
}

func NewNode(
	id uuid.UUID,
	worlds universe.Worlds,
	assets2D universe.Assets2D,
	assets3D universe.Assets3D,
	spaceTypes universe.SpaceTypes,
) *Node {
	return &Node{
		id:         id,
		worlds:     worlds,
		assets2D:   assets2D,
		assets3D:   assets3D,
		spaceTypes: spaceTypes,
	}
}

func (n *Node) GetID() uuid.UUID {
	n.mu.RLock()
	defer n.mu.RUnlock()

	return n.id
}

func (n *Node) Initialize(ctx context.Context) error {
	log, ok := ctx.Value(types.ContextLoggerKey).(*zap.SugaredLogger)
	if !ok {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	n.log = log
	n.ctx = ctx

	return nil
}

func (n *Node) Run() error {
	return errors.Errorf("implement me")
}

func (n *Node) Stop() error {
	return errors.Errorf("implement me")
}

func (n *Node) Load() error {
	return errors.Errorf("implement me")
}
