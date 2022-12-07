package node

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe"
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

	attributes, err := n.db.UserAttributesGetUserAttributes(c)
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

	userA, err := n.db.UsersGetUserByWallet(context.Background(), inBody.WalletA)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userA == nil {
		m := "User UUID not found for wallet:" + inBody.WalletA
		api.AbortRequest(c, http.StatusNotFound, "user_A_not_found", errors.New(m), n.log)
		return
	}

	userB, err := n.db.UsersGetUserByWallet(context.Background(), inBody.WalletB)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userB == nil {
		m := "User UUID not found for wallet:" + inBody.WalletB
		api.AbortRequest(c, http.StatusNotFound, "user_B_not_found", errors.New(m), n.log)
		return
	}

	spaceA, err := n.db.SpacesGetSpaceByID(c, userA.UserID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to SpacesGetSpaceByID for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if spaceA == nil {
		err = errors.New("space not found for userID=:" + userA.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_A_not_found", err, n.log)
		return
	}

	spaceB, err := n.db.SpacesGetSpaceByID(c, userB.UserID)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to SpacesGetSpaceByID for userB")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if spaceB == nil {
		err = errors.New("space not found for userID=:" + userB.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_B_not_found", err, n.log)
		return
	}

	name := "Docking station"
	dockStationType := n.getSpaceTypeByName(name)
	if dockStationType == nil {
		err = errors.New("dockStationType not found for name=:" + name)
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}
	fmt.Println(dockStationType)

	type Out struct {
		Status          string     `json:"status"`
		UserA           *uuid.UUID `json:"userA"`
		UserB           *uuid.UUID `json:"userB"`
		SpaceA          *uuid.UUID `json:"spaceA"`
		SpaceB          *uuid.UUID `json:"spaceB"`
		DockStationType string     `json:"dockStationType"`
	}
	out := Out{
		Status:          "ok",
		UserA:           &userA.UserID,
		UserB:           &userB.UserID,
		SpaceA:          &spaceA.SpaceID,
		SpaceB:          &spaceB.SpaceID,
		DockStationType: dockStationType.GetID().String(),
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) getSpaceTypeByName(name string) universe.SpaceType {
	types := n.GetSpaceTypes().GetSpaceTypes()
	for _, v := range types {
		if v.GetName() == name {
			return v
		}
	}

	return nil
}

func (n *Node) getDockStation(dockStationType universe.SpaceType) (universe.Space, error) {
	return nil, nil
}
