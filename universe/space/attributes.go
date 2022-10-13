package space

import (
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (s *Space) GetSpaceAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := s.spaceAttributes.Load(entry.NewSpaceAttributeID(attributeID, s.id))
	if !ok {
		return nil, false
	}
	return payload.Value, true
}

func (s *Space) GetSpaceAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := s.spaceAttributes.Load(entry.NewSpaceAttributeID(attributeID, s.id))
	if !ok {
		return nil, false
	}
	return payload.Options, false
}

func (s *Space) GetSpaceAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attr, ok := universe.GetNode().GetAttributes().GetAttribute(attributeID)
	if !ok {
		return nil, false
	}
	payload, ok := s.spaceAttributes.Load(entry.NewSpaceAttributeID(attributeID, s.id))
	if !ok {
		return nil, false
	}
	return utils.MergePTRs(payload.Options, attr.GetOptions()), true
}

func (s *Space) SetSpaceAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[entry.NewSpaceAttributeID(attributeID, s.id)]
	if !ok {
		return errors.Errorf("space attribute not found")
	}

	payload.Value = modifyFn(payload.Value)

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeValue(
			s.ctx, attributeID.PluginID, attributeID.Name, s.id, payload.Value,
		); err != nil {
			return errors.WithMessage(err, "failed to udpate db")
		}
	}

	return nil
}

func (s *Space) SetSpaceAttributeOptions(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeOptions], updateDB bool,
) error {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[entry.NewSpaceAttributeID(attributeID, s.id)]
	if !ok {
		return errors.Errorf("space attribute not found")
	}

	payload.Options = modifyFn(payload.Options)

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeOptions(
			s.ctx, attributeID.PluginID, attributeID.Name, s.id, payload.Options,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	return nil
}

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	node := universe.GetNode()
	for _, instance := range entries {
		if _, ok := node.GetAttributes().GetAttribute(instance.AttributeID); ok {
			s.spaceAttributes.Store(
				instance.SpaceAttributeID,
				entry.NewAttributePayload(instance.Value, instance.Options),
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
		if _, ok := node.GetAttributes().GetAttribute(instance.AttributeID); ok {
			s.spaceUserAttributes.Store(
				entry.NewSpaceUserAttributeID(instance.AttributeID, instance.SpaceID, instance.UserID),
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
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
