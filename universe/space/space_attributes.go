package space

import (
	"github.com/gorilla/websocket"
	"github.com/momentum-xyz/ubercontroller/pkg/message"
	"github.com/momentum-xyz/ubercontroller/utils/merge"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

func (s *Space) GetSpaceAttributeValue(attributeID entry.AttributeID) (*entry.AttributeValue, bool) {
	payload, ok := s.spaceAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}
	return payload.Value, true
}

func (s *Space) GetSpaceAttributeOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	payload, ok := s.spaceAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}
	return payload.Options, false
}

func (s *Space) GetSpaceAttributeEffectiveOptions(attributeID entry.AttributeID) (*entry.AttributeOptions, bool) {
	attr, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(attributeID))
	if !ok {
		return nil, false
	}
	payload, ok := s.spaceAttributes.Load(attributeID)
	if !ok {
		return nil, false
	}

	effectiveOptions, err := merge.Auto(payload.Options, attr.GetOptions())
	if err != nil {
		s.log.Error(
			errors.WithMessagef(
				err, "Space: GetSpaceAttributeEffectiveOptions: failed to merge space attribute effective options: %s: %+v",
				s.id, attributeID,
			),
		)
		return nil, false
	}

	return effectiveOptions, true
}

func (s *Space) SetSpaceAttributeValue(
	attributeID entry.AttributeID, modifyFn modify.Fn[entry.AttributeValue], updateDB bool,
) error {
	s.spaceAttributes.Mu.Lock()
	defer s.spaceAttributes.Mu.Unlock()

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("space attribute not found")
	}

	value, err := modifyFn(payload.Value)
	if err != nil {
		return errors.WithMessage(err, "failed to modify value")
	}
	payload.Value = value

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeValue(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.id), payload.Value,
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

	payload, ok := s.spaceAttributes.Data[attributeID]
	if !ok {
		return errors.Errorf("space attribute not found")
	}

	options, err := modifyFn(payload.Options)
	if err != nil {
		return errors.WithMessage(err, "failed to modify options")
	}
	payload.Options = options

	if updateDB {
		if err := s.db.SpaceAttributesUpdateSpaceAttributeOptions(
			s.ctx, entry.NewSpaceAttributeID(attributeID, s.id), payload.Options,
		); err != nil {
			return errors.WithMessage(err, "failed to update db")
		}
	}

	return nil
}

func (s *Space) CheckIfRendered(instance *entry.SpaceAttribute) {

	attr, ok := universe.GetNode().GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(instance.AttributeID))
	if !ok {
		return
	}
	var opts entry.AttributeOptions
	//if instance.Options != nil {
	//	opts = *instance.Options
	//} else {
	//	opts = entry.AttributeOptions{}
	//}

	if attr.GetOptions() != nil {
		opts = *attr.GetOptions()
	} else {
		return
	}
	//utils.MergePTRs(&opts, attr.GetOptions())
	//if opts == nil {
	//	return
	//}
	if v, ok := opts["render_type"]; ok && v.(string) == "texture" {
		if c, ok := (map[string]any)(*instance.Value)["render_hash"]; ok {
			s.renderTextureAttr[attr.GetName()] = c.(string)
		}
	}
	s.textMsg.Store(message.GetBuilder().SetObjectTextures(s.id, s.renderTextureAttr))
}

func (s *Space) SendTextures(f func(*websocket.PreparedMessage) error, recursive bool) {

	f(s.textMsg.Load())
	if recursive {
		s.Children.Mu.RLock()
		for _, space := range s.Children.Data {
			space.SendTextures(f, true)
		}
		s.Children.Mu.RUnlock()
	}

}

func (s *Space) loadSpaceAttributes() error {
	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.id)
	if err != nil {
		return errors.WithMessage(err, "failed to get space attributes")
	}

	node := universe.GetNode()
	for _, instance := range entries {
		if _, ok := node.GetAttributeTypes().GetAttributeType(entry.AttributeTypeID(instance.AttributeID)); ok {
			s.CheckIfRendered(instance)
			s.spaceAttributes.Store(
				instance.AttributeID,
				entry.NewAttributePayload(instance.Value, instance.Options),
			)
		}
	}

	return nil
}

//func (s *Space) SetSpaceAttribute() {
//	s.spaceUserAttributes.SetAttributeInstance(
//		NewUserAttributeIndex(instance.PluginID, instance.Name, instance.UserID), instance.Value,
//		instance.Options, attr,
//	)
//}
