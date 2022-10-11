package space

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
)

//var _ universe.AttributeIndexType = (*AttributeIndex)(nil)
//var _ universe.AttributeIndexType = (*UserAttributeIndex)(nil)

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	node := universe.GetNode()
	for _, instance := range entries {
		attr, ok := node.GetAttributes().GetAttribute(
			entry.AttributeID{
				PluginID: instance.PluginID,
				Name:     instance.Name,
			},
		)
		if ok {
			s.spaceAttributes.SetAttributeInstance(
				types.NewSpaceAttributeIndex(instance.PluginID, instance.Name),
				attr, instance.Value, instance.Options,
			)
		}
	}

	return nil
}

func (s *Space) loadSpaceUserAttributes() error {
	entries, err := s.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to load space user attributes")
	}

	node := universe.GetNode()
	for _, instance := range entries {
		attr, ok := node.GetAttributes().GetAttribute(
			entry.AttributeID{
				PluginID: instance.PluginID,
				Name:     instance.Name,
			},
		)
		if ok {
			ai := s.spaceUserAttributes.SetAttributeInstance(
				types.NewSpaceUserAttributeIndex(instance.PluginID, instance.Name, instance.UserID),
				attr, instance.Value, instance.Options,
			)
			opt := ai.GetEffectiveOptions()

			if v, ok := (*opt)["render_type"]; ok && v == "texture" {

			}
		}
	}

	return nil
}

func (s *Space) SetAttributesMsg(kind, name string, msg *websocket.PreparedMessage) {
	m, ok := s.attributesMsg.Load(kind)
	if !ok {
		m = generic.NewSyncMap[string, *websocket.PreparedMessage]()
		s.attributesMsg.Store(kind, m)
	}
	m.Store(name, msg)
}

//func (s *Space) SetSpaceAttribute() {
//	s.spaceUserAttributes.SetAttributeInstance(
//		NewUserAttributeIndex(instance.PluginID, instance.Name, instance.UserID), instance.Value,
//		instance.Options, attr,
//	)
//}
