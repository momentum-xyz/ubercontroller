package attribute_data

import (
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type AttributeData struct {
	id               entry.AttributeID
	options          *entry.AttributeOptions
	effectiveOptions *entry.AttributeOptions
	value            *entry.AttributeValue
	attribute        universe.Attribute
}

func (a *AttributeData) GetOptions() *entry.AttributeOptions {

}
func (a *AttributeData) SetOptions(modifyFn modify.Fn[entry.AttributeOptions], updateDB bool) error {

}

func (a *AttributeData) GetValue() *string {

}
func (a *AttributeData) SetValue(modifyFn modify.Fn[string], updateDB bool) error {

}

func (a *AttributeData) GetEntry() *entry.Attribute {

}
func (a *AttributeData) LoadFromEntry(entry *entry.Attribute) error {

}

func (a *AttributeData) GetEffectiveOptions() *entry.AttributeOptions {
	if a.effectiveOptions == nil {
		a.effectiveOptions = utils.MergePTRs(a.options, a.attribute.GetOptions())
	}
	return a.effectiveOptions
}

//func (s *Space) LoadAttributes() error {
//	entries, err := s.db.SpaceAttributesGetSpaceAttributesBySpaceId(s.ctx, s.GetID())
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
//	entries, err := s.db.SpaceUserAttributesGetSpaceUserAttributesBySpaceId(s.ctx, s.GetID())
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
