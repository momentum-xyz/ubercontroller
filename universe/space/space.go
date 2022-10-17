package space

import (
	"context"
	"fmt"
	"github.com/sasha-s/go-deadlock"
	"sync/atomic"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.Space = (*Space)(nil)

type Space struct {
	id       uuid.UUID
	world    universe.World
	ctx      context.Context
	log      *zap.SugaredLogger
	db       database.DB
	Users    *generic.SyncMap[uuid.UUID, universe.User]
	Children *generic.SyncMap[uuid.UUID, universe.Space]
	//mu               sync.RWMutex
	mu               deadlock.RWMutex
	ownerID          uuid.UUID
	position         *cmath.Vec3
	options          *entry.SpaceOptions
	parent           universe.Space
	asset2d          universe.Asset2d
	asset3d          universe.Asset3d
	spaceType        universe.SpaceType
	effectiveOptions *entry.SpaceOptions

	spaceAttributes     *generic.SyncMap[entry.AttributeID, *entry.AttributePayload]
	spaceUserAttributes *generic.SyncMap[entry.UserAttributeID, *entry.AttributePayload]

	spawnMsg          atomic.Pointer[websocket.PreparedMessage]
	attributesMsg     *generic.SyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]]
	actualPosition    atomic.Pointer[cmath.Vec3]
	broadcastPipeline chan *websocket.PreparedMessage
	messageAccept     atomic.Bool
	numSendsQueued    atomic.Int64

	// TODO: replace theta with full calculation of orientation, once Unity is read
	theta float64
}

func NewSpace(id uuid.UUID, db database.DB, world universe.World) *Space {
	return &Space{
		id:                  id,
		db:                  db,
		Users:               generic.NewSyncMap[uuid.UUID, universe.User](),
		Children:            generic.NewSyncMap[uuid.UUID, universe.Space](),
		spaceAttributes:     generic.NewSyncMap[entry.AttributeID, *entry.AttributePayload](),
		spaceUserAttributes: generic.NewSyncMap[entry.UserAttributeID, *entry.AttributePayload](),
		attributesMsg:       generic.NewSyncMap[string, *generic.SyncMap[string, *websocket.PreparedMessage]](),
		world:               world,
	}
}

func (s *Space) GetID() uuid.UUID {
	return s.id
}

// todo: implement this via spaceAttributes
func (s *Space) GetName() string {
	return "unknown"
}

func (s *Space) Initialize(ctx context.Context) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if log == nil {
		return errors.Errorf("failed to get logger from context: %T", ctx.Value(types.LoggerContextKey))
	}

	s.ctx = ctx
	s.log = log
	s.numSendsQueued.Store(chanIsClosed)
	s.actualPosition.Store(new(cmath.Vec3))
	go s.Run()

	return nil
}

func (s *Space) GetWorld() universe.World {
	return s.world
}

func (s *Space) GetParent() universe.Space {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.parent
}

func (s *Space) SetParent(parent universe.Space, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if parent != nil && parent.GetWorld().GetID() != s.world.GetID() {
		return errors.Errorf("worlds mismatch: %s != %s", parent.GetWorld().GetID(), s.world.GetID())
	}

	if updateDB {
		if parent == nil {
			return errors.Errorf("parent is nil")
		}
		if err := s.db.SpacesUpdateSpaceParentID(s.ctx, s.id, parent.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.parent = parent

	return nil
}

func (s *Space) GetOwnerID() uuid.UUID {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.ownerID
}

func (s *Space) SetOwnerID(ownerID uuid.UUID, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpaceOwnerID(s.ctx, s.id, ownerID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.ownerID = ownerID

	return nil
}

func (s *Space) GetAsset2D() universe.Asset2d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset2d
}

func (s *Space) SetAsset2D(asset2d universe.Asset2d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		var asset2dID *uuid.UUID
		if asset2d != nil {
			asset2dID = utils.GetPTR(asset2d.GetID())
		}
		if err := s.db.SpacesUpdateSpaceAsset2dID(s.ctx, s.id, asset2dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset2d = asset2d

	return nil
}

func (s *Space) GetAsset3D() universe.Asset3d {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.asset3d
}

func (s *Space) SetAsset3D(asset3d universe.Asset3d, updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		var asset3dID *uuid.UUID
		if asset3d != nil {
			asset3dID = utils.GetPTR(asset3d.GetID())
		}
		if err := s.db.SpacesUpdateSpaceAsset3dID(s.ctx, s.id, asset3dID); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.asset3d = asset3d

	return nil
}

func (s *Space) GetSpaceType() universe.SpaceType {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.spaceType
}

func (s *Space) SetSpaceType(spaceType universe.SpaceType, updateDB bool) error {
	if spaceType == nil {
		return errors.Errorf("space type is nil")
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if updateDB {
		if err := s.db.SpacesUpdateSpaceSpaceTypeID(s.ctx, s.id, spaceType.GetID()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.spaceType = spaceType

	return nil
}

func (s *Space) GetOptions() *entry.SpaceOptions {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.options
}

func (s *Space) SetOptions(modifyFn modify.Fn[entry.SpaceOptions], updateDB bool) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	options := modifyFn(s.options)

	if updateDB {
		if err := s.db.SpacesUpdateSpaceOptions(s.ctx, s.id, options); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	s.options = options
	s.effectiveOptions = nil

	return nil
}

func (s *Space) GetEffectiveOptions() *entry.SpaceOptions {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.effectiveOptions == nil {
		s.effectiveOptions = utils.MergePTRs(s.options, s.spaceType.GetOptions())
	}
	return s.effectiveOptions
}

func (s *Space) GetEntry() *entry.Space {
	s.mu.RLock()
	defer s.mu.RUnlock()

	entry := &entry.Space{
		SpaceID:  s.id,
		OwnerID:  &s.ownerID,
		Options:  s.options,
		Position: s.position,
	}
	if s.spaceType != nil {
		entry.SpaceTypeID = utils.GetPTR(s.spaceType.GetID())
	}
	if s.parent != nil {
		entry.ParentID = utils.GetPTR(s.parent.GetID())
	}
	if s.asset2d != nil {
		entry.Asset2dID = utils.GetPTR(s.asset2d.GetID())
	}
	if s.asset3d != nil {
		entry.Asset3dID = utils.GetPTR(s.asset3d.GetID())
	}

	return entry
}

func (s *Space) Update(recursive bool) error {
	s.mu.Lock()
	s.effectiveOptions = nil
	s.mu.Unlock()

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

	if entry.SpaceID != s.id {
		return errors.Errorf("space ids mismatch: %s != %s", entry.SpaceID, s.id)
	}

	if err := s.loadSelfData(entry); err != nil {
		return errors.WithMessage(err, "failed to load self data")
	}
	if err := s.loadDependencies(entry); err != nil {
		return errors.WithMessage(err, "failed to load dependencies")
	}

	if err := s.SetPosition(entry.Position, false); err != nil {
		return errors.WithMessage(err, "failed to set position")
	}

	s.UpdateSpawnMessage()

	if !recursive {
		return nil
	}

	entries, err := s.db.SpacesGetSpacesByParentID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessagef(err, "failed to get spaces by parent id: %s", s.id)
	}

	for i := range entries {
		space, err := s.CreateSpace(entries[i].SpaceID)
		if err != nil {
			return errors.WithMessagef(err, "failed to create new space: %s", entries[i].SpaceID)
		}
		if err := space.LoadFromEntry(entries[i], recursive); err != nil {
			return errors.WithMessagef(err, "failed to load space from entry: %s", entries[i].SpaceID)
		}
		s.Children.Store(entries[i].SpaceID, space)
	}

	return nil
}

func (s *Space) loadSelfData(spaceEntry *entry.Space) error {
	if err := s.SetOwnerID(*spaceEntry.OwnerID, false); err != nil {
		return errors.WithMessagef(err, "failed to set owner id: %s", spaceEntry.OwnerID)
	}

	if err := s.SetOptions(modify.ReplaceWith(spaceEntry.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	if err := s.loadSpaceAttributes(); err != nil {
		return errors.WithMessage(err, "failed to load space attributes")
	}
	if err := s.loadSpaceUserAttributes(); err != nil {
		return errors.WithMessage(err, "failed to load space_user attributes")
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

	if err := s.loadSpaceAttributes(); err != nil {
		return errors.WithMessage(err, "failed to load space attributes")
	}
	if err := s.loadSpaceUserAttributes(); err != nil {
		return errors.WithMessage(err, "failed to load space user attributes")
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

func (s *Space) UpdateSpawnMessage() {
	name := " "
	v, ok := s.GetSpaceAttributeValue(entry.NewAttributeID(universe.GetSystemPluginID(), "name"))
	if ok {
		// QUESTION: what to do here? maybe return an error?
		name = utils.GetFromAnyMap(*v, "name", "")
	} else {
		fmt.Println("smsg2", s.GetID())
	}

	opts := s.GetEffectiveOptions()
	parent := uuid.Nil
	if s.GetParent() != nil {
		parent = s.GetParent().GetID()
	}
	asset3d := uuid.Nil
	if s.asset3d != nil {
		asset3d = s.asset3d.GetID()
	}
	uuidNilPtr := utils.GetPTR(uuid.Nil)
	falsePtr := utils.GetPTR(false)
	msg := message.GetBuilder().MsgObjectDefinition(
		message.ObjectDefinition{
			ObjectID:         s.id,
			ParentID:         parent,
			AssetType:        asset3d,
			Name:             name,
			Position:         s.GetActualPosition(),
			TetheredToParent: true,
			Minimap:          *utils.GetFromAny(opts.Minimap, falsePtr),
			InfoUI:           *utils.GetFromAny(opts.InfoUIID, uuidNilPtr),
		},
	)
	s.spawnMsg.Store(msg)
	s.world.Broadcast(s.spawnMsg.Load(), false)
}

func (s *Space) GetSpawnMessage() *websocket.PreparedMessage {
	return s.spawnMsg.Load()
}

func (s *Space) SendSpawnMessage(f func(*websocket.PreparedMessage) error, recursive bool) {
	f(s.spawnMsg.Load())
	if recursive {
		s.Children.Mu.RLock()
		for _, space := range s.Children.Data {
			space.SendSpawnMessage(f, true)
		}
		s.Children.Mu.RUnlock()
	}

}

func (s *Space) SendAttributes(f func(*websocket.PreparedMessage), recursive bool) {
	s.attributesMsg.Mu.RLock()
	for _, g := range s.attributesMsg.Data {
		for _, a := range g.Data {
			f(a)
		}
	}
	s.attributesMsg.Mu.RUnlock()

	f(s.spawnMsg.Load())
	if recursive {
		s.Children.Mu.RLock()
		defer s.Children.Mu.RUnlock()

		for _, space := range s.Children.Data {
			space.SendAttributes(f, true)
		}
		s.Children.Mu.RUnlock()
	}
}

func (s *Space) SetAttributesMsg(kind, name string, msg *websocket.PreparedMessage) {
	m, ok := s.attributesMsg.Load(kind)
	if !ok {
		m = generic.NewSyncMap[string, *websocket.PreparedMessage]()
		s.attributesMsg.Store(kind, m)
	}
	m.Store(name, msg)
}
