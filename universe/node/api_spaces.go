package node

import (
	"net/http"

	"github.com/momentum-xyz/ubercontroller/universe/common/helper"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

// @Summary Create space
// @Schemes
// @Description Creates a space base on body
// @Tags spaces
// @Accept json
// @Produce json
// @Param body body node.apiSpacesCreateSpace.InBody true "body params"
// @Success 201 {object} node.apiSpacesCreateSpace.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces [post]
func (n *Node) apiSpacesCreateSpace(c *gin.Context) {
	// TODO: use "helper.SpaceTemplate" alternative here to have ability to create composite objects
	// QUESTION: can we automatically clone "helper.SpaceTemplate" definition and add validation tags to it?
	type InBody struct {
		SpaceName   string               `json:"space_name" binding:"required"`
		ParentID    string               `json:"parent_id" binding:"required"`
		SpaceTypeID string               `json:"space_type_id" binding:"required"`
		Asset2dID   *string              `json:"asset_2d_id"`
		Asset3dID   *string              `json:"asset_3d_id"`
		Position    *cmath.SpacePosition `json:"position"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	parentID, err := uuid.Parse(inBody.ParentID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse parent id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_parent_id", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to get user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	isAdmin, err := n.db.UserSpaceCheckIsUserIndirectSpaceAdmin(c, userID, parentID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to check space indirect admin")
		api.AbortRequest(c, http.StatusBadRequest, "admin_check_failed", err, n.log)
		return
	}

	if !isAdmin {
		err := errors.New("Node: apiSpacesCreateSpace: operation is not permitted for user")
		api.AbortRequest(c, http.StatusUnauthorized, "space_creation_not_permitted", err, n.log)
		return
	}

	position := inBody.Position
	if position == nil {
		position, err = helper.CalcSpaceSpawnPosition(parentID, userID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to calc space spawn position")
			api.AbortRequest(c, http.StatusBadRequest, "calc_spawn_position_failed", err, n.log)
			return
		}
	}

	spaceTypeID, err := uuid.Parse(inBody.SpaceTypeID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse space type id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_type_id", err, n.log)
		return
	}

	var asset2dID *uuid.UUID
	if inBody.Asset2dID != nil {
		assetID, err := uuid.Parse(*inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse asset 2d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
			return
		}
		asset2dID = &assetID
	}

	var asset3dID *uuid.UUID
	if inBody.Asset3dID != nil {
		assetID, err := uuid.Parse(*inBody.Asset3dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse asset 3d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3dID = &assetID
	}

	spaceTemplate := helper.SpaceTemplate{
		SpaceName:   &inBody.SpaceName,
		SpaceTypeID: spaceTypeID,
		ParentID:    parentID,
		OwnerID:     &userID,
		Asset2dID:   asset2dID,
		Asset3dID:   asset3dID,
		Position:    position,
	}

	spaceID, err := helper.AddSpaceFromTemplate(&spaceTemplate, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to add space from template")
		api.AbortRequest(c, http.StatusInternalServerError, "add_space_failed", err, n.log)
		return
	}

	type Out struct {
		SpaceID string `json:"space_id"`
	}
	out := Out{
		SpaceID: spaceID.String(),
	}

	c.JSON(http.StatusCreated, out)
}

// TODO: it was created only for tests, fix or remove
func (n *Node) apiSpacesCreateSpaceFromTemplate(c *gin.Context) {
	var template helper.SpaceTemplate

	if err := c.ShouldBindJSON(&template); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	spaceID, err := helper.AddSpaceFromTemplate(&template, true)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"space_id": spaceID,
	})
}
