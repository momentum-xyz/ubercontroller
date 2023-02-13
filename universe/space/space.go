package space

import (
	"context"
	"fmt"
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

var _ universe.Space = (*Space)(nil)

type Space struct {
	id       uuid.UUID
	world    universe.World
	ctx      context.Context
	log      *zap.SugaredLogger
	db       database.DB
	enabled  atomic.Bool
	Users    *generic.SyncMap[uuid.UUID, universe.User]
	Children *generic.SyncMap[uuid.UUID, universe.Space]
	//Mu               sync.RWMutex
	Mu               deadlock.RWMutex
	ownerID          uuid.UUID
	position         *cmath.SpacePosition
	options          *entry.SpaceOptions
	Parent           universe.Space
	asset2d          universe.Asset2d
	asset3d          universe.Asset3d
	spaceType        universe.SpaceType
	effectiveOptions *entry.SpaceOptions
	spaceAttributes  *spaceAttributes // WARNING: the Space is sharing the same mutex ("Mu") with it

	spawnMsg          atomic.Pointer[websocket.PreparedMessage]
	attributesMsg     *generic.SyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]]
	renderTextureMap  *generic.SyncMap[string, string]
	renderStringMap   *generic.SyncMap[string, string]
	textureMsg        atomic.Pointer[websocket.PreparedMessage]
	stringMsg         atomic.Pointer[websocket.PreparedMessage]
	actualPosition    atomic.Pointer[cmath.SpacePosition]
	broadcastPipeline chan *websocket.PreparedMessage
	messageAccept     atomic.Bool
	numSendsQueued    atomic.Int64

	lockedBy atomic.Value

	// TODO: replace theta with full calculation of orientation, once Unity is read
	theta float64
}

func NewSpace(id uuid.UUID, db database.DB, world universe.World) *Space {
	space := &Space{
		id:               id,
		db:               db,
		Users:            generic.NewSyncMap[uuid.UUID, universe.User](0),
		Children:         generic.NewSyncMap[uuid.UUID, universe.Space](0),
		attributesMsg:    generic.NewSyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]](0),
		renderTextureMap: generic.NewSyncMap[string, string](0),
		renderStringMap:  generic.NewSyncMap[string, string](0),
		world:            world,
	}
	space.spaceAttributes = newSpaceAttributes(space)

	return space
}

func (s *Space) GetID() uuid.UUID {
	return s.id
}

func (s *Space) GetEnabled() bool {
	return s.enabled.Load()
}

func (s *Space) SetEnabled(enabled bool) {
	s.enabled.Store(enabled)
}

func (s *Space) GetName() string {
	name := "unknown"
	value, ok := s.GetSpaceAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Space.Name.Name),
	)
	if !ok || value == nil {
		return name
	}
	return utils.GetFromAnyMap(*value, universe.ReservedAttributes.Space.Name.Key, name)
}

func (s *Space) SetName(name string, updateDB bool) error {
	if _, err := s.GetSpaceAttributes().Upsert(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Space.Name.Name),
		modify.MergeWith(
			entry.NewAttributePayload(
				&entry.AttributeValue{
					universe.ReservedAttributes.Space.Name.Key: name,
				},
				nil,
			),
		), updateDB,
	); err != nil {
		return errors.WithMessage(err, "failed to upsert space attribute")
	}
	return nil
}

func (s *Space) GetSpaceAttributes() universe.SpaceAttributes {
	return s.spaceAttributes
}

func (s *Space) Initialize(ctx context.Context) error {
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

func (s *Space) GetWorld() universe.World {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.world
}

func (s *Space) GetParent() universe.Space {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.Parent
}

func (s *Space) SetParent(parent universe.Space, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if parent == s {
		return errors.Errorf("space can't be a parent of itself")
	} else if parent != nil && parent.GetWorld().GetID() != s.world.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.world.GetID())
	}

	if updateDB {
		if parent == nil {
			return errors.Errorf("parent is nil")
		}
		if err := s.db.GetSpacesDB().UpdateSpaceParentID(s.ctx, s.GetID(), parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.Parent = parent

	return nil
}

func (s *Space) GetOwnerID() uuid.UUID {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.ownerID
}

func (s *Space) SetOwnerID(ownerID uuid.UUID, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		if err := s.db.GetSpacesDB().UpdateSpaceOwnerID(s.ctx, s.GetID(), ownerID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.ownerID = ownerID

	return nil
}

func (s *Space) GetAsset2D() universe.Asset2d {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.asset2d
}

func (s *Space) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		var asset2dID *uuid.UUID
		if asset2d != nil {
			asset2dID = utils.GetPTR(asset2d.GetID())
		}
		if err := s.db.GetSpacesDB().UpdateSpaceAsset2dID(s.ctx, s.GetID(), asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d

	return nil
}

func (s *Space) GetAsset3D() universe.Asset3d {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.asset3d
}

func (s *Space) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		var asset3dID *uuid.UUID
		if asset3d != nil {
			asset3dID = utils.GetPTR(asset3d.GetID())
		}
		if err := s.db.GetSpacesDB().UpdateSpaceAsset3dID(s.ctx, s.GetID(), asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset3d = asset3d

	return nil
}

func (s *Space) GetSpaceType() universe.SpaceType {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.spaceType
}

func (s *Space) SetSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	if spaceType == nil {
		return errors.Errorf("space type is nil")
	}

	s.Mu.Lock()
	defer s.Mu.Unlock()

	if updateDB {
		if err := s.db.GetSpacesDB().UpdateSpaceSpaceTypeID(s.ctx, s.GetID(), spaceType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceType = spaceType
	s.dropCache()

	return nil
}

func (s *Space) GetOptions() *entry.SpaceOptions {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	return s.options
}

func (s *Space) SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) (*entry.SpaceOptions, error) {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	options, err := modifyFn(s.options)
	if err != nil {
		return nil, errors.WithMessage(err, "failed to modify options")
	}

	if updateDB {
		if err := s.db.GetSpacesDB().UpdateSpaceOptions(s.ctx, s.GetID(), options); err != nil {
			return nil, errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options
	s.dropCache()

	return options, nil
}

func (s *Space) GetEffectiveOptions() *entry.SpaceOptions {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if s.effectiveOptions == nil {
		effectiveOptions, err := merge.Auto(s.options, s.spaceType.GetOptions())
		if err != nil {
			s.log.Error(
				errors.WithMessagef(
					err, "Space: GetEffectiveOptions: failed to merge space effective options: %s", s.GetID(),
				),
			)
			return nil
		}

		s.effectiveOptions = effectiveOptions
	}

	return s.effectiveOptions
}

func (s *Space) DropCache() {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	s.dropCache()
}

func (s *Space) dropCache() {
	s.effectiveOptions = nil
}

func (s *Space) GetEntry() *entry.Space {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	entry := &entry.Space{
		SpaceID:  s.id,
		OwnerID:  &s.ownerID,
		Options:  s.options,
		Position: s.position,
	}
	if s.spaceType != nil {
		entry.SpaceTypeID = utils.GetPTR(s.spaceType.GetID())
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

func (s *Space) Run() error {
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

func (s *Space) Stop() error {
	ns := s.numSendsQueued.Add(1)
	if ns >= 0 {
		s.broadcastPipeline <- nil
	}
	return nil
}

func (s *Space) Update(recursive bool) error {
	s.UpdateSpawnMessage()

	if s.GetEnabled() {
		world := s.GetWorld()
		if world != nil {
			world.Send(s.spawnMsg.Load(), true)
			s.SendAllAutoAttributes(
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

func (s *Space) LoadFromEntry(entry *entry.Space, recursive bool) error {
	s.log.Debugf("Loading space %s...", entry.SpaceID)

	if entry.SpaceID != s.GetID() {
		return errors.Errorf("space ids mismatch: %s != %s", entry.SpaceID, s.GetID())
	}

	group, _ := errgroup.WithContext(s.ctx)
	group.Go(
		func() error {
			if err := s.loadSpaceAttributes(); err != nil {
				return errors.WithMessage(err, "failed to load space attributes")
			}
			return nil
		},
	)
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

			entries, err := s.db.GetSpacesDB().GetSpacesByParentID(s.ctx, s.GetID())
			if err != nil {
				return errors.WithMessagef(err, "failed to get spaces by parent id: %s", s.GetID())
			}

			for i := range entries {
				child, err := s.CreateSpace(entries[i].SpaceID)
				if err != nil {
					return errors.WithMessagef(err, "failed to create new space: %s", entries[i].SpaceID)
				}
				if err := child.LoadFromEntry(entries[i], recursive); err != nil {
					return errors.WithMessagef(err, "failed to load space from entry: %s", entries[i].SpaceID)
				}
			}

			return nil
		},
	)
	return group.Wait()
}

func (s *Space) loadSelfData(spaceEntry *entry.Space) error {
	if err := s.SetOwnerID(*spaceEntry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner id: %s", spaceEntry.OwnerID)
	}
	if _, err := s.SetOptions(modify.MergeWith(spaceEntry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}
	return nil
}

func (s *Space) loadDependencies(entry *entry.Space) error {
	node := universe.GetNode()

	spaceType, ok := node.GetSpaceTypes().GetSpaceType(*entry.SpaceTypeID)
	if !ok {
		return errors.Errorf("failed to get space type: %s", entry.SpaceTypeID)
	}
	if err := s.SetSpaceType(spaceType, false); err != nil {
		return errors.WithMessagef(err, "failed to set space type: %s", entry.SpaceTypeID)
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

func (s *Space) UpdateSpawnMessage() error {
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
	spaceType := s.GetSpaceType()
	assetFormat := dto.AddressableAssetType
	if asset3d == nil && spaceType != nil {
		asset3d = spaceType.GetAsset3d()
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

	effectiveOptions := s.GetEffectiveOptions()

	// TODO: discuss is it ok to rely on "ReactSpaceVisibleType"?
	var visible bool
	if effectiveOptions.Visible != nil {
		if *effectiveOptions.Visible == entry.ReactSpaceVisibleType || *effectiveOptions.Visible == entry.ReactUnitySpaceVisibleType {
			visible = true
		}
	}

	t := *utils.GetFromAny(effectiveOptions.Minimap, &visible)
	if s.Parent != nil && s.Parent.GetID() == uuid.MustParse("e7860a2c-0cb0-4dbb-8e95-74fc8070b565") {
		fmt.Printf("ttt___ %+v %+v %+v %+v \n", s.GetName(), s.GetID(), s.spaceType.GetName(), t)
	}

	msg := message.GetBuilder().MsgObjectDefinition(
		message.ObjectDefinition{
			ObjectID:         s.GetID(),
			ParentID:         parentID,
			AssetType:        asset3dID,
			AssetFormat:      assetFormat,
			Name:             s.GetName(),
			Position:         *s.GetActualPosition(),
			Editable:         *utils.GetFromAny(effectiveOptions.Editable, utils.GetPTR(true)),
			TetheredToParent: true,
			Minimap:          *utils.GetFromAny(effectiveOptions.Minimap, &visible),
			InfoUI:           *utils.GetFromAny(effectiveOptions.InfoUIID, utils.GetPTR(uuid.Nil)),
		},
	)
	s.spawnMsg.Store(msg)

	return nil
}

func (s *Space) GetSpawnMessage() *websocket.PreparedMessage {
	return s.spawnMsg.Load()
}

func (s *Space) SendSpawnMessage(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	sendFn(s.spawnMsg.Load())
	//time.Sleep(time.Millisecond * 100)
	if !recursive {
		return
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, space := range s.Children.Data {
		space.SendSpawnMessage(sendFn, recursive)
	}

}

func (s *Space) SendAllAutoAttributes(sendFn func(*websocket.PreparedMessage) error, recursive bool) {
	msg := s.stringMsg.Load()
	if msg != nil {
		sendFn(msg)
	}

	msg = s.textureMsg.Load()
	if msg != nil {
		sendFn(msg)
	}

	if !recursive {
		return
	}

	s.Children.Mu.RLock()
	defer s.Children.Mu.RUnlock()

	for _, space := range s.Children.Data {
		space.SendAllAutoAttributes(sendFn, recursive)
	}
}

// QUESTION: why this method is never called?
func (s *Space) SendAttributes(sendFn func(*websocket.PreparedMessage), recursive bool) {
	s.attributesMsg.Mu.RLock()
	for _, g := range s.attributesMsg.Data {
		for _, a := range g.Data {
			sendFn(a)
		}
	}
	s.attributesMsg.Mu.RUnlock()

	sendFn(s.spawnMsg.Load())
	if recursive {
		s.Children.Mu.RLock()
		defer s.Children.Mu.RUnlock()

		for _, space := range s.Children.Data {
			space.SendAttributes(sendFn, recursive)
		}
	}
}

// QUESTION: why this method is never called?
func (s *Space) SetAttributesMsg(kind, name string, msg *websocket.PreparedMessage) {
	m, ok := s.attributesMsg.Load(kind)
	if !ok {
		m = generic.NewSyncMap[string, *websocket.PreparedMessage](0)
		s.attributesMsg.Store(kind, m)
	}
	m.Store(name, msg)
}

func (s *Space) LockUnityObject(user universe.User, state uint32) bool {
	if state == 1 {
		return s.lockedBy.CompareAndSwap(uuid.Nil, user.GetID())
	} else {
		return s.lockedBy.CompareAndSwap(user.GetID(), uuid.Nil)
	}
}
