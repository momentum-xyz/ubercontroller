package node

import (
	"encoding/hex"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

// @Summary Get current user's stakes list
// @Schemes
// @Description Return stakes list
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/me/stakes [get]
func (n *Node) apiGetMyStakes(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetMyStakes: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	userEntry, err := n.db.GetUsersDB().GetUserByID(c, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersGetMe: user not found")
		api.AbortRequest(c, http.StatusNotFound, "user_not_found", err, n.log)
		return
	}
	_ = userEntry

	//wallets, err := n.db.GetUsersDB().GetUserWalletsByUserID(c, umid.MustParse("f4c90bda-34c9-4e6f-9d8e-328164c6a019"))
	wallets, err := n.db.GetUsersDB().GetUserWalletsByUserID(c, userID)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiUsersGetMe: wallets not found for given user_id:%s", userID)
		api.AbortRequest(c, http.StatusNotFound, "wallets_not_found", err, n.log)
		return
	}

	result := make([]*map[string]any, 0)

	for _, w := range wallets {
		r, err := n.db.GetStakesDB().GetStakes(c, HexToAddress(*w))
		if err != nil {
			err := errors.WithMessagef(err, "Node: apiUsersGetMe: can not get stakes for wallet:%s", *w)
			api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
			return
		}
		result = append(result, r...)
	}

	c.JSON(http.StatusOK, result)
}

func HexToAddress(s string) []byte {
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		panic(err)
	}
	return b
}
