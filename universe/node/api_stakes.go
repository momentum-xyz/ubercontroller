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

// @Summary Get current user's wallets list
// @Schemes
// @Description Return wallets list
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/me/wallets [get]
func (n *Node) apiGetMyWallets(c *gin.Context) {
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetMyWallets: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	wallets, err := n.db.GetUsersDB().GetUserWalletsByUserID(c, userID)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiGetMyWallets: wallets not found for given user_id:%s", userID)
		api.AbortRequest(c, http.StatusNotFound, "wallets_not_found", err, n.log)
		return
	}

	walletAddresses := make([][]byte, 0)
	for _, w := range wallets {
		walletAddresses = append(walletAddresses, HexToAddress(*w))
	}

	result, err := n.db.GetStakesDB().GetWalletsInfo(c, walletAddresses)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiGetMyWallets: can not get wallets for user:%s", userID)
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
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
