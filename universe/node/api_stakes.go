package node

import (
	"encoding/hex"
	"math/big"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
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

	wallets, err := n.db.GetUsersDB().GetUserWalletsByUserID(c, userID)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiUsersGetMe: wallets not found for given user_id:%s", userID)
		api.AbortRequest(c, http.StatusNotFound, "wallets_not_found", err, n.log)
		return
	}

	result := make([]*map[string]any, 0)
	for _, w := range wallets {
		r, err := n.db.GetStakesDB().GetStakes(c, utils.HexToAddress(*w))
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
		walletAddresses = append(walletAddresses, utils.HexToAddress(*w))
	}

	result, err := n.db.GetStakesDB().GetWalletsInfo(c, walletAddresses)
	if err != nil {
		err := errors.WithMessagef(err, "Node: apiGetMyWallets: can not get wallets for user:%s", userID)
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, result)
}

// @Summary Add pending stake transaction
// @Schemes
// @Description Add pending transaction
// @Tags users
// @Accept json
// @Produce json
// @Param body body node.apiAddPendingStakeTransaction.Body true "body params"
// @Success 200 {object} node.apiAddPendingStakeTransaction.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/users/me/stakes [post]
func (n *Node) apiAddPendingStakeTransaction(c *gin.Context) {
	type Body struct {
		TransactionID string    `json:"transaction_id" binding:"required"`
		OdysseyID     umid.UMID `json:"odyssey_id" binding:"required"`
		Wallet        string    `json:"wallet" binding:"required"`
		Comment       string    `json:"comment" binding:"required"`
		Amount        string    `json:"amount" binding:"required"`
		Kind          string    `json:"kind" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiAddPendingStakeTransaction: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	transactionID, err := hexToAddress(inBody.TransactionID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiAddPendingStakeTransaction: failed to parse transaction_id to byte array")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	wallet, err := hexToAddress(inBody.Wallet)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiAddPendingStakeTransaction: failed to parse wallet to byte array")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	big := big.NewInt(0)
	amount, ok := big.SetString(inBody.Amount, 10)
	if !ok {
		err := errors.New("Node: apiAddPendingStakeTransaction: failed to parse amount")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	err = n.db.GetStakesDB().InsertIntoPendingStakes(c, transactionID,
		inBody.OdysseyID, wallet, umid.MustParse("ccccaaaa-1111-2222-3333-222222222222"), amount, inBody.Comment, 0)
	if err != nil {
		err := errors.New("Node: apiAddPendingStakeTransaction: failed to insert into pending stakes")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	type Out struct {
		Success bool `json:"success"`
	}
	out := Out{
		Success: true,
	}

	c.JSON(http.StatusOK, out)
}

func hexToAddress(s string) ([]byte, error) {
	b, err := hex.DecodeString(s[2:])
	if err != nil {
		return nil, err
	}
	return b, nil
}
