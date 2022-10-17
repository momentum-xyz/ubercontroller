package world

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/posbus-protocol/posbus"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"sync/atomic"
	"time"
)

var _ universe.World = (*World)(nil)

type DecorationMetadataNew struct {
	AssetID  uuid.UUID  `json:"asset_id" db:"asset_id" mapstructure:"asset_id"`
	Position cmath.Vec3 `json:"position" db:"position" mapstructure:"position"`
	rotation cmath.Vec3
}

type WorldMeta struct {
	LOD              []uint32                `json:"lod" db:"lod" mapstructure:"lod"`
	Decorations      []DecorationMetadataNew `json:"decorations,omitempty" db:"decorations,omitempty" mapstructure:"decorations"`
	AvatarController uuid.UUID               `json:"avatar_controller" db:"avatar_controller" mapstructure:"avatar_controller"`
	SkyboxController uuid.UUID               `json:"skybox_controller" db:"skybox_controller" mapstructure:"skybox_controller"`
}

type World struct {
	*space.Space
	ctx              context.Context
	log              *zap.SugaredLogger
	db               database.DB
	pluginController *mplugin.PluginController
	//corePluginInstance  mplugin.PluginInstance
	corePluginInterface mplugin.PluginInterface
	broadcast           chan *websocket.PreparedMessage
	metaMsg             atomic.Pointer[websocket.PreparedMessage]
	metaData            WorldMeta
}

func NewWorld(id uuid.UUID, db database.DB) *World {
	world := &World{
		db: db,
	}
	world.Space = space.NewSpace(id, db, world)
	world.pluginController = mplugin.NewPluginController(id)
	//world.corePluginInstance, _ = world.pluginController.AddPlugin(world.GetID(), world.corePluginInitFunc)
	world.pluginController.AddPlugin(universe.GetSystemPluginID(), world.corePluginInitFunc)
	return world
}

func (w *World) corePluginInitFunc(pi mplugin.PluginInterface) (mplugin.PluginInstance, error) {
	instance := CorePluginInstance{PluginInterface: pi}
	w.corePluginInterface = pi
	return instance, nil
}

func (w *World) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.ContextLoggerKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.ContextLoggerKey))
	}

	w.ctx = ctx
	w.log = log

	return w.Space.Initialize(ctx)
}

// TODO: implement
func (w *World) Run() error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		//case message := <-cu.broadcast:
		// v := reflect.ValueOf(cu.broadcast)
		// fmt.Println(color.Red, "Bcast", wc.users.Num(), v.Len(), color.Reset)
		//go cu.PerformBroadcast(message)
		// logger.Logln(4, "BcastE")
		case <-ticker.C:
			// fmt.Println(color.Red, "Ticker", color.Reset)
			go w.broadcastPositions()
		}
	}
	return nil
}

func (w *World) broadcastPositions() {
	flag := false
	w.Users.Mu.RLock()
	numClients := len(w.Users.Data)
	msg := posbus.NewUserPositionsMsg(numClients)
	if numClients > 0 {
		flag = true
		i := 0
		for _, u := range w.Users.Data {
			msg.SetPosition(i, u.GetPosBuffer())
			i++
		}
	}
	w.Users.Mu.RUnlock()
	if flag {
		w.Space.Broadcast(msg.WebsocketMessage(), false)
	}
}

// TODO: implement
func (w *World) Stop() error {
	return nil
}

func (w *World) Load() error {
	w.log.Infof("Loading world: %s", w.GetID())

	entry, err := w.db.SpacesGetSpaceByID(w.ctx, w.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space by id")
	}

	if err := w.LoadFromEntry(entry, true); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}
	w.UpdateWorldMetadata()
	universe.GetNode().AddAPIRegister(w)

	go w.Run()
	//cu.BroadcastPositions()

	w.log.Infof("World loaded: %s", w.GetID())

	return nil
}

func (w *World) UpdateWorldMetadata() error {
	fmt.Printf("%+v\n", uuid.UUID(w.corePluginInterface.GetId()).String())

	meta, ok := w.GetSpaceAttributeValue(
		entry.NewAttributeID(
			uuid.UUID(w.corePluginInterface.GetId()), "world_meta",
		),
	)
	if !ok {
		w.metaMsg.Store(nil)
		return nil
	}
	metaMap := (map[string]any)(*meta)

	utils.MapDecode(metaMap, &w.metaData)

	fmt.Printf("Meta: %+v\n", w.metaData)
	//TODO: Ut is all ugly with circular deps
	dec := make([]message.DecorationMetadata, len(w.metaData.Decorations))
	for i, decoration := range w.metaData.Decorations {
		dec[i].AssetID = decoration.AssetID
		dec[i].Position = decoration.Position
	}

	w.metaMsg.Store(
		message.GetBuilder().MsgSetWorld(
			w.GetID(), w.GetName(), w.metaData.AvatarController, w.metaData.SkyboxController, w.metaData.LOD,
			dec,
		),
	)

	return nil
}

func (w *World) Save() error {
	w.log.Infof("Saving world: %s", w.GetID())

	spaces := w.GetSpaces(true)

	entries := make([]*entry.Space, 0, len(spaces))
	for _, space := range spaces {
		entries = append(entries, space.GetEntry())
	}

	if err := w.db.SpacesUpsertSpaces(w.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert spaces")
	}

	w.log.Infof("World saved: %s", w.GetID())

	return nil
}
