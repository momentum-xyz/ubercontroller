package attribute_instances

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type AttributeInstance struct {
	attribute        universe.Attribute
	options          *entry.AttributeOptions
	effectiveOptions *entry.AttributeOptions
	value            *entry.AttributeValue
}

func (a AttributeInstance) GetOptions() *entry.AttributeOptions {
	return a.options

}

func (a AttributeInstance) SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) error {
	panic("implement me")
	return nil
}

func (a AttributeInstance) GetValue() *entry.AttributeValue {
	return a.value
}

func (a AttributeInstance) SetValue(modifyFn modify.Fn[string], updateDB bool) error {
	panic("implement me")
	return nil

}

func (a AttributeInstance) GetEntry() *entry.Attribute {
	panic("implement me")
	return nil

}

func (a AttributeInstance) LoadFromEntry(entry *entry.Attribute) error {
	panic("implement me")
	return nil
}

func (a AttributeInstance) GetEffectiveOptions() *entry.AttributeOptions {
	if a.effectiveOptions == nil {
		a.effectiveOptions = utils.MergePTRs(a.options, a.attribute.GetOptions())
	}
	return a.effectiveOptions
}

//func (s *Space) LoadAttributes() error {
//	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceID(s.ctx, s.GetID())
//	if err != nil {
//		return errors.WithMessage(err, "failed to load attribute data")
//	}
//
//	for _, e := range entries {
//		s.LoadAttributeFromEntry(e)
//	}
//
//	return nil
//}
//
//func (s *Space) LoadUserAttributes() error {
//	entries, err := s.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceID(s.ctx, s.GetID())
//	if err != nil {
//		return errors.WithMessage(err, "failed to load user attribute data")
//	}
//
//	for _, e := range entries {
//		s.LoadUserAttributeFromEntry(e)
//	}
//
//	return nil
//}
//
//func (s *Space) LoadAttributeFromEntry(e *entry.SpaceAttribute) {
//	s.attributes.Store(
//		AttributeIndex{PluginId: e.PluginID, Name: e.Name}, Attribute{Value: e.Value, Options: e.Options},
//	)
//}
//
//func (s *Space) LoadUserAttributeFromEntry(e *entry.SpaceUserAttribute) {
//	s.userAttributes.Store(
//		UserAttributeIndex{PluginId: e.PluginID, Name: e.Name, UserId: e.UserID},
//		Attribute{Value: e.Value, Options: e.Options},
//	)
//}
