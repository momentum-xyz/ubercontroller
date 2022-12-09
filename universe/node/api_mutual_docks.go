package node

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"net/http"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Create mutual docks for teleport if users staked to each other
// @Schemes
// @Description After users has been made mutual staking this EP will add mutual teleport docks to user's Odysseys
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} node.apiUsersCreateMutualDocks.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/mutual-docks [post]
func (n *Node) apiUsersCreateMutualDocks(c *gin.Context) {
	type Body struct {
		WalletA string `json:"walletA" binding:"required"`
		WalletB string `json:"walletB" binding:"required"`
	}
	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userA, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletA)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userA == nil {
		m := "User UUID not found for wallet:" + inBody.WalletA
		api.AbortRequest(c, http.StatusNotFound, "user_A_not_found", errors.New(m), n.log)
		return
	}

	userB, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletB)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userB == nil {
		m := "User UUID not found for wallet:" + inBody.WalletB
		api.AbortRequest(c, http.StatusNotFound, "user_B_not_found", errors.New(m), n.log)
		return
	}

	worldA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err = errors.New("Node: apiUsersCreateMutualDocks: failed to GetSpace for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if worldA == nil {
		err = errors.New("space not found for userID=:" + userA.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_A_not_found", err, n.log)
		return
	}

	worldB, ok := n.GetWorlds().GetWorld(userB.UserID)
	if !ok {
		err = errors.New("Node: apiUsersCreateMutualDocks: failed to GetSpace for userB")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if worldB == nil {
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

	bulbAID, err := n.addDockingBulb(worldA, dockStationType.GetID(), userA, userB)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to addDockingBulb for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	bulbBID, err := n.addDockingBulb(worldB, dockStationType.GetID(), userB, userA)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to addDockingBulb for userB")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	userSpaces := []*entry.UserSpace{
		{
			SpaceID: worldA.GetID(),
			UserID:  userB.UserID,
			Value:   map[string]any{"role": "admin"},
		},
		{
			SpaceID: worldB.GetID(),
			UserID:  userA.UserID,
			Value:   map[string]any{"role": "admin"},
		},
	}

	err = n.db.UserSpacesUpsertUserSpaces(n.ctx, userSpaces)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to UserSpacesUpsertUserSpaces")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	type Out struct {
		Status          string     `json:"status"`
		UserA           *uuid.UUID `json:"userA"`
		UserB           *uuid.UUID `json:"userB"`
		SpaceA          string     `json:"spaceA"`
		SpaceB          string     `json:"spaceB"`
		DockStationType string     `json:"dockStationType"`
		BulbA           string     `json:"bulbA"`
		BulbB           string     `json:"bulbB"`
	}
	out := Out{
		Status:          "ok",
		UserA:           &userA.UserID,
		UserB:           &userB.UserID,
		SpaceA:          worldA.GetID().String(),
		SpaceB:          worldB.GetID().String(),
		DockStationType: dockStationType.GetID().String(),
		BulbA:           bulbAID.String(),
		BulbB:           bulbBID.String(),
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Remove mutual docks
// @Schemes
// @Description After any of linked users has been unstake to another this EP removes docking bulbs and world admin rights for both
// @Tags users
// @Accept json
// @Produce json
// @Success 202 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/mutual-docks [delete]
func (n *Node) apiUsersRemoveMutualDocks(c *gin.Context) {
	type Body struct {
		WalletA string `json:"walletA" binding:"required"`
		WalletB string `json:"walletB" binding:"required"`
	}
	inBody := Body{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userA, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletA)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userA == nil {
		m := "User UUID not found for wallet:" + inBody.WalletA
		api.AbortRequest(c, http.StatusNotFound, "user_A_not_found", errors.New(m), n.log)
		return
	}

	userB, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletB)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to UsersGetUserByWallet")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_server_error", err, n.log)
		return
	}
	if userB == nil {
		m := "User UUID not found for wallet:" + inBody.WalletB
		api.AbortRequest(c, http.StatusNotFound, "user_B_not_found", errors.New(m), n.log)
		return
	}

	worldA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err = errors.New("Node: apiUsersRemoveMutualDocks: failed to GetSpace for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if worldA == nil {
		err = errors.New("space not found for userID=:" + userA.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_A_not_found", err, n.log)
		return
	}

	worldB, ok := n.GetWorlds().GetWorld(userB.UserID)
	if !ok {
		err = errors.New("Node: apiUsersRemoveMutualDocks: failed to GetSpace for userB")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if worldB == nil {
		err = errors.New("space not found for userID=:" + userB.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_B_not_found", err, n.log)
		return
	}

	userSpaces := []*entry.UserSpace{
		{
			SpaceID: worldA.GetID(),
			UserID:  userB.UserID,
		},
		{
			SpaceID: worldB.GetID(),
			UserID:  userA.UserID,
		},
	}

	for i := range userSpaces {
		if err := n.db.UserSpaceRemoveUserSpace(n.ctx, userSpaces[i]); err != nil {
			err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to remove userSpace")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}
	}

	bulbsA := n.findDockingBulbsByTargetWorldID(worldA, userB.UserID)
	bulbsB := n.findDockingBulbsByTargetWorldID(worldB, userA.UserID)
	bulbs := make(map[uuid.UUID]universe.Space, len(bulbsA)+len(bulbsB))
	for _, bulb := range bulbs {
		parent := bulb.GetParent()

		if _, err := parent.RemoveSpace(bulb, false, true); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: apiUsersRemoveMutualDocks: failed to remove bulb: %s", bulb.GetID(),
				),
			)
		}

		if err := parent.UpdateChildrenPosition(true); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: apiUsersRemoveMutualDocks: failed to update children position: %s", parent.GetID(),
				),
			)
		}

		if err := bulb.Stop(); err != nil {
			n.log.Error(
				errors.WithMessagef(
					err, "Node: apiUsersRemoveMutualDocks: failed to stop bulb: %s", bulb.GetID(),
				),
			)
		}

		bulb.SetEnabled(false)
	}

	c.JSON(http.StatusAccepted, nil)
}

func (n *Node) addDockingBulb(world universe.World, dockStationTypeID uuid.UUID, fromUser *entry.User, toUser *entry.User) (uuid.UUID, error) {
	bulbs := n.findDockingBulbsByTargetWorldID(world, toUser.UserID)
	if len(bulbs) > 0 {
		for spaceID, _ := range bulbs {
			return spaceID, nil
		}
	}

	dockStation := n.getDockStation(world, dockStationTypeID)

	bulbSpaceID := uuid.New()
	err := n.createAndAddBulbSpace(bulbSpaceID, dockStation, *fromUser, *toUser)
	if err != nil {
		return bulbSpaceID, errors.WithMessage(err, "Node: addDockingBulb: failed to createAndAddBulbSpace")
	}

	return bulbSpaceID, nil
}

func (n *Node) createAndAddBulbSpace(spaceID uuid.UUID, dock universe.Space, fromUser entry.User, toUser entry.User) error {
	toUserName := ""
	if toUser.Profile != nil {
		if toUser.Profile.Name != nil {
			toUserName = *toUser.Profile.Name
		}
	}

	bulbSpaceTypeName := "Docking bulb"
	spaceType := n.getSpaceTypeByName(bulbSpaceTypeName)
	if spaceType == nil {
		return errors.Errorf("failed to get space type: %s", bulbSpaceTypeName)
	}

	spaceTemplate := SpaceTemplate{
		SpaceID:     &spaceID,
		SpaceName:   toUserName,
		SpaceTypeID: spaceType.GetID(),
		ParentID:    dock.GetID(),
		OwnerID:     &fromUser.UserID,
		SpaceAttributes: []*Attribute{
			{
				AttributeID: entry.NewAttributeID(universe.GetSystemPluginID(), "teleport"),
				AttributePayload: entry.AttributePayload{
					Value: &entry.AttributeValue{
						"DestinationWorldID": toUser.UserID.String(),
					},
				},
			},
		},
	}

	if _, err := n.addSpaceFromTemplate(&spaceTemplate); err != nil {
		return errors.WithMessagef(err, "failed to add space from template: %+v", spaceTemplate)
	}

	return nil
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

func (n *Node) getDockStation(world universe.World, dockStationTypeID uuid.UUID) universe.Space {
	allSpaces := world.GetAllSpaces()

	for _, space := range allSpaces {
		spaceType := space.GetSpaceType()
		if spaceType == nil {
			continue
		}
		if spaceType.GetID() == dockStationTypeID {
			return space
		}
	}

	return nil
}

func (n *Node) findDockingBulbsByTargetWorldID(world universe.World, targetWorldID uuid.UUID) map[uuid.UUID]universe.Space {
	targetWorld := targetWorldID.String()
	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), "teleport")

	predicateFn := func(spaceID uuid.UUID, space universe.Space) bool {
		value, ok := space.GetSpaceAttributeValue(attributeID)
		if !ok || value == nil {
			return false
		}

		worldID := utils.GetFromAnyMap(*value, "DestinationWorldID", "")
		if worldID != targetWorld {
			return false
		}

		return true
	}

	return world.FilterAllSpaces(predicateFn)
}
