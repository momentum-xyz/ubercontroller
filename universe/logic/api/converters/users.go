package converters

import (
	"github.com/momentum-xyz/ubercontroller/utils/umid"
	"time"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func ToUserDTO(userEntry *entry.User, guestUserTypeID umid.UMID, includeWallet bool) *dto.User {
	profileEntry := userEntry.Profile

	userDTO := dto.User{
		ID:         userEntry.UserID.String(),
		UserTypeID: userEntry.UserTypeID.String(),
		Profile: dto.Profile{
			Bio:         profileEntry.Bio,
			Location:    profileEntry.Location,
			AvatarHash:  profileEntry.AvatarHash,
			ProfileLink: profileEntry.ProfileLink,
		},
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
		UpdatedAt: userEntry.UpdatedAt.Format(time.RFC3339),
		IsGuest:   false,
	}
	if userEntry.UserTypeID == guestUserTypeID {
		userDTO.IsGuest = true
	}
	if profileEntry.Name != nil {
		userDTO.Name = *profileEntry.Name
	}

	if includeWallet {
		var wallet *string

		walletValue, ok := universe.GetNode().GetUserAttributes().GetValue(
			entry.NewUserAttributeID(
				entry.NewAttributeID(
					universe.GetKusamaPluginID(), universe.ReservedAttributes.Kusama.User.Wallet.Name,
				),
				userEntry.UserID,
			),
		)
		if ok && walletValue != nil {
			wallets := utils.GetFromAnyMap(*walletValue, universe.ReservedAttributes.Kusama.User.Wallet.Key, []any(nil))
			if len(wallets) > 0 {
				wallet = utils.GetPTR(utils.GetFromAny(wallets[0], ""))
			}
		}

		userDTO.Wallet = wallet
	}

	return &userDTO
}

func ToUserDTOs(userEntries []*entry.User, guestUserTypeID umid.UMID, includeWallet bool) []*dto.User {
	userDTOs := make([]*dto.User, len(userEntries))
	for i := range userEntries {
		userDTOs[i] = ToUserDTO(userEntries[i], guestUserTypeID, includeWallet)
	}
	return userDTOs
}
