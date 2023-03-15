package common

import (
	"github.com/momentum-xyz/ubercontroller/utils/mid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetGuestUserTypeID() (mid.ID, error) {
	guestUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.GuestUserType.Name),
	)
	if !ok || guestUserTypeValue == nil {
		err := errors.New("failed to get guest user type attribute value")
		return mid.Nil, err
	}

	guestUserType := utils.GetFromAnyMap(
		*guestUserTypeValue, universe.ReservedAttributes.Node.GuestUserType.Key, "",
	)
	guestUserTypeID, err := mid.Parse(guestUserType)
	if err != nil {
		err := errors.New("failed to parse guest user type mid")
		return mid.Nil, err
	}

	return guestUserTypeID, err
}

func GetNormalUserTypeID() (mid.ID, error) {
	normUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.NormalUserType.Name),
	)
	if !ok || normUserTypeValue == nil {
		return mid.Nil, errors.Errorf("failed to get normal user type attribute value")
	}

	normUserType := utils.GetFromAnyMap(*normUserTypeValue, universe.ReservedAttributes.Node.NormalUserType.Key, "")
	normUserTypeID, err := mid.Parse(normUserType)
	if err != nil {
		return mid.Nil, errors.Errorf("failed to parse normal user type mid")
	}

	return normUserTypeID, nil
}

func GetPortalObjectTypeID() (mid.ID, error) {
	portalObjectTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.PortalObjectType.Name),
	)
	if !ok || portalObjectTypeValue == nil {
		return mid.Nil, errors.Errorf("failed to get portal object type attribute value")
	}

	portalObjectType := utils.GetFromAnyMap(
		*portalObjectTypeValue, universe.ReservedAttributes.Node.PortalObjectType.Key, "",
	)
	portalObjectTypeID, err := mid.Parse(portalObjectType)
	if err != nil {
		return mid.Nil, errors.Errorf("failed to parse portal object type mid")
	}

	return portalObjectTypeID, nil
}
