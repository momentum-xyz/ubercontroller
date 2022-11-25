package node

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/pkg/errors"
	"net/http"
)

// @Summary Create space
// @Schemes
// @Description Creates a space base on body
// @Tags spaces
// @Accept json
// @Produce json
// @Param body body node.apiCreateSpace.InBody true "body params"
// @Success 201 {object} node.apiCreateSpace.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces [post]
func (n *Node) apiCreateSpace(c *gin.Context) {
	type InBody struct {
		SpaceName   string      `json:"space_name" binding:"required"`
		ParentID    string      `json:"parent_id" binding:"required"`
		SpaceTypeID string      `json:"space_type_id" binding:"required"`
		Asset2dID   string      `json:"asset_2d_id"`
		Asset3dID   string      `json:"asset_3d_id"`
		Position    *cmath.Vec3 `json:"position"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiCreateSpace: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	parentID, err := uuid.Parse(inBody.ParentID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse parent id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_parent_id", err, n.log)
		return
	}
	parent, ok := n.GetSpaceFromAllSpaces(parentID)
	if !ok {
		err := errors.Errorf("Node: apiCreateSpace: parent not found")
		api.AbortRequest(c, http.StatusBadRequest, "parent_not_found", err, n.log)
		return
	}

	spaceID := uuid.New()

	space, err := parent.CreateSpace(spaceID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to create space")
		api.AbortRequest(c, http.StatusInternalServerError, "create_space_failed", err, n.log)
		return
	}

	// TODO: revert on error

	ownerID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to get owner id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_owner_id", err, n.log)
		return
	}
	if err := space.SetOwnerID(ownerID, false); err != nil {
		err := errors.Errorf("Node: apiCreateSpace: failed to set owner id")
		api.AbortRequest(c, http.StatusInternalServerError, "set_owner_id_failed", err, n.log)
		return
	}

	spaceTypeID, err := uuid.Parse(inBody.SpaceTypeID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse space type id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_type_id", err, n.log)
		return
	}
	spaceType, ok := n.GetSpaceTypes().GetSpaceType(spaceTypeID)
	if !ok {
		err := errors.Errorf("Node: apiCreateSpace: space type not found")
		api.AbortRequest(c, http.StatusBadRequest, "space_type_not_found", err, n.log)
		return
	}
	if err := space.SetSpaceType(spaceType, false); err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to set space type")
		api.AbortRequest(c, http.StatusInternalServerError, "set_space_type_failed", err, n.log)
		return
	}

	if inBody.Position != nil {
		if err := space.SetPosition(inBody.Position, false); err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to set position")
			api.AbortRequest(c, http.StatusInternalServerError, "set_position_failed", err, n.log)
			return
		}
	}

	// TODO: should be available for admin or owner of parent
	if inBody.Asset2dID != "" {
		asset2dID, err := uuid.Parse(inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse asset 2d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
			return
		}
		asset2d, ok := n.GetAssets2d().GetAsset2d(asset2dID)
		if !ok {
			err := errors.Errorf("Node: apiCreateSpace: asset 2d not found")
			api.AbortRequest(c, http.StatusBadRequest, "asset_2d_not_found", err, n.log)
			return
		}
		if err := space.SetAsset2D(asset2d, false); err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to set asset 2d")
			api.AbortRequest(c, http.StatusInternalServerError, "set_asset_2d_failed", err, n.log)
			return
		}
	}

	// TODO: should be available for admin or owner of parent
	if inBody.Asset3dID != "" {
		asset3dID, err := uuid.Parse(inBody.Asset3dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse asset 3d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3d, ok := n.GetAssets3d().GetAsset3d(asset3dID)
		if !ok {
			err := errors.Errorf("Node: apiCreateSpace: asset 3d not found")
			api.AbortRequest(c, http.StatusBadRequest, "asset_3d_not_found", err, n.log)
			return
		}
		if err := space.SetAsset3D(asset3d, false); err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to set asset 3d")
			api.AbortRequest(c, http.StatusInternalServerError, "set_asset_3d_failed", err, n.log)
			return
		}
	}

	if err := parent.AddSpace(space, true); err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to add space")
		api.AbortRequest(c, http.StatusInternalServerError, "add_space_failed", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Name.Name)
	value := entry.NewAttributeValue()
	(*value)[universe.Attributes.Space.Name.Key] = inBody.SpaceName
	payload := entry.NewAttributePayload(value, nil)

	if _, err := space.UpsertSpaceAttribute(attributeID, modify.MergeWith(payload), true); err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to upsert space name attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "upsert_space_attribute_failed", err, n.log)
		return
	}

	if err := parent.UpdateChildrenPosition(true, false); err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to update children position")
		api.AbortRequest(c, http.StatusInternalServerError, "update_children_position_failed", err, n.log)
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
