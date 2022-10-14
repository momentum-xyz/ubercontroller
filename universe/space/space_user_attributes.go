package space

import (
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
)

func (s *Space) loadSpaceUserAttributes() error {
	entries, err := s.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to load space user attributes")
	}

	node := universe.GetNode()
	for _, instance := range entries {
		if _, ok := node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(instance.AttributeID)); ok {
			s.spaceUserAttributes.Store(
				entry.NewUserAttributeID(instance.AttributeID, instance.UserID),
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
		}
	}

	return nil
}
