package helper

import (
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetGuestUserTypeID() (uuid.UUID, error) {
	userTypeAttributeValue, ok := universe.GetNode().GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Node.GuestUserType.Name),
	)
	if !ok || userTypeAttributeValue == nil {
		err := errors.New("failed to get user type attribute value")
		return uuid.Nil, err
	}

	guestUserType := utils.GetFromAnyMap(*userTypeAttributeValue, universe.Attributes.Node.GuestUserType.Key, "")
	guestUserTypeID, err := uuid.Parse(guestUserType)
	if err != nil {
		err := errors.New("failed to parse guest user type id")
		return uuid.Nil, err
	}

	return guestUserTypeID, err
}
