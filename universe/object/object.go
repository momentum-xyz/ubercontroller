package object

import (
	"context"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Object = (*Object)(nil)

type Object struct {
	id       uuid.UUID
	world    universe.World
	ctx      context.Context
	log      *zap.SugaredLogger
	db       database.DB
	enabled  atomic.Bool
	Users    *generic.SyncMap[uuid.UUID, universe.User]
	Children *generic.SyncMap[uuid.UUID, universe.Object]
	//Mu               sync.RWMutex
	Mu               deadlock.RWMutex
	ownerID          uuid.UUID
	position         *cmath.SpacePosition
	options          *entry.ObjectOptions
	Parent           universe.Object
	asset2d          universe.Asset2d
	asset3d          universe.Asset3d
	objectType       universe.ObjectType
	effectiveOptions *entry.ObjectOptions
	objectAttributes *objectAttributes // WARNING: the Object is sharing the same mutex ("Mu") with it

	spawnMsg          atomic.Pointer[websocket.PreparedMessage]
	attributesMsg     *generic.SyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]]
	renderTextureMap  *generic.SyncMap[string, string]
	textMsg           atomic.Pointer[websocket.PreparedMessage]
	actualPosition    atomic.Pointer[cmath.SpacePosition]
	broadcastPipeline chan *websocket.PreparedMessage
	messageAccept     atomic.Bool
	numSendsQueued    atomic.Int64

	lockedBy atomic.Value

	// TODO: replace theta with full calculation of orientation, once Unity is read
	theta float64
}

func NewObject(id uuid.UUID, db database.DB, world universe.World) *Object {
	object := &Object{
		id:               id,
		db:               db,
		Users:            generic.NewSyncMap[uuid.UUID, universe.User](0),
		Children:         generic.NewSyncMap[uuid.UUID, universe.Object](0),
		attributesMsg:    generic.NewSyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]](0),
		renderTextureMap: generic.NewSyncMap[string, string](0),
		world:            world,
	}
	object.objectAttributes = newObjectAttributes(object)

	return object
}

func (s *Object) GetID() uuid.UUID {
	return s.id
}

func (s *Object) GetEnabled() bool {
	return s.enabled.Load()
}

func (s *Object) SetEnabled(enabled bool) {
	s.enabled.Store(enabled)
}

func (s *Object) GetName() string {
	name := "unknown"
	value, ok := s.GetObjectAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Name.Name),
	)
	if !ok || value == nil {
		return name
	}
	return utils.GetFromAnyMap(*value, universe.ReservedAttributes.Object.Name.Key, name)
}

func (s *Object) SetName(name string, updateDB bool) error {
	if _, err := s.GetObjectAttributes().Upsert(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Name.Name),
		modify.MergeWith(entry.NewAttributePayload(
			&entry.AttributeValue{
				universe.ReservedAttributes.Object.Name.Key: name,
			},
			nil),
		), updateDB,
	); err != nil {
		return errors.WithMessage(err, "failed to upsert object attribute")
	}
	return nil
}

func (s *Object) GetObjectAttributes() universe.ObjectAttributes {
	return s.objectAttributes
}

func (s *Object) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	s.ctx = ctx
	s.log = log
	s.numSendsQueued.Store(chanIsClosed)
	s.lockedBy.Store(uuid.Nil)

	newPos := cmath.SpacePosition{Location: *new(cmath.Vec3), Rotation: *new(cmath.Vec3), Scale: *new(cmath.Vec3)}
	s.actualPosition.Store(&newPos)

	return nil
}

func (s *Object) GetWorld() universe.World {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.world
}

func (s *Object) GetParent() universe.Object {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Parent
}

func (s *Object) SetParent(parent universe.Object, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if parent == s {
		return errors.Errorf("object can't be a parent of itself")
	} else if parent != nil && parent.GetWorld().GetID() != s.world.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.world.GetID())
	}

	if updateDB {
		if parent == nil {
			return errors.Errorf("parent is nil")
		}
		if err := s.db.GetObjectsDB().UpdateObjectParentID(s.ctx, s.GetID(), parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.Parent = parent

	return nil
}

func (s *Object) GetOwnerID() uuid.UUID {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.ownerID
}

func (s *Object) SetOwnerID(ownerID uuid.UUID, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		if err := s.db.GetObjectsDB().UpdateObjectOwnerID(s.ctx, s.GetID(), ownerID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.ownerID = ownerID

	return nil
}

func (s *Object) GetAsset2D() universe.Asset2d {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.asset2d
}

func (s *Object) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		var asset2dID *uuid.UUID
		if asset2d != nil {
			asset2dID = utils.GetPTR(asset2d.GetID())
		}
		if err := s.db.GetObjectsDB().UpdateObjectAsset2dID(s.ctx, s.GetID(), asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d

	return nil
}

func (s *Object) GetAsset3D() universe.Asset3d {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.asset3d
}

func (s *Object) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		var asset3dID *uuid.UUID
		if asset3d != nil {
			asset3dID = utils.GetPTR(asset3d.GetID())
		}
		if err := s.db.GetObjectsDB().UpdateObjectAsset3dID(s.ctx, s.GetID(), asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset3d = asset3d

	return nil
}

func (s *Object) GetObjectType() universe.ObjectType {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.objectType
}

func (s *Object) SetObjectType(objectType universe.ObjectType, updateDB bool) error {
	if objectType == nil {
		return errors.Errorf("object type is nil")
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		if err := s.db.GetObjectsDB().UpdateObjectObjectTypeID(s.ctx, s.GetID(), objectType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.objectType = objectType
	s.dropCache()

	return nil
}

func (s *Object) GetOptions() *entry.ObjectOptions {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.options
}

func (s *Object) SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	options, err := modifyFn(s.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := s.db.GetObjectsDB().UpdateObjectOptions(s.ctx, s.GetID(), options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options
	s.dropCache()

	return options, nil
}

func (s *Object) GetEffectiveOptions() *entry.ObjectOptions {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if s.effectiveOptions == nil {
		effectiveOptions, err := merge.Auto(s.options, s.objectType.GetOptions())
		if err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Object: GetEffectiveOptions: failed to merge object effective options: %s", s.GetID(),
				),
			)
			return nil
		}

		s.effectiveOptions = effectiveOptions
	}

	return s.effectiveOptions
}

func (s *Object) DropCache() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.dropCache()
}

func (s *Object) dropCache() {
	s.effectiveOptions = nil
}

func (s *Object) GetEntry() *entry.Object {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	entry := &entry.Object{
		ObjectID: s.id,
		OwnerID:  &s.ownerID,
		Options:  s.options,
		Position: s.position,
	}
	if s.objectType != nil {
		entry.ObjectTypeID = utils.GetPTR(s.objectType.GetID())
	}
	if s.Parent != nil {
		entry.ParentID = utils.GetPTR(s.Parent.GetID())
	}
	if s.asset2d != nil {
		entry.Asset2dID = utils.GetPTR(s.asset2d.GetID())
	}
	if s.asset3d != nil {
		entry.Asset3dID = utils.GetPTR(s.asset3d.GetID())
	}

	return entry
}

func (s *Object) Run() error {
	s.numSendsQueued.Store(0)
	s.broadcastPipeline = make(chan *websocket.PreparedMessage, 100)

	go func() {
		defer func() {
			ns := s.numSendsQueued.Swap(chanIsClosed)
			for i := int64(0); i < ns; i++ {
				<-s.broadcastPipeline
			}
			close(s.broadcastPipeline)
		}()

		for {
			select {
			case message := <-s.broadcastPipeline:
				s.numSendsQueued.Add(-1)
				if message == nil {
					return
				}

				s.performBroadcast(message)
			case <-s.ctx.Done():
				s.Stop()
			}
		}
	}()

	return nil
}

func (s *Object) Stop() error {
	ns := s.numSendsQueued.Add(1)
	if ns >= 0 {
		s.broadcastPipeline <- nil
	}
	return nil
}

func (s *Object) Update(recursive bool) error {
	s.UpdateSpawnMessage()

	if s.GetEnabled() {
		world := s.GetWorld()
		if world != nil {
			world.Send(s.spawnMsg.Load(), true)
			s.SendTextures(
				func(msg *websocket.PreparedMessage) error {
					return world.Send(msg, false)
				}, false,
			)
		}
	}

	if !recursive {
		return nil
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		if err := child.Update(recursive); err != nil {
			return errors.WithMessagef(err, "failed to update child: %s", child.GetID())
		}
	}

	return nil
}

func (s *Object) LoadFromEntry(entry *entry.Object, recursive bool) error {
	s.log.Debugf("Loading object %s...", entry.ObjectID)

	if entry.ObjectID != s.GetID() {
		return errors.Errorf("object ids mismatch: %s != %s", entry.ObjectID, s.GetID())
	}

	group, _ := errgroup.WithContext(s.ctx)
	group.Go(s.GetObjectAttributes().Load)
	group.Go(
		func() error {
			if err := s.loadSelfData(entry); err != nil {
				return errors.WithMessage(err, "failed to load self data")
			}
			if err := s.loadDependencies(entry); err != nil {
				return errors.WithMessage(err, "failed to load dependencies")
			}
			if err := s.SetPosition(entry.Position, false); err != nil {
				return errors.WithMessage(err, "failed to set position")
			}

			if !recursive {
				return nil
			}

			entries, err := s.db.GetObjectsDB().GetObjectsByParentID(s.ctx, s.GetID())
			if err != nil {
				return errors.WithMessagef(err, "failed to get objects by parent id: %s", s.GetID())
			}

			for i := range entries {
				child, err := s.CreateObject(entries[i].ObjectID)
				if err != nil {
					return errors.WithMessagef(err, "failed to create new object: %s", entries[i].ObjectID)
				}
				if err := child.LoadFromEntry(entries[i], recursive); err != nil {
					return errors.WithMessagef(err, "failed to load object from entry: %s", entries[i].ObjectID)
				}
			}

			return nil
		},
	)
	return group.Wait()
}

func (s *Object) loadSelfData(objectEntry *entry.Object) error {
	if err := s.SetOwnerID(*objectEntry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner id: %s", objectEntry.OwnerID)
	}
	if _, err := s.SetOptions(modify.MergeWith(objectEntry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}
	return nil
}

func (s *Object) loadDependencies(entry *entry.Object) error {
	node := universe.GetNode()

	objectType, ok := node.GetObjectTypes().GetObjectType(*entry.ObjectTypeID)
	if !ok {
		return errors.Errorf("failed to get object type: %s", entry.ObjectTypeID)
	}
	if err := s.SetObjectType(objectType, false); err != nil {
		return errors.WithMessagef(err, "failed to set object type: %s", entry.ObjectTypeID)
	}

	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("failed to get asset 2d: %s", entry.Asset2dID)
		}
		if err := s.SetAsset2D(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("failed to get asset 3d: %s", entry.Asset3dID)
		}
		if err := s.SetAsset3D(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", entry.Asset3dID)
		}
	}

	return nil
}

func (s *Object) UpdateSpawnMessage() error {
	world := s.GetWorld()
	if world == nil {
		return errors.Errorf("world is empty")
	}

	parentID := uuid.Nil
	parent := s.GetParent()
	if parent != nil {
		parentID = parent.GetID()
	}

	asset3dID := uuid.Nil
	asset3d := s.GetAsset3D()
	objectType := s.GetObjectType()
	assetFormat := dto.AddressableAssetType
	if asset3d == nil && objectType != nil {
		asset3d = objectType.GetAsset3d()
	}
	if asset3d != nil {
		asset3dID = asset3d.GetID()
		asset3dMeta := asset3d.GetMeta()
		if asset3dMeta != nil {
			// TODO: make GetMeta return struct type
			metaData := struct {
				Type int `json:"type"`
			}{}
			utils.MapDecode(*asset3dMeta, &metaData)
			assetFormat = dto.Asset3dType(metaData.Type)
		}
	}

	uuidNilPtr := utils.GetPTR(uuid.Nil)
	falsePtr := utils.GetPTR(false)
	truePtr := utils.GetPTR(true)
	opts := s.GetEffectiveOptions()
	msg := message.GetBuilder().MsgObjectDefinition(
		message.ObjectDefinition{
			ObjectID:         s.GetID(),
			ParentID:         parentID,
			AssetType:        asset3dID,
			AssetFormat:      assetFormat,
			Name:             s.GetName(),
			Position:         *s.GetActualPosition(),
			Editable:         *utils.GetFromAny(opts.Editable, truePtr),
			TetheredToParent: true,
			Minimap:          *utils.GetFromAny(opts.Minimap, falsePtr),
			InfoUI:           *utils.GetFromAny(opts.InfoUIID, uuidNilPtr),
		},
	)
	s.spawnMsg.Store(msg)

	return nil
}

func (s *Object) GetSpawnMessage() *websocket.PreparedMessage {
	return s.spawnMsg.Load()
}

func (s *Object) SendSpawnMessage(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	sendFn(s.spawnMsg.Load())
	//time.Sleep(time.Millisecond * 100)
	if !recursive {
		return
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		child.SendSpawnMessage(sendFn, recursive)
	}

}

func (s *Object) SendTextures(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	msg := s.textMsg.Load()
	if msg != nil {
		sendFn(msg)
	}

	if !recursive {
		return
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		child.SendTextures(sendFn, recursive)
	}
}

// QUESTION: why this method is never called?
func (s *Object) SendAttributes(sendFn func(*websocket.PreparedMessage), recursive bool) {
	s.attributesMsg.Mu.RLock()
	for _, g := range s.attributesMsg.Data {
		for _, a := range g.Data {
			sendFn(a)
		}
	}
	s.attributesMsg.Mu.RUnlock()

	sendFn(s.spawnMsg.Load())

	if !recursive {
		return
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, child := range s.Children.Data {
		child.SendAttributes(sendFn, recursive)
	}
}

// QUESTION: why this method is never called?
func (s *Object) SetAttributesMsg(kind, name string, msg *websocket.PreparedMessage) {
	m, ok := s.attributesMsg.Load(kind)
	if !ok {
		m = generic.NewSyncMap[string, *websocket.PreparedMessage](0)
		s.attributesMsg.Store(kind, m)
	}
	m.Store(name, msg)
}

func (s *Object) LockUnityObject(user universe.User, state uint32) bool {
	if state == 1 {
		return s.lockedBy.CompareAndSwap(uuid.Nil, user.GetID())
	} else {
		return s.lockedBy.CompareAndSwap(user.GetID(), uuid.Nil)
	}
}
