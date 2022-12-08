package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
)

type SpaceTemplate struct {
	SpaceID         *uuid.UUID           `json:"space_id"`
	SpaceName       string               `json:"space_name"`
	SpaceTypeID     uuid.UUID            `json:"space_type_id"`
	ParentID        uuid.UUID            `json:"parent_id"`
	OwnerID         *uuid.UUID           `json:"owner_id"`
	Asset2dID       *uuid.UUID           `json:"asset_2d_id"`
	Asset3dID       *uuid.UUID           `json:"asset_3d_id"`
	Options         *entry.SpaceOptions  `json:"options"`
	Position        *cmath.SpacePosition `json:"position"`
	SpaceAttributes []*Attribute         `json:"space_attributes"`
	Spaces          []*SpaceTemplate     `json:"spaces"`
}

// workaround for mapstructure errors
type Attribute struct {
	entry.AttributeID      `json:",squash"`
	entry.AttributePayload `json:",squash"`
}

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
		SpaceName   string               `json:"space_name" binding:"required"`
		ParentID    string               `json:"parent_id" binding:"required"`
		SpaceTypeID string               `json:"space_type_id" binding:"required"`
		Asset2dID   string               `json:"asset_2d_id"`
		Asset3dID   string               `json:"asset_3d_id"`
		Position    *cmath.SpacePosition `json:"position"`
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

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to get user id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	// is admin check
	userIDs, err := n.db.UserSpaceGetIndirectAdmins(c, parent.GetID())
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to get user space entry for parent space")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_user_space_entry", err, log)
		return
	}

	isAdmin := false

	for _, uID := range userIDs {
		if uID != nil && *uID == userID {
			isAdmin = true
		}
	}

	if !isAdmin {
		err := errors.Errorf("Node: apiCreateSpace: user does not have the permissions to create a new space")
		api.AbortRequest(c, http.StatusUnauthorized, "unauthorized_create_space", err, n.log)
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

	if err := space.SetOwnerID(userID, false); err != nil {
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

	// TODO: should be available for admin or owner of parent
	var asset2dID *uuid.UUID
	if inBody.Asset2dID != "" {
		assetID, err := uuid.Parse(inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse asset 2d id")
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
			err := errors.WithMessage(err, "Node: apiCreateSpace: failed to parse asset 3d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3dID = &assetID
	}

	spaceTemplate := SpaceTemplate{
		SpaceName:   inBody.SpaceName,
		SpaceTypeID: spaceTypeID,
		ParentID:    parentID,
		OwnerID:     &userID,
		Asset2dID:   asset2dID,
		Asset3dID:   asset3dID,
		Position:    inBody.Position,
	}

	tempSpaceID, err := n.addSpaceFromTemplate(&spaceTemplate)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCreateSpace: failed to add space from template")
		api.AbortRequest(c, http.StatusInternalServerError, "add_space_failed", err, n.log)
		return
	}

	type Out struct {
		SpaceID string `json:"space_id"`
	}
	out := Out{
		SpaceID: tempSpaceID.String(),
	}

	c.JSON(http.StatusCreated, out)
}

func (n *Node) addSpaceFromTemplate(spaceTemplate *SpaceTemplate) (uuid.UUID, error) {
	// loading
	spaceType, ok := n.GetSpaceTypes().GetSpaceType(spaceTemplate.SpaceTypeID)
	if !ok {
		return uuid.Nil, errors.Errorf("failed to get space type: %s", spaceTemplate.SpaceTypeID)
	}

	parent, ok := n.GetSpaceFromAllSpaces(spaceTemplate.ParentID)
	if !ok {
		return uuid.Nil, errors.Errorf("parent space not found: %s", spaceTemplate.ParentID)
	}

	// TODO: should be available for admin or owner of parent
	var asset2d universe.Asset2d
	if spaceTemplate.Asset2dID != nil {
		asset2d, ok = n.GetAssets2d().GetAsset2d(*spaceTemplate.Asset2dID)
		if !ok {
			return uuid.Nil, errors.Errorf("asset 2d not found: %s", spaceTemplate.Asset2dID)
		}
	}

	// TODO: should be available for admin or owner of parent
	var asset3d universe.Asset3d
	if spaceTemplate.Asset3dID != nil {
		asset3d, ok = n.GetAssets3d().GetAsset3d(*spaceTemplate.Asset3dID)
		if !ok {
			return uuid.Nil, errors.Errorf("asset 3d not found: %s", spaceTemplate.Asset3dID)
		}
	}

	spaceID := spaceTemplate.SpaceID
	ownerID := spaceTemplate.OwnerID
	if spaceID == nil {
		spaceID = utils.GetPTR(uuid.New())
	}
	if ownerID == nil {
		ownerID = utils.GetPTR(parent.GetOwnerID())
	}

	// creation
	space, err := parent.CreateSpace(*spaceID)
	if err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to create space: %s", spaceID)
	}

	if err := space.SetOwnerID(*ownerID, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set owner id: %s", ownerID)
	}
	if err := space.SetSpaceType(spaceType, false); err != nil {
		return uuid.Nil, errors.WithMessagef(err, "failed to set space type: %s", spaceTemplate.SpaceTypeID)
	}
	if spaceTemplate.Position != nil {
		if err := space.SetPosition(spaceTemplate.Position, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set position: %+v", spaceTemplate.Position)
		}
	}
	if asset2d != nil {
		if err := space.SetAsset2D(asset2d, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set asset 2d: %s", spaceTemplate.Asset2dID)
		}
	}
	if asset3d != nil {
		if err := space.SetAsset3D(asset3d, false); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to set asset 3d: %s", spaceTemplate.Asset3dID)
		}
	}

	if err := parent.AddSpace(space, true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to add space")
	}

	// adding attributes
	spaceTemplate.SpaceAttributes = append(
		spaceTemplate.SpaceAttributes,
		&Attribute{
			AttributeID: entry.NewAttributeID(universe.GetSystemPluginID(), universe.Attributes.Space.Name.Name),
			AttributePayload: entry.AttributePayload{
				Value: &entry.AttributeValue{
					universe.Attributes.Space.Name.Key: spaceTemplate.SpaceName,
				},
			},
		},
	)
	for i := range spaceTemplate.SpaceAttributes {
		if _, err := space.UpsertSpaceAttribute(
			spaceTemplate.SpaceAttributes[i].AttributeID,
			modify.MergeWith(&spaceTemplate.SpaceAttributes[i].AttributePayload),
			true,
		); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to upsert space attribute: %+v", spaceTemplate.SpaceAttributes[i])
		}
	}

	// run
	if err := parent.UpdateChildrenPosition(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update children position")
	}
	if err := space.Run(); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to run space")
	}

	space.SetEnabled(true)

	if err := space.Update(true); err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to update space")
	}

	// children
	for i := range spaceTemplate.Spaces {
		spaceTemplate.Spaces[i].ParentID = *spaceID
		if _, err := n.addSpaceFromTemplate(spaceTemplate.Spaces[i]); err != nil {
			return uuid.Nil, errors.WithMessagef(err, "failed to add space from template: %+v", spaceTemplate.Spaces[i])
		}
	}

	return *spaceID, nil
}
