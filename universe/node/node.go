package node

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
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
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	n.ctx = ctx
	n.log = log

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

func (n *Node) Run(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)
	worlds := n.GetWorlds().GetWorlds()

	worlds.Mu.RLock()
	for _, world := range worlds.Data {
		world := world

		group.Go(func() error {
			if err := world.Run(gctx); err != nil {
				return errors.WithMessagef(err, "failed to run world: %s", world.GetID())
			}
			return nil
		})
	}
	worlds.Mu.RUnlock()

	return group.Wait()
}

func (n *Node) Stop() error {
	var wg sync.WaitGroup
	var errs *multierror.Error
	var errsMu sync.Mutex
	worlds := n.GetWorlds().GetWorlds()

	worlds.Mu.RLock()
	for _, world := range worlds.Data {
		wg.Add(1)

		go func(world universe.World) {
			defer wg.Done()

			if err := world.Stop(); err != nil {
				errsMu.Lock()
				defer errsMu.Unlock()

				errs = multierror.Append(errs, errors.WithMessagef(err, "failed to stop world: %s", world.GetID()))
			}
		}(world)
	}
	worlds.Mu.RUnlock()

	return errs.ErrorOrNil()
}

func (n *Node) Load(ctx context.Context) error {
	group, gctx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return n.assets2d.Load(gctx)
	})
	group.Go(func() error {
		return n.assets3d.Load(gctx)
	})
	group.Go(func() error {
		return n.spaceTypes.Load(gctx)
	})

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load node data")
	}

	return n.worlds.Load(ctx)
}
