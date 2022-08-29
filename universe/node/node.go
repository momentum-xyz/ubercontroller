package node

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	ctx        context.Context
	db         database.DB
	log        *zap.SugaredLogger
	worlds     universe.Worlds
	assets2d   universe.Assets2d
	assets3d   universe.Assets3d
	spaceTypes universe.SpaceTypes
	mu         sync.RWMutex
	id         uuid.UUID
}

func NewNode(
	id uuid.UUID,
	db database.DB,
	worlds universe.Worlds,
	assets2D universe.Assets2d,
	assets3D universe.Assets3d,
	spaceTypes universe.SpaceTypes,
) *Node {
	return &Node{
		id:         id,
		db:         db,
		worlds:     worlds,
		assets2d:   assets2D,
		assets3d:   assets3D,
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

func (n *Node) GetWorlds() universe.Worlds {
	return n.worlds
}

func (n *Node) GetAssets2d() universe.Assets2d {
	return n.assets2d
}

func (n *Node) GetAssets3d() universe.Assets3d {
	return n.assets3d
}

func (n *Node) GetSpaceTypes() universe.SpaceTypes {
	return n.spaceTypes
}

func (n *Node) Run() error {
	return errors.Errorf("implement me")
}

func (n *Node) Stop() error {
	return errors.Errorf("implement me")
}

func (n *Node) Load(ctx context.Context) error {
	if err := n.LoadGlobalData(ctx); err != nil {
		return errors.WithMessage(err, "failed to load global data")
	}

	_, err := n.db.WorldsGetWorldIDs(ctx)
	if err != nil {
		return errors.WithMessage(err, "failed to WorldsGetWorldIDs")
	}

	return nil
}

func (n *Node) LoadGlobalData(ctx context.Context) error {
	group, ctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return n.assets2d.Load(ctx)
	})
	group.Go(func() error {
		return n.assets3d.Load(ctx)
	})
	group.Go(func() error {
		return n.spaceTypes.Load(ctx)
	})

	return group.Wait()
}
