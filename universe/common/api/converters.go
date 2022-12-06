package api

import (
	"time"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func ToUserDTO(userEntry *entry.User) (*dto.User, error) {
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

	profileEntry := userEntry.Profile

	outUser := dto.User{
		ID:     userEntry.UserID.String(),
		Wallet: wallet,
		Profile: dto.Profile{
			Bio:         profileEntry.Bio,
			Location:    profileEntry.Location,
			AvatarHash:  profileEntry.AvatarHash,
			ProfileLink: profileEntry.ProfileLink,
		},
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if profileEntry != nil {
		if profileEntry.Name != nil {
			outUser.Name = *profileEntry.Name
		}
	}

	return &outUser, nil
}
