package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Edit user profile
// @Schemes
// @Description Edits a user profile
// @Tags profile
// @Accept json
// @Produce json
// @Param body body node.apiProfileUpdate.Body true "body params"
// @Success 200 {object} dto.User
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/profile [patch]
func (n *Node) apiProfileUpdate(c *gin.Context) {
	type Body struct {
		Name    *string `json:"name"`
		Profile *struct {
			Bio         *string `json:"bio"`
			ProfileLink *string `json:"profileLink"`
			Location    *string `json:"location"`
			AvatarHash  *string `json:"avatarHash"`
		} `json:"profile"`
	}

	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user id from context")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	userProfile, err := n.db.GetUsersDB().GetUserProfileByUserID(c, userID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiProfileUpdate: failed to get user profile by user id")
		api.AbortRequest(c, http.StatusNotFound, "not_found", err, n.log)
		return
	}

	nameChanged := inBody.Name != nil && *userProfile.Name != *inBody.Name
	avatarChanged := inBody.Profile != nil && inBody.Profile.AvatarHash != nil && *userProfile.AvatarHash != *inBody.Profile.AvatarHash

	if inBody.Name != nil {
		// TODO: check name unique
		userProfile.Name = inBody.Name
	}

	inProfile := inBody.Profile
	if inProfile != nil {
		if inProfile.Bio != nil {
			userProfile.Bio = inProfile.Bio
		}
		if inProfile.Location != nil {
			userProfile.Location = inProfile.Location
		}
		if inProfile.AvatarHash != nil {
			userProfile.AvatarHash = inProfile.AvatarHash
		}
		if inProfile.ProfileLink != nil {
			userProfile.ProfileLink = inProfile.ProfileLink
		}
	}

	userProfile.OnBoarded = utils.GetPTR(true)

	if nameChanged || avatarChanged {
		// Need update NFT metadata
		a := userProfile.Name != inBody.Name
		fmt.Println(a)
		if err := n.updateNFTMetadata(c, userID, userProfile.Name, userProfile.AvatarHash); err != nil {
			api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update_user_profile", err, n.log)
			return
		}
	}

	if err := n.db.GetUsersDB().UpdateUserProfile(c, userID, userProfile); err != nil {
		err = errors.WithMessage(err, "failed to update user profile")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_update_user_profile", err, n.log)
		return
	}

	n.apiUsersGetMe(c)
}

func (n *Node) updateNFTMetadata(c *gin.Context, userID uuid.UUID, userName *string, avatarHash *string) error {
	wallet, err := n.db.GetUsersDB().GetUserWalletByUserID(c, userID)
	if err != nil {
		return errors.WithMessage(err, "failed to get user wallet by userID")
	}

	fmt.Println(wallet)

	//walletMeta, err := n.getWalletMetadata(*wallet)
	//if err != nil {
	//	return errors.WithMessagef(err, "failed to get wallet metadata: %s", *wallet)
	//}

	output, err := exec.Command("node", "./nodejs/check-nft/update-nft.js", *wallet).Output()
	if err != nil {
		return errors.WithMessage(err, "failed to execute node script update-nft.js")
	}

	fmt.Println(output)
	fmt.Println(string(output))

	var nodeJSOut NodeJSOut
	if err := json.Unmarshal(output, &nodeJSOut); err != nil {
		return errors.WithMessage(err, "failed to unmarshal output")
	}

	fmt.Println(nodeJSOut)

	return nil
}
