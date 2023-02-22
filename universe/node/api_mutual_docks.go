package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/common"
	"github.com/momentum-xyz/ubercontroller/universe/logic/tree"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
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
	abPortalImage := ""
	baPortalImage := ""
	if userB.Profile != nil && userB.Profile.Name != nil {
		abPortalName = *userB.Profile.Name
	}
	if userB.Profile != nil && userB.Profile.AvatarHash != nil {
		abPortalImage = *userB.Profile.AvatarHash
	}
	if userA.Profile != nil && userA.Profile.Name != nil {
		baPortalName = *userA.Profile.Name
	}
	if userA.Profile != nil && userA.Profile.AvatarHash != nil {
		baPortalImage = *userA.Profile.AvatarHash
	}

	if _, err := createWorldPortal(abPortalName, worldA, worldB, abPortalImage); err != nil {
		err := errors.WithMessagef(
			err,
			"Node: apiUsersCreateMutualDocks: failed to create world portal from %s to %s",
			worldA.GetID(), worldB.GetID(),
		)
		api.AbortRequest(c, http.StatusInternalServerError, "ab_portal_create_failed", err, n.log)
		return
	}

	if _, err := createWorldPortal(baPortalName, worldB, worldA, baPortalImage); err != nil {
		err := errors.WithMessagef(
			err,
			"Node: apiUsersCreateMutualDocks: failed to create world portal from %s to %s",
			worldB.GetID(), worldA.GetID(),
		)
		api.AbortRequest(c, http.StatusInternalServerError, "ba_portal_create_failed", err, n.log)
		return
	}

	permissions := []*entry.UserObject{
		entry.NewUserObject(
			entry.NewUserObjectID(userA.UserID, worldB.GetID()),
			&entry.UserObjectValue{"role": "admin"},
		),
		entry.NewUserObject(
			entry.NewUserObjectID(userB.UserID, worldA.GetID()),
			&entry.UserObjectValue{"role": "admin"},
		),
	}

	for i := range permissions {
		if _, err := n.GetUserObjects().Upsert(
			permissions[i].UserObjectID, modify.MergeWith(permissions[i].Value), true,
		); err != nil {
			err := errors.WithMessagef(
				err, "Node: apiUsersCreateMutualDocks: failed to upsert user object: %+v", permissions[i].UserObjectID,
			)
			api.AbortRequest(c, http.StatusInternalServerError, "upsert_user_object_failed", err, n.log)
			return
		}
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
		if _, err := tree.RemoveObjectFromParent(portal.GetParent(), portal, true); err != nil {
			err := errors.WithMessagef(
				err, "Node: apiUsersRemoveMutualDocks: failed to remove portal: %s", portal.GetID(),
			)
			api.AbortRequest(c, http.StatusInternalServerError, "portal_remove_failed", err, n.log)
			return
		}
	}

	permissions := []entry.UserObjectID{
		entry.NewUserObjectID(userA.UserID, worldB.GetID()),
		entry.NewUserObjectID(userB.UserID, worldA.GetID()),
	}

	if _, err := n.GetUserObjects().RemoveMany(permissions, true); err != nil {
		err := errors.WithMessage(err, "Node: apiUsersRemoveMutualDocks: failed to remove user objects")
		api.AbortRequest(c, http.StatusInternalServerError, "user_objects_remove_failed", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, nil)
}

func createWorldPortal(portalName string, from, to universe.World, portalImage string) (uuid.UUID, error) {
	var objectAttributes []*entry.Attribute

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

	portalObjectTypeID, err := common.GetPortalObjectTypeID()
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get portal object type id")
	}

	objectAttributes = append(objectAttributes, entry.NewAttribute(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.World.TeleportDestination.Name),
		entry.NewAttributePayload(
			&entry.AttributeValue{
				universe.ReservedAttributes.World.TeleportDestination.Key: to.GetID().String(),
			},
			nil,
		),
	))

	if portalImage != "" {
		objectAttributes = append(objectAttributes, entry.NewAttribute(
			entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Object.PortalDockFace.Name),
			entry.NewAttributePayload(
				&entry.AttributeValue{
					universe.ReservedAttributes.Object.PortalDockFace.Key: portalImage,
				},
				nil,
			),
		))
	}

	template := tree.ObjectTemplate{
		ObjectName:       &portalName,
		ObjectTypeID:     portalObjectTypeID,
		ParentID:         dockingStation.GetID(),
		ObjectAttributes: objectAttributes,
	}

	return tree.AddObjectFromTemplate(&template, true)
}

func getWorldDockingStation(world universe.World) (universe.Object, error) {
	dockingStationID := world.GetSettings().Objects["docking_station"]
	dockingStation, ok := world.GetObjectFromAllObjects(dockingStationID)
	if !ok {
		return nil, errors.Errorf("failed to get docking station object: %s", dockingStationID)
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
	findPortalFn := func(objectID uuid.UUID, object universe.Object) bool {
		value, ok := object.GetObjectAttributes().GetValue(attributeID)
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
