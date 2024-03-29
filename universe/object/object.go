package object

import (
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/sasha-s/go-deadlock"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/media"
	"github.com/momentum-xyz/ubercontroller/pkg/posbus"
	"github.com/momentum-xyz/ubercontroller/seed"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

var _ universe.Object = (*Object)(nil)

type Object struct {
	id               umid.UMID
	world            universe.World
	ctx              types.NodeContext
	log              *zap.SugaredLogger
	CFG              *config.Config
	db               database.DB
	enabled          atomic.Bool
	Users            *generic.SyncMap[umid.UMID, universe.User]
	Children         *generic.SyncMap[umid.UMID, universe.Object]
	Mu               deadlock.RWMutex
	ownerID          umid.UMID
	transform        *cmath.Transform
	options          *entry.ObjectOptions
	Parent           universe.Object
	asset2d          universe.Asset2d
	asset3d          universe.Asset3d
	objectType       universe.ObjectType
	effectiveOptions *entry.ObjectOptions
	activities       *entry.Activity
	objectAttributes *objectAttributes // WARNING: the Object is sharing the same mutex ("Mu") with it

	media *media.Media

	spawnMsg      atomic.Pointer[websocket.PreparedMessage]
	attributesMsg *generic.SyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]]

	renderDataMap *generic.SyncMap[entry.SlotType, *posbus.StringAnyMap]
	dataMsg       atomic.Pointer[websocket.PreparedMessage]

	actualPosition    atomic.Pointer[cmath.Transform]
	broadcastPipeline chan *websocket.PreparedMessage
	messageAccept     atomic.Bool
	numSendsQueued    atomic.Int64

	lockedBy atomic.Value

	createdAt time.Time
	updatedAt time.Time

	// TODO: replace theta with full calculation of orientation, once Unity is read
	theta float64
}

func NewObject(id umid.UMID, db database.DB, world universe.World, media *media.Media) *Object {
	now := time.Now()
	object := &Object{
		id:            id,
		db:            db,
		Users:         generic.NewSyncMap[umid.UMID, universe.User](0),
		Children:      generic.NewSyncMap[umid.UMID, universe.Object](0),
		attributesMsg: generic.NewSyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]](0),
		renderDataMap: generic.NewSyncMap[entry.SlotType, *posbus.StringAnyMap](0),
		world:         world,
		media:         media,
		createdAt:     now,
		updatedAt:     now,
	}
	object.objectAttributes = newObjectAttributes(object)

	return object
}

func (o *Object) GetID() umid.UMID {
	return o.id
}

func (o *Object) GetEnabled() bool {
	return o.enabled.Load()
}

func (o *Object) SetEnabled(enabled bool) {
	o.enabled.Store(enabled)
}

func (o *Object) GetDescription() string {
	defaultDescription := "No Description Set"
	value, ok := o.GetObjectAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Description.Name),
	)
	if !ok || value == nil {
		return defaultDescription
	}
	return utils.GetFromAnyMap(*value, universe.ReservedAttributes.Object.Description.Key, defaultDescription)
}

func (o *Object) GetChildIDs() []umid.UMID {
	ids := make([]umid.UMID, 0, o.Children.Len())

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		ids = append(ids, child.GetID())
	}

	return ids
}

func (o *Object) GetName() string {
	name := o.GetID().String()
	value, ok := o.GetObjectAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Name.Name),
	)
	if !ok || value == nil {
		return name
	}
	return utils.GetFromAnyMap(*value, universe.ReservedAttributes.Object.Name.Key, name)
}

func (o *Object) SetName(name string, updateDB bool) error {
	if _, err := o.GetObjectAttributes().Upsert(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.Name.Name),
		modify.MergeWith(
			entry.NewAttributePayload(
				&entry.AttributeValue{
					universe.ReservedAttributes.Object.Name.Key: name,
				},
				nil,
			),
		), updateDB,
	); err != nil {
		return errors.WithMessage(err, "failed to upsert object attribute")
	}
	return nil
}

func (o *Object) GetObjectAttributes() universe.ObjectAttributes {
	return o.objectAttributes
}

func (o *Object) Initialize(ctx types.NodeContext) error {
	o.ctx = ctx
	o.log = ctx.Logger()
	o.CFG = ctx.Config()
	o.numSendsQueued.Store(chanIsClosed)
	o.lockedBy.Store(umid.Nil)

	newPos := cmath.Transform{Position: *new(cmath.Vec3), Rotation: *new(cmath.Vec3), Scale: *new(cmath.Vec3)}
	o.actualPosition.Store(&newPos)

	return nil
}

func (o *Object) GetWorld() universe.World {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.world
}

func (o *Object) GetParent() universe.Object {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.Parent
}

func (o *Object) SetParent(parent universe.Object, updateDB bool) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if parent == o {
		return errors.Errorf("object can't be a parent of itself")
	} else if parent != nil && parent.GetWorld().GetID() != o.world.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), o.world.GetID())
	}

	if updateDB {
		if parent == nil {
			return errors.Errorf("parent is nil")
		}
		if err := o.db.GetObjectsDB().UpdateObjectParentID(o.ctx, o.GetID(), parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	o.Parent = parent

	return nil
}

func (o *Object) GetOwnerID() umid.UMID {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.ownerID
}

func (o *Object) SetOwnerID(ownerID umid.UMID, updateDB bool) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if updateDB {
		if err := o.db.GetObjectsDB().UpdateObjectOwnerID(o.ctx, o.GetID(), ownerID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	o.ownerID = ownerID

	return nil
}

func (o *Object) GetAsset2D() universe.Asset2d {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.asset2d
}

func (o *Object) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if updateDB {
		var asset2dID *umid.UMID
		if asset2d != nil {
			asset2dID = utils.GetPTR(asset2d.GetID())
		}
		if err := o.db.GetObjectsDB().UpdateObjectAsset2dID(o.ctx, o.GetID(), asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	o.asset2d = asset2d

	return nil
}

func (o *Object) GetAsset3D() universe.Asset3d {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.asset3d
}

func (o *Object) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if updateDB {
		var asset3dID *umid.UMID
		if asset3d != nil {
			asset3dID = utils.GetPTR(asset3d.GetID())
		}
		if err := o.db.GetObjectsDB().UpdateObjectAsset3dID(o.ctx, o.GetID(), asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	o.asset3d = asset3d

	return nil
}

func (o *Object) GetObjectType() universe.ObjectType {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.objectType
}

func (o *Object) SetObjectType(objectType universe.ObjectType, updateDB bool) error {
	if objectType == nil {
		return errors.Errorf("object type is nil")
	}

	o.Mu.Lock()
	defer o.Mu.Unlock()

	if updateDB {
		if err := o.db.GetObjectsDB().UpdateObjectObjectTypeID(o.ctx, o.GetID(), objectType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	o.objectType = objectType
	o.invalidateCache()

	return nil
}

func (o *Object) GetOptions() *entry.ObjectOptions {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	return o.options
}

func (o *Object) SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error) {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	options, err := modifyFn(o.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := o.db.GetObjectsDB().UpdateObjectOptions(o.ctx, o.GetID(), options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	o.options = options
	o.invalidateCache()

	return options, nil
}

func (o *Object) GetEffectiveOptions() *entry.ObjectOptions {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	if o.effectiveOptions == nil {
		effectiveOptions, err := merge.Auto(o.options, o.objectType.GetOptions())
		if err != nil {
			o.log.Error(
				errors.WithMessagef(
					err, "Object: GetEffectiveOptions: failed to merge object effective options: %s", o.GetID(),
				),
			)
			return nil
		}

		o.effectiveOptions = effectiveOptions
	}

	return o.effectiveOptions
}

func (o *Object) InvalidateCache() {
	o.Mu.Lock()
	defer o.Mu.Unlock()

	o.invalidateCache()
}

func (o *Object) invalidateCache() {
	o.effectiveOptions = nil
}

func (o *Object) GetEntry() *entry.Object {
	o.Mu.RLock()
	defer o.Mu.RUnlock()

	entry := &entry.Object{
		ObjectID:  o.id,
		OwnerID:   o.ownerID,
		Options:   o.options,
		Transform: o.transform,
	}
	if o.objectType != nil {
		entry.ObjectTypeID = o.objectType.GetID()
	}
	if o.Parent != nil {
		entry.ParentID = o.Parent.GetID()
	}
	if o.asset2d != nil {
		entry.Asset2dID = utils.GetPTR(o.asset2d.GetID())
	}
	if o.asset3d != nil {
		entry.Asset3dID = utils.GetPTR(o.asset3d.GetID())
	}

	if o.objectType != nil && o.objectType.GetID() == umid.MustParse(seed.NodeObjectTypeID) {
		// TODO Think how to avoid this hack
		entry.ParentID = o.id // By convention Node has parentID of itself
	}

	return entry
}

func (o *Object) Run() error {
	o.numSendsQueued.Store(0)
	o.broadcastPipeline = make(chan *websocket.PreparedMessage, 100)

	go func() {
		defer func() {
			ns := o.numSendsQueued.Swap(chanIsClosed)
			for i := int64(0); i < ns; i++ {
				<-o.broadcastPipeline
			}
			close(o.broadcastPipeline)
		}()

		for {
			select {
			case message := <-o.broadcastPipeline:
				o.numSendsQueued.Add(-1)
				if message == nil {
					return
				}

				o.performBroadcast(message)
			case <-o.ctx.Done():
				o.Stop()
			}
		}
	}()

	return nil
}

func (o *Object) Stop() error {
	ns := o.numSendsQueued.Add(1)
	if ns >= 0 {
		o.broadcastPipeline <- nil
	}
	return nil
}

func (o *Object) Update(recursive bool) error {
	o.UpdateSpawnMessage()

	if o.GetEnabled() {
		world := o.GetWorld()
		if world != nil {
			world.Send(o.spawnMsg.Load(), true)
			o.SendAllAutoAttributes(
				func(msg *websocket.PreparedMessage) error {
					return world.Send(msg, true)
				}, false,
			)
		}
	}

	if !recursive {
		return nil
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		if err := child.Update(true); err != nil {
			return errors.WithMessagef(err, "failed to update child: %s", child.GetID())
		}
	}

	return nil
}

func (o *Object) LoadFromEntry(entry *entry.Object, recursive bool) error {
	o.log.Debugf("Loading object: %s...", entry.ObjectID)

	if entry.ObjectID != o.GetID() {
		return errors.Errorf("object ids mismatch: %s != %s", entry.ObjectID, o.GetID())
	}

	group, ctx := errgroup.WithContext(o.ctx)
	group.Go(o.GetObjectAttributes().Load)
	group.Go(
		func() error {
			if err := o.load(entry); err != nil {
				return errors.WithMessage(err, "failed to load data")
			}

			if !recursive {
				return nil
			}

			entries, err := o.db.GetObjectsDB().GetObjectsByParentID(ctx, o.GetID())
			if err != nil {
				return errors.WithMessagef(err, "failed to get objects by parent umid: %s", o.GetID())
			}

			for _, childEntry := range entries {
				child, err := o.CreateObject(childEntry.ObjectID)
				if err != nil {
					return errors.WithMessagef(err, "failed to create new object: %s", childEntry.ObjectID)
				}
				if err := child.LoadFromEntry(childEntry, true); err != nil {
					return errors.WithMessagef(err, "failed to load object from entry: %s", childEntry.ObjectID)
				}
			}

			return nil
		},
	)
	return group.Wait()
}

func (o *Object) Save() error {
	return o.saveObjects(
		map[umid.UMID]universe.Object{
			o.GetID(): o,
		},
	)
}

func (o *Object) saveObjects(objects map[umid.UMID]universe.Object) error {
	if len(objects) < 1 {
		return nil
	}

	objList := make([]universe.Object, 0, len(objects))
	entries := make([]*entry.Object, 0, len(objects))
	for _, object := range objects {
		objList = append(objList, object)
		entries = append(entries, object.GetEntry())
	}

	// saving objects
	if err := o.db.GetObjectsDB().UpsertObjects(o.ctx, entries); err != nil {
		return errors.WithMessage(err, "failed to upsert objects")
	}

	// saving objects attributes
	if err := generic.NewButcher(objList).HandleItems(
		o.ctx,
		int(o.CFG.Postgres.MAXCONNS), // modify batchSize when database consumption while saving will be changed
		func(object universe.Object) error {
			if err := object.GetObjectAttributes().Save(); err != nil {
				return errors.WithMessagef(err, "failed to save object attributes: %s", object.GetID())
			}
			return nil
		},
	); err != nil {
		return errors.WithMessage(err, "failed to save objects attributes")
	}

	// saving children
	for _, object := range objects {
		if err := o.saveObjects(object.GetObjects(false)); err != nil {
			return errors.WithMessagef(err, "failed to save children: %s", object.GetID())
		}
	}

	return nil
}

func (o *Object) load(entry *entry.Object) error {
	node := universe.GetNode()

	if err := o.SetOwnerID(entry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner umid: %s", entry.OwnerID)
	}
	if _, err := o.SetOptions(modify.MergeWith(entry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	objectType, ok := node.GetObjectTypes().GetObjectType(entry.ObjectTypeID)
	if !ok {
		return errors.Errorf("failed to get object type: %s", entry.ObjectTypeID)
	}
	if err := o.SetObjectType(objectType, false); err != nil {
		return errors.WithMessagef(err, "failed to set object type: %s", entry.ObjectTypeID)
	}

	if entry.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*entry.Asset2dID)
		if !ok {
			return errors.Errorf("failed to get asset 2d: %s", entry.Asset2dID)
		}
		if err := o.SetAsset2D(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", entry.Asset2dID)
		}
	}

	if entry.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*entry.Asset3dID)
		if !ok {
			return errors.Errorf("failed to get asset 3d: %s", entry.Asset3dID)
		}
		if err := o.SetAsset3D(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", entry.Asset3dID)
		}
	}

	if err := o.SetTransform(entry.Transform, false); err != nil {
		return errors.WithMessage(err, "failed to set transform")
	}

	o.createdAt = entry.CreatedAt
	o.updatedAt = entry.UpdatedAt

	return nil
}

func (o *Object) UpdateSpawnMessage() error {
	world := o.GetWorld()
	if world == nil {
		return errors.Errorf("world is empty")
	}

	parentID := umid.Nil
	parent := o.GetParent()
	if parent != nil {
		parentID = parent.GetID()
	}

	asset3dID := umid.Nil
	asset3d := o.GetAsset3D()
	objectType := o.GetObjectType()
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
			utils.MapDecode(asset3dMeta, &metaData)
			assetFormat = dto.Asset3dType(metaData.Type)
		}
	}

	effectiveOptions := o.GetEffectiveOptions()

	// TODO: discuss is it ok to rely on "ReactSpaceVisibleType"?
	var visible bool
	if effectiveOptions.Visible != nil && *effectiveOptions.Visible == entry.UI2DObjectVisibleType {
		visible = true
	}

	mData := make([]posbus.ObjectDefinition, 1)
	mData[0] = posbus.ObjectDefinition{ID: o.GetID(), ParentID: parentID, AssetType: asset3dID, AssetFormat: assetFormat, Name: o.GetName(), IsEditable: *utils.GetFromAny(
		effectiveOptions.Editable, utils.GetPTR(true),
	),
		ShowOnMiniMap: *utils.GetFromAny(effectiveOptions.Minimap, &visible), Transform: *o.GetActualTransform()}
	msg := posbus.WSMessage(&posbus.AddObjects{Objects: mData})
	o.spawnMsg.Store(msg)

	return nil
}

func (o *Object) GetSpawnMessage() *websocket.PreparedMessage {
	return o.spawnMsg.Load()
}

func (o *Object) SendSpawnMessage(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	sendFn(o.spawnMsg.Load())
	//time.Sleep(time.Millisecond * 100)
	if !recursive {
		return
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		child.SendSpawnMessage(sendFn, true)
	}

}

func (o *Object) SendAllAutoAttributes(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	msg := o.dataMsg.Load()
	if msg != nil {
		sendFn(msg)
	}

	if !recursive {
		return
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		child.SendAllAutoAttributes(sendFn, recursive)
	}
}

// QUESTION: why this method is never called?
func (o *Object) SendAttributes(sendFn func(*websocket.PreparedMessage), recursive bool) {
	o.attributesMsg.Mu.RLock()
	for _, g := range o.attributesMsg.Data {
		for _, a := range g.Data {
			sendFn(a)
		}
	}
	o.attributesMsg.Mu.RUnlock()

	sendFn(o.spawnMsg.Load())

	if !recursive {
		return
	}

	o.Children.Mu.RLock()
	defer o.Children.Mu.RUnlock()

	for _, child := range o.Children.Data {
		child.SendAttributes(sendFn, true)
	}
}

// QUESTION: why this method is never called?
func (o *Object) SetAttributesMsg(kind, name string, msg *websocket.PreparedMessage) {
	m, ok := o.attributesMsg.Load(kind)
	if !ok {
		m = generic.NewSyncMap[string, *websocket.PreparedMessage](0)
		o.attributesMsg.Store(kind, m)
	}
	m.Store(name, msg)
}

func (o *Object) LockUIObject(user universe.User, state uint32) bool {
	if state == 1 {
		return o.lockedBy.CompareAndSwap(umid.Nil, user.GetID())
	} else {
		return o.lockedBy.CompareAndSwap(user.GetID(), umid.Nil)
	}
}

func (o *Object) IsLockedByUser(user universe.User) bool {
	return o.lockedBy.Load() == user.GetID()
}

func (o *Object) GetLockUserID() umid.UMID {
	return o.lockedBy.Load().(umid.UMID)
}

func (o *Object) GetCreatedAt() time.Time {
	return o.createdAt
}

func (o *Object) GetUpdatedAt() time.Time {
	return o.updatedAt
}
