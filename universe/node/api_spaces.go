package node

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/universe/common/helper"
	"net/http"

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
	type InBody struct {
		SpaceName   string               `json:"space_name" binding:"required"`
		ParentID    string               `json:"parent_id" binding:"required"`
		SpaceTypeID string               `json:"space_type_id" binding:"required"`
		Asset2dID   string               `json:"asset_2d_id"`
		Asset3dID   string               `json:"asset_3d_id"`
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

	// TODO: fix this bloody stuff
	position := inBody.Position
	if position == nil {
		parent, ok := n.GetSpaceFromAllSpaces(parentID)
		if !ok {
			err := errors.Errorf("Node: apiSpacesCreateSpace: parent space not found")
			api.AbortRequest(c, http.StatusBadRequest, "parent_not_found", err, n.log)
			return
		}
		options := parent.GetOptions()
		if options == nil || len(options.ChildPlacements) == 0 {
			parentWorld := parent.GetWorld()
			if parentWorld != nil {
				user, ok := parentWorld.GetUser(userID, true)
				if ok {
					fmt.Printf("User rotation: %v", user.GetRotation())
					//distance := float32(10)
					position = &cmath.SpacePosition{
						// TODO: recalc based on euler angles, not lookat: Location: cmath.Add(user.GetPosition(), cmath.MultiplyN(user.GetRotation(), distance)),
						Location: user.GetPosition(),
						Rotation: cmath.Vec3{},
						Scale:    cmath.Vec3{X: 1, Y: 1, Z: 1},
					}
				}
			}
		}
	}

	spaceTypeID, err := uuid.Parse(inBody.SpaceTypeID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse space type id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_type_id", err, n.log)
		return
	}

	// TODO: should be available for admin or owner of parent
	var asset2dID *uuid.UUID
	if inBody.Asset2dID != "" {
		assetID, err := uuid.Parse(inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiSpacesCreateSpace: failed to parse asset 2d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
			return
		}
		asset2dID = &assetID
	}

	// TODO: should be available for admin or owner of parent
	var asset3dID *uuid.UUID
	if inBody.Asset3dID != "" {
		assetID, err := uuid.Parse(inBody.Asset3dID)
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
