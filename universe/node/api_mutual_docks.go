package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/common/helper"
	"github.com/momentum-xyz/ubercontroller/utils"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Create mutual docks
// @Schemes
// @Description Creates mutual worlds portals and worlds admin permissions
// @Tags users
// @Accept json
// @Produce json
// @Param body body node.apiUsersCreateMutualDocks.InBody true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/mutual-docks [post]
func (n *Node) apiUsersCreateMutualDocks(c *gin.Context) {
	type InBody struct {
		WalletA string `json:"walletA" binding:"required"`
		WalletB string `json:"walletB" binding:"required"`
	}
	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userA, err := n.db.GetUsersDB().GetUserByWallet(c, inBody.WalletA)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to get user A by wallet")
		api.AbortRequest(c, http.StatusNotFound, "user_a_not_found", err, n.log)
		return
	}

	userB, err := n.db.GetUsersDB().GetUserByWallet(c, inBody.WalletB)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to get user B by wallet")
		api.AbortRequest(c, http.StatusNotFound, "user_b_not_found", err, n.log)
		return
	}

	worldA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err := errors.Errorf("Node: apiUsersCreateMutualDocks: failed to get user A world: %s", userA.UserID)
		api.AbortRequest(c, http.StatusNotFound, "user_a_world_not_found", err, n.log)
		return
	}

	worldB, ok := n.GetWorlds().GetWorld(userB.UserID)
	if !ok {
		err := errors.Errorf("Node: apiUsersCreateMutualDocks: failed to get user B world: %s", userB.UserID)
		api.AbortRequest(c, http.StatusNotFound, "user_b_world_not_found", err, n.log)
		return
	}

	abPortalName := userB.UserID.String()
	baPortalName := userA.UserID.String()
	if userB.Profile != nil && userB.Profile.Name != nil {
		abPortalName = *userB.Profile.Name
	}
	if userA.Profile != nil && userA.Profile.Name != nil {
		baPortalName = *userA.Profile.Name
	}

	if _, err := createWorldPortal(abPortalName, worldA, worldB); err != nil {
		err := errors.WithMessagef(
			err,
			"Node: apiUsersCreateMutualDocks: failed to create world portal from %s to %s",
			worldA.GetID(), worldB.GetID(),
		)
		api.AbortRequest(c, http.StatusInternalServerError, "ab_portal_create_failed", err, n.log)
		return
	}

	if _, err := createWorldPortal(baPortalName, worldB, worldA); err != nil {
		err := errors.WithMessagef(
			err,
			"Node: apiUsersCreateMutualDocks: failed to create world portal from %s to %s",
			worldB.GetID(), worldA.GetID(),
		)
		api.AbortRequest(c, http.StatusInternalServerError, "ba_portal_create_failed", err, n.log)
		return
	}

	permissions := []*entry.UserSpace{
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

	if err := n.db.GetUserSpaceDB().UpsertUserSpaces(n.ctx, permissions); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersCreateMutualDocks: failed to upsert user spaces")
		api.AbortRequest(c, http.StatusInternalServerError, "upsert_user_spaces_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Remove mutual docks
// @Schemes
// @Description Removes mutual worlds portals and worlds admin permissions
// @Tags users
// @Accept json
// @Produce json
// @Param body body node.apiUsersRemoveMutualDocks.InBody true "body params"
// @Success 202 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/users/mutual-docks [delete]
func (n *Node) apiUsersRemoveMutualDocks(c *gin.Context) {
	type InBody struct {
		WalletA string `json:"walletA" binding:"required"`
		WalletB string `json:"walletB" binding:"required"`
	}
	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userA, err := n.db.GetUsersDB().GetUserByWallet(c, inBody.WalletA)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to get user A by wallet")
		api.AbortRequest(c, http.StatusNotFound, "user_a_not_found", err, n.log)
		return
	}

	userB, err := n.db.GetUsersDB().GetUserByWallet(c, inBody.WalletB)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to get user B by wallet")
		api.AbortRequest(c, http.StatusNotFound, "user_b_not_found", err, n.log)
		return
	}

	worldA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err := errors.Errorf("Node: apiUsersRemoveMutualDocks: failed to get user A world: %s", userA.UserID)
		api.AbortRequest(c, http.StatusNotFound, "user_a_world_not_found", err, n.log)
		return
	}

	worldB, ok := n.GetWorlds().GetWorld(userB.UserID)
	if !ok {
		err := errors.Errorf("Node: apiUsersRemoveMutualDocks: failed to get user B world: %s", userB.UserID)
		api.AbortRequest(c, http.StatusNotFound, "user_b_world_not_found", err, n.log)
		return
	}

	portalsA := getWorldPortals(worldA, worldB)
	portalsB := getWorldPortals(worldB, worldA)
	for _, portal := range utils.MergeMaps(portalsA, portalsB) {
		if _, err := helper.RemoveSpaceFromParent(portal.GetParent(), portal, true); err != nil {
			err := errors.WithMessagef(
				err, "Node: apiUsersRemoveMutualDocks: failed to remove portal: %s", portal.GetID(),
			)
			api.AbortRequest(c, http.StatusInternalServerError, "portal_remove_failed", err, n.log)
			return
		}
	}

	permissions := []*entry.UserSpace{
		{
			SpaceID: worldA.GetID(),
			UserID:  userB.UserID,
		},
		{
			SpaceID: worldB.GetID(),
			UserID:  userA.UserID,
		},
	}

	if err := n.db.GetUserSpaceDB().RemoveUserSpaces(n.ctx, permissions); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to remove user spaces")
		api.AbortRequest(c, http.StatusInternalServerError, "user_spaces_remove_failed", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func createWorldPortal(portalName string, from, to universe.World) (uuid.UUID, error) {
	portals := getWorldPortals(from, to)
	if len(portals) > 0 {
		for portalID, _ := range portals {
			return portalID, nil
		}
	}

	dockingStation, err := getWorldDockingStation(from)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get docking station")
	}

	portalSpaceTypeID, err := helper.GetPortalSpaceTypeID()
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get portal space type id")
	}

	template := helper.SpaceTemplate{
		SpaceName:   &portalName,
		SpaceTypeID: portalSpaceTypeID,
		ParentID:    dockingStation.GetID(),
		SpaceAttributes: []*entry.Attribute{
			entry.NewAttribute(
				entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.TeleportDestination.Name),
				entry.NewAttributePayload(
					&entry.AttributeValue{
						universe.ReservedAttributes.World.TeleportDestination.Key: to.GetID().String(),
					},
					nil,
				),
			),
		},
	}

	return helper.AddSpaceFromTemplate(&template, true)
}

func getWorldDockingStation(world universe.World) (universe.Object, error) {
	dockingStationID := world.GetSettings().Spaces["docking_station"]
	dockingStation, ok := world.GetObjectFromAllObjects(dockingStationID)
	if !ok {
		return nil, errors.Errorf("failed to get docking station space: %s", dockingStationID)
	}
	return dockingStation, nil
}

func getWorldPortals(from, to universe.World) map[uuid.UUID]universe.Object {
	dockingStation, err := getWorldDockingStation(from)
	if err != nil {
		return nil
	}

	toWorld := to.GetID().String()
	attributeID := entry.NewAttributeID(
		universe.GetSystemPluginID(), universe.ReservedAttributes.World.TeleportDestination.Name,
	)
	findPortalFn := func(spaceID uuid.UUID, space universe.Object) bool {
		value, ok := space.GetObjectAttributes().GetValue(attributeID)
		if !ok || value == nil {
			return false
		}

		destination := utils.GetFromAnyMap(*value, universe.ReservedAttributes.World.TeleportDestination.Key, "")
		if destination == toWorld {
			return true
		}

		return false
	}

	return dockingStation.FilterObjects(findPortalFn, false)
}
