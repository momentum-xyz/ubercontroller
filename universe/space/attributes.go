package space

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
)

//var _ universe.AttributeIndexType = (*AttributeIndex)(nil)
//var _ universe.AttributeIndexType = (*UserAttributeIndex)(nil)

type AttributeIndex struct {
	entry.AttributeID
}

func (a AttributeIndex) GetAttributeID() entry.AttributeID {
	return a.AttributeID
}

func NewAttributeIndex(pluginId uuid.UUID, name string) AttributeIndex {
	return AttributeIndex{entry.AttributeID{PluginID: pluginId, Name: name}}

}

type UserAttributeIndex struct {
	AttributeID entry.AttributeID
	UserId      uuid.UUID
}

func (u UserAttributeIndex) GetAttributeID() entry.AttributeID {
	return u.AttributeID
}

func NewUserAttributeIndex(pluginId uuid.UUID, name string, userId uuid.UUID) UserAttributeIndex {
	return UserAttributeIndex{AttributeID: entry.AttributeID{PluginID: pluginId, Name: name}, UserId: userId}
}

func (s *Space) loadSpaceAttributes() error {

	if entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceId(s.ctx, s.id); err != nil {
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
					NewAttributeIndex(instance.PluginID, instance.Name), instance.Value, instance.Options, attr,
				)
			}

		}

	}
	return nil
}

func (s *Space) loadSpaceUserAttributes() error {

	if entries, err := s.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceId(s.ctx, s.id); err != nil {
		node := universe.GetNode()
		for _, instance := range entries {
			attr, ok := node.GetAttributes().GetAttribute(
				entry.AttributeID{
					PluginID: instance.PluginID,
					Name:     instance.Name,
				},
			)
			if ok {
				ai := s.userSpaceAttributes.SetAttributeInstance(
					NewUserAttributeIndex(instance.PluginID, instance.Name, instance.UserID), instance.Value,
					instance.Options, attr,
				)
				opt := ai.GetEffectiveOptions()

				if v, ok := (*opt)["render_type"]; ok && v == "texture" {

				}
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
//	s.userSpaceAttributes.SetAttributeInstance(
//		NewUserAttributeIndex(instance.PluginID, instance.Name, instance.UserID), instance.Value,
//		instance.Options, attr,
//	)
//}
