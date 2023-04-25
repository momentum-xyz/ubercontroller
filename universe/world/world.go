package world

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/mplugin"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/calendar"
	"github.com/momentum-xyz/ubercontroller/universe/object"
	"github.com/momentum-xyz/ubercontroller/utils"
)

var _ universe.World = (*World)(nil)

// MaxPosUpdateInterval : send user position at least ones per 5 min, even if user is not moving
const MaxPosUpdateInterval = 60 * 5

const PosUpdateInterval = 500 * time.Millisecond

type World struct {
	*object.Object
	ctx              context.Context
	log              *zap.SugaredLogger
	db               database.DB
	cancel           context.CancelFunc
	pluginController *mplugin.PluginController
	//corePluginInstance  mplugin.PluginInstance
	corePluginInterface mplugin.PluginInterface
	metaMsg             atomic.Pointer[websocket.PreparedMessage]
	metaData            Metadata
	settings            atomic.Pointer[universe.WorldSettings]
	allObjects          *generic.SyncMap[umid.UMID, universe.Object]
	calendar            *calendar.Calendar
	skyBoxMsg           atomic.Pointer[websocket.PreparedMessage]
	lastPosUpdate       int64
}

func (w *World) LockUIObject(user universe.User, state uint32) bool {
	//TODO implement me
	panic("implement me")
}

func (w *World) GetTotalStake() uint8 {
	//TODO implement me
	panic("implement me")
}

func (w *World) TempSetSkybox(msg *websocket.PreparedMessage) {
	w.skyBoxMsg.Store(msg)
}

func (w *World) TempGetSkybox() *websocket.PreparedMessage {
	return w.skyBoxMsg.Load()
}

func NewWorld(id umid.UMID, db database.DB) *World {
	world := &World{
		db:         db,
		allObjects: generic.NewSyncMap[umid.UMID, universe.Object](0),
	}
	world.Object = object.NewObject(id, db, world)
	world.settings.Store(&universe.WorldSettings{})
	world.pluginController = mplugin.NewPluginController(id)
	//world.corePluginInstance, _ = world.pluginController.AddPlugin(world.GetID(), world.corePluginInitFunc)
	world.pluginController.AddPlugin(universe.GetSystemPluginID(), world.corePluginInitFunc)
	world.calendar = calendar.NewCalendar(world)
	return world
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
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.Errorf("failed to get config from context: %T", ctx.Value(types.ConfigContextKey))
	}

	w.ctx, w.cancel = context.WithCancel(ctx)
	w.log = log

	if err := w.calendar.Initialize(ctx); err != nil {
		return errors.WithMessage(err, "failed to initialize calendar")
	}

	return w.ToObject().Initialize(ctx)
}

func (w *World) ToObject() universe.Object {
	return w.Object
}

func (w *World) GetSettings() *universe.WorldSettings {
	return w.settings.Load()
}

func (w *World) GetCalendar() universe.Calendar {
	return w.calendar
}

func (w *World) SetParent(parent universe.Object, updateDB bool) error {
	if parent == nil {
		return errors.Errorf("parent is nil")
	} else if parent.GetID() != universe.GetNode().GetID() {
		return errors.Errorf("parent is not the node")
	}

	w.Object.Mu.Lock()
	defer w.Object.Mu.Unlock()

	if updateDB {
		if err := w.db.GetObjectsDB().UpdateObjectParentID(w.ctx, w.GetID(), parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	w.Object.Parent = parent

	return nil
}

func (w *World) Run() error {
	go func() {
		go func() {
			if err := w.runObjects(); err != nil {
				w.log.Error(errors.WithMessagef(err, "World: Run: failed to run objects: %s", w.GetID()))
			}
		}()
		go w.calendar.Run()
		ticker := time.NewTicker(PosUpdateInterval)

		defer func() {
			w.calendar.Stop()
			ticker.Stop()
			if err := w.stopObjects(); err != nil {
				w.log.Error(errors.WithMessagef(err, "World: Run: failed to stop objects: %s", w.GetID()))
			}
		}()

		for {
			select {
			case <-ticker.C:
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

func (w *World) runObjects() error {
	w.allObjects.Mu.RLock()
	defer w.allObjects.Mu.RUnlock()

	var errs *multierror.Error
	for _, object := range w.allObjects.Data {
		if err := object.Run(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to run object: %s", object.GetID()))
		}
		object.SetEnabled(true)
	}

	return errs.ErrorOrNil()
}

// TODO: optimize
func (w *World) stopObjects() error {
	w.allObjects.Mu.RLock()
	defer w.allObjects.Mu.RUnlock()

	var errs *multierror.Error
	for _, object := range w.allObjects.Data {
		if err := object.Stop(); err != nil {
			errs = multierror.Append(errs, errors.WithMessagef(err, "failed to stop object: %s", object.GetID()))
		}
		object.SetEnabled(false)
	}

	return errs.ErrorOrNil()
}

func (w *World) broadcastPositions() {
	w.Users.Mu.RLock()
	numClients := len(w.Users.Data)
	currentTime := time.Now().Unix()

	// Need to do some 'reasonable' batching.
	// - fit inside max message size of a receiving client.
	// - something reasonable to process in frontend and not lockup UI thread for too long?
	// Size is about 40b per client with some overhead,
	// 'default' setting of our websocket client is 32kb
	// so as a start, something that fits and won't often be triggered?
	batchSize := 768

	var msgBatches []posbus.UsersTransformList
	if numClients > 0 {
		uTransforms := make([]posbus.UserTransform, 0)
		for _, u := range w.Users.Data {
			if (u.GetLastPosTime() >= w.lastPosUpdate) || ((currentTime - u.GetLastSendPosTime()) > MaxPosUpdateInterval) {
				u.SetLastSendPosTime(currentTime)
				uTransforms = append(uTransforms, posbus.UserTransform{ID: u.GetID(), Transform: *u.GetTransform()})
			}
		}
		nrUpdates := len(uTransforms)
		msgBatches = make([]posbus.UsersTransformList, 0, (nrUpdates+batchSize-1)/batchSize)

		generic.NewButcher(uTransforms).HandleBatchesSync(
			batchSize,
			func(batch []posbus.UserTransform) error {
				msg := posbus.UsersTransformList{}
				msg.Value = batch
				msgBatches = append(msgBatches, msg)
				return nil
			},
		)
	}
	w.lastPosUpdate = currentTime

	w.Users.Mu.RUnlock()
	for _, msg := range msgBatches {
		w.Send(posbus.WSMessage(&msg), true)
	}
}

// Send posbus.AddUsers containing all current users in the world (including themself).
// Accepts a function, which should be the function to send the message to a user.
//
// Similar to Object.SendSpawnMessage, but not prepared like objects (stored on world).
// This changes more often and would require some fine-grained hooks into the add/remove user logic.
// (Also it just changed, since this new user was added)
func (w *World) SendUsersSpawnMessage(sendFn func(*websocket.PreparedMessage) error) {
	// See broadcastPositions, same logic, just different contents
	w.Users.Mu.RLock()
	numClients := len(w.Users.Data)
	batchSize := 100 // Sane size for UserData? contains variable name string.
	var msgBatches []posbus.AddUsers
	if numClients > 0 {
		uDatas := make([]posbus.UserData, 0)
		for _, u := range w.Users.Data {
			uDatas = append(uDatas, *u.GetUserDefinition())
		}
		nrUpdates := len(uDatas)
		msgBatches = make([]posbus.AddUsers, 0, (nrUpdates+batchSize-1)/batchSize)

		generic.NewButcher(uDatas).HandleBatchesSync(
			batchSize,
			func(batch []posbus.UserData) error {
				msg := posbus.AddUsers{}
				msg.Users = batch
				msgBatches = append(msgBatches, msg)
				return nil
			},
		)
	}
	w.Users.Mu.RUnlock()
	for _, msg := range msgBatches {
		sendFn(posbus.WSMessage(&msg))
	}
}

func (w *World) Load() error {
	w.log.Infof("Loading world: %s...", w.GetID())

	worldEntry, err := w.db.GetObjectsDB().GetObjectByID(w.ctx, w.GetID())
	if err != nil {
		return errors.WithMessage(err, "failed to get object by umid")
	}

	if err := w.LoadFromEntry(worldEntry, true); err != nil {
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

func (w *World) Save() error {
	w.log.Infof("Saving world: %s...", w.GetID())

	if err := w.ToObject().Save(); err != nil {
		return errors.WithMessage(err, "failed to save world object")
	}

	w.log.Infof("World saved: %s", w.GetID())

	return nil
}

func (w *World) Update(recursive bool) error {
	if err := w.UpdateWorldMetadata(); err != nil {
		w.log.Error(errors.WithMessagef(err, "World: Update: failed to update world metadata: %s", w.GetID()))
	}
	if err := w.UpdateWorldSettings(); err != nil {
		w.log.Error(errors.WithMessagef(err, "World: Update: failed to update world settings: %s", w.GetID()))
	}

	return w.ToObject().Update(recursive)
}

func (w *World) UpdateWorldSettings() error {
	value, ok := w.GetObjectAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.Settings.Name),
	)
	if !ok || value == nil {
		return errors.Errorf("object attribute not found")
	}

	var settings universe.WorldSettings
	if err := utils.MapDecode(*value, &settings); err != nil {
		return errors.WithMessage(err, "failed to decode map")
	}

	w.settings.Store(&settings)

	return nil
}

func (w *World) UpdateWorldMetadata() error {
	meta, ok := w.GetObjectAttributes().GetValue(
		entry.NewAttributeID(
			umid.UMID(w.corePluginInterface.GetId()), universe.ReservedAttributes.World.Meta.Name,
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

	w.metaMsg.Store(
		posbus.WSMessage(&posbus.SetWorld{ID: w.GetID(), Name: w.GetName(), Avatar: w.metaData.AvatarController, Avatar3DAssetID: w.metaData.AvatarController, Owner: w.GetOwnerID()}),
	)

	return nil
}
