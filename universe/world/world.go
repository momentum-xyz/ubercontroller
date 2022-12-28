package world

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/posbus-protocol/posbus"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/calendar"
	"github.com/momentum-xyz/ubercontroller/universe/space"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.World = (*World)(nil)

type World struct {
	*space.Space
	ctx              context.Context
	cancel           context.CancelFunc
	log              *zap.SugaredLogger
	db               database.DB
	pluginController *mplugin.PluginController
	//corePluginInstance  mplugin.PluginInstance
	corePluginInterface mplugin.PluginInterface
	metaMsg             atomic.Pointer[websocket.PreparedMessage]
	metaData            Metadata
	settings            atomic.Pointer[universe.WorldSettings]
	allSpaces           *generic.SyncMap[uuid.UUID, universe.Space]
	calendar            *calendar.Calendar
}

func NewWorld(id uuid.UUID, db database.DB) *World {
	world := &World{
		db:        db,
		allSpaces: generic.NewSyncMap[uuid.UUID, universe.Space](0),
	}
	world.Space = space.NewSpace(id, db, world)
	world.settings.Store(&universe.WorldSettings{})
	world.pluginController = mplugin.NewPluginController(id)
	//world.corePluginInstance, _ = world.pluginController.AddPlugin(world.GetID(), world.corePluginInitFunc)
	world.pluginController.AddPlugin(universe.GetSystemPluginID(), world.corePluginInitFunc)
	world.calendar = calendar.NewCalendar(world)
	return world
}

func (w *World) GetID() uuid.UUID {
	return w.Space.GetID()
}

func (w *World) corePluginInitFunc(pi mplugin.PluginInterface) (mplugin.PluginInstance, error) {
	instance := CorePluginInstance{PluginInterface: pi}
	w.corePluginInterface = pi
	return instance, nil
}

func (w *World) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.log = log

	if err := w.calendar.Initialize(ctx); err != nil {
		return errors.WithMessage(err, "failed to initialize calendar")
	}

	return w.Space.Initialize(ctx)
}

func (w *World) ToSpace() universe.Space {
	return w.Space
}

func (w *World) GetSettings() *universe.WorldSettings {
	return w.settings.Load()
}

func (w *World) GetCalendar() universe.Calendar {
	return w.calendar
}

func (w *World) Run() error {
	go func() {
		go func() {
			if err := w.runSpaces(); err != nil {
				w.log.Error(errors.WithMessagef(err, "World: Run: failed to run spaces: %s", w.GetID()))
			}
		}()
		go w.calendar.Run()
		ticker := time.NewTicker(500 * time.Millisecond)

		defer func() {
			w.calendar.Stop()
			ticker.Stop()
			if err := w.stopSpaces(); err != nil {
				w.log.Error(errors.WithMessagef(err, "World: Run: failed to stop spaces: %s", w.GetID()))
			}
		}()

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
			case <-w.ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (w *World) Stop() error {
	w.cancel()
	return nil
}

func (w *World) runSpaces() error {
	w.allSpaces.Mu.RLock()
	defer w.allSpaces.Mu.RUnlock()

	var errs *multierror.Error
	for _, space := range w.allSpaces.Data {
		if err := space.Run(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to run space: %s", space.GetID()))
		}
		space.SetEnabled(true)
	}

	return errs.ErrorOrNil()
}

// TODO: optimize
func (w *World) stopSpaces() error {
	w.allSpaces.Mu.RLock()
	defer w.allSpaces.Mu.RUnlock()

	var errs *multierror.Error
	for _, space := range w.allSpaces.Data {
		if err := space.Stop(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to stop space: %s", space.GetID()))
		}
		space.SetEnabled(false)
	}

	return errs.ErrorOrNil()
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
		w.Send(msg.WebsocketMessage(), true)
	}
}

func (w *World) Load() error {
	w.log.Infof("Loading world: %s...", w.GetID())

	entry, err := w.db.SpacesGetSpaceByID(w.ctx, w.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get space by id")
	}

	if err := w.LoadFromEntry(entry, true); err != nil {
		return errors.WithMessage(err, "failed to load from entry")
	}
	if err := w.UpdateChildrenPosition(true); err != nil {
		return errors.WithMessage(err, "failed to update children position")
	}
	if err := w.Update(true); err != nil {
		w.log.Error(errors.WithMessagef(err, "failed to update world"))
	}

	w.log.Infof("World loaded: %s", w.GetID())

	return nil
}

func (w *World) Update(recursive bool) error {
	if err := w.UpdateWorldMetadata(); err != nil {
		w.log.Error(errors.WithMessagef(err, "World: Update: failed to update world metadata: %s", w.GetID()))
	}
	if err := w.UpdateWorldSettings(); err != nil {
		w.log.Error(errors.WithMessagef(err, "World: Update: failed to update world settings: %s", w.GetID()))
	}

	return w.Space.Update(recursive)
}

func (w *World) UpdateWorldSettings() error {
	value, ok := w.GetSpaceAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.Settings.Name),
	)
	if !ok || value == nil {
		return errors.Errorf("space attribute not found")
	}

	var settings universe.WorldSettings
	if err := utils.MapDecode(*value, &settings); err != nil {
		return errors.WithMessage(err, "failed to decode map")
	}

	w.settings.Store(&settings)

	return nil
}

func (w *World) UpdateWorldMetadata() error {
	meta, ok := w.GetSpaceAttributes().GetValue(
		entry.NewAttributeID(
			uuid.UUID(w.corePluginInterface.GetId()), universe.ReservedAttributes.World.Meta.Name,
		),
	)

	if ok {
		if err := utils.MapDecode(*meta, &w.metaData); err != nil {
			return errors.WithMessage(err, "failed to decode meta")
		}
	} else {
		// TODO: print warning and call stack here
		w.metaData = Metadata{}
	}

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

	spaces := w.GetAllSpaces()

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

func (w *World) GetAllSpaces() map[uuid.UUID]universe.Space {
	w.allSpaces.Mu.RLock()
	defer w.allSpaces.Mu.RUnlock()

	spaces := make(map[uuid.UUID]universe.Space, len(w.allSpaces.Data))

	for id, space := range w.allSpaces.Data {
		spaces[id] = space
	}

	return spaces
}
