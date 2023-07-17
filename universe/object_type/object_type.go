package object_type

import (
	"context"
	"sync"

	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

var _ universe.ObjectType = (*ObjectType)(nil)

type ObjectType struct {
	id           umid.UMID
	ctx          context.Context
	log          *zap.SugaredLogger
	db           database.DB
	mu           sync.RWMutex
	name         string
	categoryName string
	description  *string
	options      *entry.ObjectOptions
	asset2d      universe.Asset2d
	asset3d      universe.Asset3d
}

func NewObjectType(id umid.UMID, db database.DB) *ObjectType {
	return &ObjectType{
		id: id,
		db: db,
		options: &entry.ObjectOptions{
			AllowedChildren: []umid.UMID{},
			Minimap:         utils.GetPTR(true),
			Visible:         utils.GetPTR(entry.AllObjectVisibleType),
			Editable:        utils.GetPTR(true),
			Private:         utils.GetPTR(false),
		},
	}
}

func (ot *ObjectType) GetID() umid.UMID {
	return ot.id
}

func (ot *ObjectType) Initialize(ctx types.LoggerContext) error {
	ot.ctx = ctx
	ot.log = ctx.Logger()

	return nil
}

func (ot *ObjectType) GetName() string {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.name
}

func (ot *ObjectType) SetName(name string, updateDB bool) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if updateDB {
		if err := ot.db.GetObjectTypesDB().UpdateObjectTypeName(ot.ctx, ot.GetID(), name); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.name = name

	return nil
}

func (ot *ObjectType) GetCategoryName() string {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.categoryName
}

func (ot *ObjectType) SetCategoryName(categoryName string, updateDB bool) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if updateDB {
		if err := ot.db.GetObjectTypesDB().UpdateObjectTypeCategoryName(ot.ctx, ot.GetID(), categoryName); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.categoryName = categoryName

	return nil
}

func (ot *ObjectType) GetDescription() *string {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.description
}

func (ot *ObjectType) SetDescription(description *string, updateDB bool) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if updateDB {
		if err := ot.db.GetObjectTypesDB().UpdateObjectTypeDescription(ot.ctx, ot.GetID(), description); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.description = description

	return nil
}

func (ot *ObjectType) GetAsset2d() universe.Asset2d {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.asset2d
}

func (ot *ObjectType) SetAsset2d(asset2d universe.Asset2d, updateDB bool) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if updateDB {
		if err := ot.db.GetAssets2dDB().UpsertAsset(ot.ctx, asset2d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.asset2d = asset2d

	return nil
}

func (ot *ObjectType) GetAsset3d() universe.Asset3d {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.asset3d
}

func (ot *ObjectType) SetAsset3d(asset3d universe.Asset3d, updateDB bool) error {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	if updateDB {
		if err := ot.db.GetAssets3dDB().UpsertAsset(ot.ctx, asset3d.GetEntry()); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	ot.asset3d = asset3d

	return nil
}

func (ot *ObjectType) GetOptions() *entry.ObjectOptions {
	ot.mu.RLock()
	defer ot.mu.RUnlock()

	return ot.options
}

func (ot *ObjectType) SetOptions(modifyFn modify.Fn[entry.ObjectOptions], updateDB bool) (*entry.ObjectOptions, error) {
	options, err := func() (*entry.ObjectOptions, error) {
		ot.mu.Lock()
		defer ot.mu.Unlock()

		options, err := modifyFn(ot.options)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to modify options")
		}

		if updateDB {
			if err := ot.db.GetObjectTypesDB().UpdateObjectTypeOptions(ot.ctx, ot.GetID(), options); err != nil {
				return nil, errors.WithMessage(err, "failed to update db")
			}
		}

		ot.options = options

		return options, nil
	}()
	if err != nil {
		return nil, err
	}

	for _, object := range universe.GetNode().GetAllObjects() {
		objectType := object.GetObjectType()
		if objectType == nil {
			continue
		}
		if objectType.GetID() != ot.GetID() {
			continue
		}
		object.InvalidateCache()
	}

	return options, nil
}

func (ot *ObjectType) GetEntry() *entry.ObjectType {
	ot.mu.Lock()
	defer ot.mu.Unlock()

	entry := &entry.ObjectType{
		ObjectTypeID:   ot.id,
		ObjectTypeName: ot.name,
		CategoryName:   ot.categoryName,
		Description:    ot.description,
		Options:        ot.options,
	}
	if ot.asset2d != nil {
		entry.Asset2dID = utils.GetPTR(ot.asset2d.GetID())
	}
	if ot.asset3d != nil {
		entry.Asset3dID = utils.GetPTR(ot.asset3d.GetID())
	}

	return entry
}

func (ot *ObjectType) LoadFromEntry(row *entry.ObjectType) error {
	if row.ObjectTypeID != ot.GetID() {
		return errors.Errorf("object type ids mismatch: %s != %s", row.ObjectTypeID, ot.GetID())
	}

	if err := ot.SetName(row.ObjectTypeName, false); err != nil {
		return errors.WithMessage(err, "failed to set name")
	}
	if err := ot.SetCategoryName(row.CategoryName, false); err != nil {
		return errors.WithMessage(err, "failed to set category name")
	}
	if err := ot.SetDescription(row.Description, false); err != nil {
		return errors.WithMessage(err, "failed to set description")
	}
	if _, err := ot.SetOptions(modify.MergeWith(row.Options), false); err != nil {
		return errors.WithMessage(err, "failed to set options")
	}

	node := universe.GetNode()
	if row.Asset2dID != nil {
		asset2d, ok := node.GetAssets2d().GetAsset2d(*row.Asset2dID)
		if !ok {
			return errors.Errorf("asset 2d not found: %s", row.Asset2dID)
		}
		if err := ot.SetAsset2d(asset2d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 2d: %s", row.Asset2dID)
		}
	}

	if row.Asset3dID != nil {
		asset3d, ok := node.GetAssets3d().GetAsset3d(*row.Asset3dID)
		if !ok {
			return errors.Errorf("asset 3d not found: %s", row.Asset3dID)
		}
		if err := ot.SetAsset3d(asset3d, false); err != nil {
			return errors.WithMessagef(err, "failed to set asset 3d: %s", row.Asset3dID)
		}
	}

	return nil
}
