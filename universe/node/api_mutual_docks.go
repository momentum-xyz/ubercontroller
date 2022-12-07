package node

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
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

	spaceA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err = errors.New("Node: apiUsersMutualDocks: failed to GetSpace for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	if spaceA == nil {
		err = errors.New("space not found for userID=:" + userA.UserID.String())
		api.AbortRequest(c, http.StatusNotFound, "space_A_not_found", err, n.log)
		return
	}

	spaceB, ok := n.GetWorlds().GetWorld(userB.UserID)
	if !ok {
		err = errors.New("Node: apiUsersMutualDocks: failed to GetSpace for userB")
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

	bulbAID, err := n.addDockingBulb(spaceA, dockStationType.GetID(), userA, userB)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to addDockingBulb for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	bulbBID, err := n.addDockingBulb(spaceB, dockStationType.GetID(), userB, userA)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to addDockingBulb for userB")
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
		SpaceA:          spaceA.GetID().String(),
		SpaceB:          spaceB.GetID().String(),
		DockStationType: dockStationType.GetID().String(),
		BulbA:           bulbAID.String(),
		BulbB:           bulbBID.String(),
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) addDockingBulb(space universe.Space, dockStationTypeID uuid.UUID, fromUser *entry.User, toUser *entry.User) (uuid.UUID, error) {

	dockStation := n.getDockStation(space, dockStationTypeID)

	bulbSpaceID := uuid.New()
	err := n.createAndAddBulbSpace(bulbSpaceID, dockStation, *fromUser, *toUser)
	if err != nil {
		return bulbSpaceID, errors.WithMessage(err, "Node: addDockingBulb: failed to createAndAddBulbSpace")
	}

	return bulbSpaceID, nil
}

func (n *Node) createAndAddBulbSpace(spaceID uuid.UUID, dock universe.Space, fromUser entry.User, toUser entry.User) error {
	space, err := dock.CreateSpace(spaceID)
	if err != nil {
		return errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to create space")
	}

	err = space.SetOwnerID(fromUser.UserID, false)
	if err != nil {
		return errors.Errorf("Node: createAndAddBulbSpace: failed to set owner id")
	}

	bulbSpaceTypeName := "Docking bulb"
	spaceType := n.getSpaceTypeByName(bulbSpaceTypeName)
	if spaceType == nil {
		return errors.New("Node: createAndAddBulbSpace: Can not find space type for name=" + bulbSpaceTypeName)
	}

	err = space.SetSpaceType(spaceType, false)
	if err != nil {
		return errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to set space type")
	}

	asset3d := spaceType.GetAsset3d()
	err = space.SetAsset3D(asset3d, false)
	if err != nil {
		err = errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to set asset 3d")
		return err
	}

	err = dock.AddSpace(space, true)
	if err != nil {
		return errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to add space")
	}

	toUserName := ""
	if toUser.Profile != nil {
		if toUser.Profile.Name != nil {
			toUserName = *toUser.Profile.Name
		}
	}

	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Name.Name)
	value := entry.NewAttributeValue()
	(*value)[universe.Attributes.Space.Name.Key] = toUserName
	payload := entry.NewAttributePayload(value, nil)

	_, err = space.UpsertSpaceAttribute(attributeID, modify.MergeWith(payload), true)
	if err != nil {
		err = errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to upsert space name attribute")
		return err
	}

	attributeID = entry.NewAttributeID(universe.GetSystemPluginID(), "teleport")
	value = entry.NewAttributeValue()
	(*value)["DestinationWorldID"] = toUser.UserID.String()
	payload = entry.NewAttributePayload(value, nil)

	_, err = space.UpsertSpaceAttribute(attributeID, modify.MergeWith(payload), true)
	if err != nil {
		err = errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to upsert teleport name attribute")
		return err
	}

	err = space.Run()
	if err != nil {
		err = errors.WithMessage(err, "Node: createAndAddBulbSpace: failed to run space")
		return err
	}

	space.SetEnabled(true)

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

func (n *Node) getDockStation(space universe.Space, dockStationTypeID uuid.UUID) universe.Space {
	allSpaces := space.GetSpaces(false)

	for _, space := range allSpaces {
		if space.GetSpaceType().GetID() == dockStationTypeID {
			return space
		}
	}

	return nil
}
