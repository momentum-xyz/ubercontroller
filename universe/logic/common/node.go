package common

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetGuestUserTypeID() (umid.UMID, error) {
	guestUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.GuestUserType.Name),
	)
	if !ok || guestUserTypeValue == nil {
		err := errors.New("failed to get guest user type attribute value")
		return umid.Nil, err
	}

	guestUserType := utils.GetFromAnyMap(
		*guestUserTypeValue, universe.ReservedAttributes.Node.GuestUserType.Key, "",
	)
	guestUserTypeID, err := umid.Parse(guestUserType)
	if err != nil {
		err := errors.New("failed to parse guest user type umid")
		return umid.Nil, err
	}

	return guestUserTypeID, err
}

func GetNormalUserTypeID() (umid.UMID, error) {
	normUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.NormalUserType.Name),
	)
	if !ok || normUserTypeValue == nil {
		return umid.Nil, errors.Errorf("failed to get normal user type attribute value")
	}

	normUserType := utils.GetFromAnyMap(*normUserTypeValue, universe.ReservedAttributes.Node.NormalUserType.Key, "")
	normUserTypeID, err := umid.Parse(normUserType)
	if err != nil {
		return umid.Nil, errors.Errorf("failed to parse normal user type umid")
	}

	return normUserTypeID, nil
}

func GetPortalObjectTypeID() (umid.UMID, error) {
	portalObjectTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.PortalObjectType.Name),
	)
	if !ok || portalObjectTypeValue == nil {
		return umid.Nil, errors.Errorf("failed to get portal object type attribute value")
	}

	portalObjectType := utils.GetFromAnyMap(
		*portalObjectTypeValue, universe.ReservedAttributes.Node.PortalObjectType.Key, "",
	)
	portalObjectTypeID, err := umid.Parse(portalObjectType)
	if err != nil {
		return umid.Nil, errors.Errorf("failed to parse portal object type umid")
	}

	return portalObjectTypeID, nil
}
