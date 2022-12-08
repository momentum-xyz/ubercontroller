package node

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

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

	userA, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletA)
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

	userB, err := n.db.UsersGetUserByWallet(n.ctx, inBody.WalletB)
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

	worldA, ok := n.GetWorlds().GetWorld(userA.UserID)
	if !ok {
		err = errors.New("Node: apiUsersMutualDocks: failed to GetSpace for userA")
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
		err = errors.New("Node: apiUsersMutualDocks: failed to GetSpace for userB")
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
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to addDockingBulb for userA")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	bulbBID, err := n.addDockingBulb(worldB, dockStationType.GetID(), userB, userA)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to addDockingBulb for userB")
		api.AbortRequest(c, http.StatusInternalServerError, "server_error", err, n.log)
		return
	}

	userSpaces := make([]*entry.UserSpace, 2)

	userSpaces[0] = &entry.UserSpace{
		SpaceID:   bulbAID,
		UserID:    userB.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Value:     map[string]any{"role": "admin"},
	}

	userSpaces[1] = &entry.UserSpace{
		SpaceID:   bulbBID,
		UserID:    userA.UserID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Value:     map[string]any{"role": "admin"},
	}

	err = n.db.UserSpacesUpsertUserSpaces(n.ctx, userSpaces)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUsersMutualDocks: failed to UserSpacesUpsertUserSpaces")
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

func (n *Node) addDockingBulb(world universe.World, dockStationTypeID uuid.UUID, fromUser *entry.User, toUser *entry.User) (uuid.UUID, error) {
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
