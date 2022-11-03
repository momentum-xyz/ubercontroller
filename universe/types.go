package universe

import (
	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
)

type WorldsFilterPredicateFn func(worldID uuid.UUID, world World) bool
type SpacesFilterPredicateFn func(spaceID uuid.UUID, space Space) bool
type Assets2dFilterPredicateFn func(asset2dID uuid.UUID, asset2d Asset2d) bool
type Assets3dFilterPredicateFn func(asset3dID uuid.UUID, asset3d Asset3d) bool
type PluginsFilterPredicateFn func(pluginID uuid.UUID, plugin Plugin) bool
type AttributeTypesFilterPredicateFn func(attributeTypeID entry.AttributeTypeID, attributeType AttributeType) bool
type SpaceTypesFilterPredicateFn func(spaceTypeID uuid.UUID, spaceType SpaceType) bool
type UserTypesFilterPredicateFn func(userTypeID uuid.UUID, userType UserType) bool
