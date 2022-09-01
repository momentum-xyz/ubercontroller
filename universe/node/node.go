package node

import (
	"context"
	"fmt"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.Node = (*Node)(nil)

type Node struct {
	cfg        *config.Config
	ctx        context.Context
	log        *zap.SugaredLogger
	db         database.DB
	router     *gin.Engine
	worlds     universe.Worlds
	assets2d   universe.Assets2d
	assets3d   universe.Assets3d
	spaceTypes universe.SpaceTypes
	mu         sync.RWMutex
	id         uuid.UUID
}

func NewNode(
	id uuid.UUID,
	cfg *config.Config,
	db database.DB,
	worlds universe.Worlds,
	assets2D universe.Assets2d,
	assets3D universe.Assets3d,
	spaceTypes universe.SpaceTypes,
) *Node {
	return &Node{
		id:         id,
		cfg:        cfg,
		db:         db,
		router:     gin.Default(),
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

func (n *Node) AddAPIRegister(register types.APIRegister) {
	register.RegisterAPI(n.router)
}

func (n *Node) Run() error {
	if err := n.worlds.Run(); err != nil {
		return errors.WithMessage(err, "failed to run worlds")
	}
	return n.router.Run(fmt.Sprintf("%s:%d", n.cfg.Settings.Address, n.cfg.Settings.Port))
}

func (n *Node) Stop() error {
	return n.worlds.Stop()
}

func (n *Node) Load() error {
	group, _ := errgroup.WithContext(n.ctx)

	group.Go(func() error {
		return n.assets2d.Load()
	})
	group.Go(func() error {
		return n.assets3d.Load()
	})
	group.Go(func() error {
		return n.spaceTypes.Load()
	})

	if err := group.Wait(); err != nil {
		return errors.WithMessage(err, "failed to load node data")
	}

	if err := n.worlds.Load(); err != nil {
		return errors.WithMessage(err, "failed to load worlds")
	}

	n.AddAPIRegister(n)

	return nil
}

func (n *Node) Save() error {
	group, _ := errgroup.WithContext(n.ctx)

	group.Go(func() error {
		return n.assets2d.Save()
	})
	group.Go(func() error {
		return n.assets3d.Save()
	})
	group.Go(func() error {
		return n.spaceTypes.Save()
	})
	group.Go(func() error {
		return n.worlds.Save()
	})

	return group.Wait()
}
