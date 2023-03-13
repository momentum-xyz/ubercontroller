package common

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetGuestUserTypeID() (uuid.UUID, error) {
	guestUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.GuestUserType.Name),
	)
	if !ok || guestUserTypeValue == nil {
		err := errors.New("failed to get guest user type attribute value")
		return uuid.Nil, err
	}

	guestUserType := utils.GetFromAnyMap(
		*guestUserTypeValue, universe.ReservedAttributes.Node.GuestUserType.Key, "",
	)
	guestUserTypeID, err := uuid.Parse(guestUserType)
	if err != nil {
		err := errors.New("failed to parse guest user type id")
		return uuid.Nil, err
	}

	return guestUserTypeID, err
}

func GetNormalUserTypeID() (uuid.UUID, error) {
	normUserTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.NormalUserType.Name),
	)
	if !ok || normUserTypeValue == nil {
		return uuid.Nil, errors.Errorf("failed to get normal user type attribute value")
	}

	normUserType := utils.GetFromAnyMap(*normUserTypeValue, universe.ReservedAttributes.Node.NormalUserType.Key, "")
	normUserTypeID, err := uuid.Parse(normUserType)
	if err != nil {
		return uuid.Nil, errors.Errorf("failed to parse normal user type id")
	}

	return normUserTypeID, nil
}

func GetPortalObjectTypeID() (uuid.UUID, error) {
	portalObjectTypeValue, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.PortalObjectType.Name),
	)
	if !ok || portalObjectTypeValue == nil {
		return uuid.Nil, errors.Errorf("failed to get portal object type attribute value")
	}

	portalObjectType := utils.GetFromAnyMap(
		*portalObjectTypeValue, universe.ReservedAttributes.Node.PortalObjectType.Key, "",
	)
	portalObjectTypeID, err := uuid.Parse(portalObjectType)
	if err != nil {
		return uuid.Nil, errors.Errorf("failed to parse portal object type id")
	}

	return portalObjectTypeID, nil
}
