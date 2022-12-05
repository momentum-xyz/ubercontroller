package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Create mutual docks for teleport if users staked to each other
// @Schemes
// @Description After users has been made mutual staking this EP will add mutual teleport docks to user's Odysseys
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} node.apiUsersMutualDocks.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/mutual-docks [post]
func (n *Node) apiUsersMutualDocks(c *gin.Context) {
	type Body struct {
		WalletA string `json:"walletA" binding:"required"`
		WalletB string `json:"walletB" binding:"required"`
	}
	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	attributes, err := n.db.UserAttributesGetUserAttributes(context.Background())
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to get users attributes")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_users_attributes", err, n.log)
		return
	}

	type WalletAttribute struct {
	}

	fmt.Println(attributes)

	//userEntry, httpCode, err := n.apiGetOrCreateUserFromTokens(c, api.GetTokenFromRequest(c))
	//if err != nil {
	//	err = errors.WithMessage(err, "Node: apiUsersCheck: failed get or create user from tokens")
	//	api.AbortRequest(c, httpCode, "failed_to_get_or_create_user", err, n.log)
	//	return
	//}

	A, B := findUserIDs(inBody.WalletA, inBody.WalletB, attributes)

	if A == nil {
		m := "User UUID not found for wallet:" + inBody.WalletA
		api.AbortRequest(c, http.StatusNotFound, "user_A_not_found", errors.New(m), n.log)
		return
	}

	if B == nil {
		m := "User UUID not found for wallet:" + inBody.WalletB
		api.AbortRequest(c, http.StatusNotFound, "user_B_not_found", errors.New(m), n.log)
		return
	}

	type Out struct {
		Status string     `json:"status"`
		UserA  *uuid.UUID `json:"userA"`
		UserB  *uuid.UUID `json:"userB"`
	}
	out := Out{
		Status: "ok",
		UserA:  A,
		UserB:  B,
	}

	c.JSON(http.StatusOK, out)
}

func findUserIDs(walletA string, walletB string, allUserAttributes []*entry.UserAttribute) (*uuid.UUID, *uuid.UUID) {
	var userIDA *uuid.UUID
	var userIDB *uuid.UUID

	for _, a := range allUserAttributes {
		if a == nil {
			continue
		}
		if a.AttributePayload == nil {
			continue
		}
		if a.AttributePayload.Value != nil {
			value := *a.AttributePayload.Value
			wallets, ok := value["wallet"]
			if ok {
				list, ok := wallets.([]any)
				if ok {
					for _, w := range list {
						if w == walletA {
							userIDA = &a.UserID
						}
						if w == walletB {
							userIDB = &a.UserID
						}
					}
				}
			}
		}
	}

	return userIDA, userIDB
}
