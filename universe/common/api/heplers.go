package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func GetUserDTOByID(db database.DB, ctx context.Context, userID uuid.UUID) (*dto.User, int, error) {
	userEntry, err := db.UsersGetUserByID(ctx, userID)
	if err != nil {
		return nil, http.StatusNotFound, errors.WithMessage(err, "failed to get user by id")
	}
	userProfileEntry := userEntry.Profile

	walletValue, ok := universe.GetNode().GetUserAttributeValue(
		entry.NewUserAttributeID(
			entry.NewAttributeID(
				universe.GetKusamaPluginID(), universe.Attributes.Kusama.User.Wallet.Name,
			),
			userID,
		),
	)
	if !ok || walletValue == nil {
		return nil, http.StatusInternalServerError, errors.Errorf("failed to get wallet user attriubte value")
	}

	wallets := utils.GetFromAnyMap(*walletValue, universe.Attributes.Kusama.User.Wallet.Key, []any(nil))
	if len(wallets) < 1 {
		return nil, http.StatusInternalServerError, errors.Errorf("invalid wallet user attribute value")
	}
	wallet := utils.GetFromAny(wallets[0], "")

	outUser := dto.User{
		ID:     userEntry.UserID.String(),
		Wallet: utils.GetPTR(wallet),
		Profile: dto.Profile{
			Bio:         userProfileEntry.Bio,
			Location:    userProfileEntry.Location,
			AvatarHash:  userProfileEntry.AvatarHash,
			ProfileLink: userProfileEntry.ProfileLink,
		},
		CreatedAt: userEntry.CreatedAt.Format(time.RFC3339),
	}
	if userEntry.UserTypeID != nil {
		outUser.UserTypeID = userEntry.UserTypeID.String()
	}
	if userEntry.UpdatedAt != nil {
		outUser.UpdatedAt = utils.GetPTR(userEntry.UpdatedAt.Format(time.RFC3339))
	}
	if userProfileEntry != nil {
		if userProfileEntry.Name != nil {
			outUser.Name = *userProfileEntry.Name
		}
	}

	return &outUser, 0, nil
}
