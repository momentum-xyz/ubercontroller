package api

import (
	"time"

	"github.com/google/uuid"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func ToUserDTO(userEntry *entry.User, guestUserTypeID uuid.UUID, includeWallet bool) *dto.User {
	profileEntry := userEntry.Profile

	userDTO := dto.User{
		ID: userEntry.UserID.String(),
		Profile: dto.Profile{
			Bio:         profileEntry.Bio,
			Location:    profileEntry.Location,
			AvatarHash:  profileEntry.AvatarHash,
			ProfileLink: profileEntry.ProfileLink,
		},
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
		IsGuest:   false,
	}
	if userEntry.UserTypeID != nil {
		userDTO.UserTypeID = userEntry.UserTypeID.String()
		if *userEntry.UserTypeID == guestUserTypeID {
			userDTO.IsGuest = true
		}
	}
	if userEntry.UpdatedAt != nil {
		userDTO.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if profileEntry != nil {
		if profileEntry.Name != nil {
			userDTO.Name = *profileEntry.Name
		}
	}

	if includeWallet {
		var wallet *string

		walletValue, ok := universe.GetNode().GetUserAttributeValue(
			entry.NewUserAttributeID(
				entry.NewAttributeID(
					universe.GetKusamaPluginID(), universe.Attributes.Kusama.User.Wallet.Name,
				),
				userEntry.UserID,
			),
		)
		if ok && walletValue != nil {
			wallets := utils.GetFromAnyMap(*walletValue, universe.Attributes.Kusama.User.Wallet.Key, []any(nil))
			if len(wallets) > 0 {
				wallet = utils.GetPTR(utils.GetFromAny(wallets[0], ""))
			}
		}

		userDTO.Wallet = wallet
	}

	return &userDTO
}

func ToUserDTOs(userEntries []*entry.User, guestUserTypeID uuid.UUID, includeWallet bool) []*dto.User {
	userDTOs := make([]*dto.User, len(userEntries))
	for i := range userEntries {
		userDTOs[i] = ToUserDTO(userEntries[i], guestUserTypeID, includeWallet)
	}
	return userDTOs
}
