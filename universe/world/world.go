package world

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"reflect"
)

var _ universe.World = (*World)(nil)

type DecorationMetadataNew struct {
	AssetID  uuid.UUID  `json:"asset_id" db:"asset_id"`
	Position cmath.Vec3 `json:"position" db:"position"`
	rotation cmath.Vec3
}

type WorldMeta struct {
	LOD              []uint32                `json:"lod" db:"lod"`
	Decorations      []DecorationMetadataNew `json:"decorations,omitempty" db:"decorations,omitempty"`
	AvatarController uuid.UUID               `json:"avatar_controller" db:"avatar_controller"`
	SkyboxController uuid.UUID               `json:"skybox_controller" db:"skybox_controller"`
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
	metaMsg             *websocket.PreparedMessage
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
	return nil
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
		w.metaMsg = nil
		return nil
	}
	metaMap := (map[string]any)(*meta)

	if v, ok := metaMap["skybox_controller"]; ok {
		w.metaData.SkyboxController = uuid.MustParse(v.(string))
	}

	if v, ok := metaMap["avatar_controller"]; ok {
		w.metaData.AvatarController = uuid.MustParse(v.(string))
	}

	lods := utils.GetFromAnyMap(metaMap, "lod", make([]any, 0))
	w.metaData.LOD = make([]uint32, len(lods))
	for i := 0; i < len(lods); i++ {
		w.metaData.LOD[i] = uint32(lods[i].(float64))
	}

	//q := WorldMeta{}
	//
	//mapstructure.Decode(metaMap, &q)

	decs := utils.GetFromAnyMap(metaMap, "decorations", make([]any, 0))

	fmt.Printf("%+v\n", decs)
	if len(decs) > 0 {
		fmt.Printf("%+v\n", reflect.ValueOf(decs[0]).Type())
	}

	fmt.Printf("Meta: %+v\n", w.metaData)

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
