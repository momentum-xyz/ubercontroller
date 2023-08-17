package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/pkg/cmath"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/tree"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

// @Summary Create object
// @Description Creates a object base on body
// @Tags objects
// @Security Bearer
// @Param body body node.apiObjectsCreateObject.InBody true "body params"
// @Success 201 {object} node.apiObjectsCreateObject.Out
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects [post]
func (n *Node) apiObjectsCreateObject(c *gin.Context) {
	// TODO: use "helper.ObjectTemplate" alternative here to have ability to create composite objects
	// QUESTION: can we automatically clone "helper.ObjectTemplate" definition and add validation tags to it?
	type InBody struct {
		ObjectName   string           `json:"object_name" binding:"required"`
		ParentID     string           `json:"parent_id" binding:"required"`
		ObjectTypeID string           `json:"object_type_id" binding:"required"`
		Asset2dID    *string          `json:"asset_2d_id"`
		Asset3dID    *string          `json:"asset_3d_id"`
		Transform    *cmath.Transform `json:"transform"`
		Minimap      bool             `json:"minimap"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	parentID, err := umid.Parse(inBody.ParentID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to parse parent umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_parent_id", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to get user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	isAdmin, err := n.db.GetUserObjectsDB().CheckIsIndirectAdminByID(c, entry.NewUserObjectID(userID, parentID))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to check object indirect admin")
		api.AbortRequest(c, http.StatusBadRequest, "admin_check_failed", err, n.log)
		return
	}

	if !isAdmin {
		err := errors.New("Node: apiObjectsCreateObject: operation is not permitted for user")
		api.AbortRequest(c, http.StatusForbidden, "object_creation_not_permitted", err, n.log)
		return
	}

	transform := inBody.Transform
	if transform == nil {
		transform, err = tree.CalcObjectSpawnPosition(parentID, userID, nil)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to calc object spawn position")
			api.AbortRequest(c, http.StatusBadRequest, "calc_spawn_position_failed", err, n.log)
			return
		}
	}

	parent, ok := n.GetObjectFromAllObjects(parentID)
	if !ok {
		err := errors.Errorf("Node: apiObjectsCreateObject: failed to get parent object")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_parent", err, n.log)
		return
	}

	parentObjectTypeOptions := parent.GetEffectiveOptions()
	if parentObjectTypeOptions == nil {
		err := errors.Errorf("Node: apiObjectsCreateObject: parent object type does not have options")
		api.AbortRequest(c, http.StatusInternalServerError, "no_options_in_object_type", err, n.log)
		return
	}

	parentObjects := uint(len(parent.GetObjects(false)))
	if parentObjectTypeOptions.ChildLimit != nil && parentObjects > *parentObjectTypeOptions.ChildLimit {
		err := errors.Errorf("Node: apiObjectsCreateObject: child limit reached for parent object")
		api.AbortRequest(c, http.StatusBadRequest, "child_limit_reached", err, n.log)
		return
	}

	if !utils.Contains(parentObjectTypeOptions.AllowedChildren, inBody.ObjectTypeID) {
		err := errors.Errorf("Node: apiObjectsCreateObject: object type is not allowed")
		api.AbortRequest(c, http.StatusBadRequest, "object_type_not_permitted", err, n.log)
		return
	}

	objectTypeID, err := umid.Parse(inBody.ObjectTypeID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to parse object type umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_type_id", err, n.log)
		return
	}

	var asset2dID *umid.UMID
	if inBody.Asset2dID != nil {
		assetID, err := umid.Parse(*inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to parse asset 2d umid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
			return
		}
		asset2dID = &assetID
	}

	var asset3dID *umid.UMID
	if inBody.Asset3dID != nil {
		assetID, err := umid.Parse(*inBody.Asset3dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to parse asset 3d umid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3dID = &assetID
	}

	objectTemplate := tree.ObjectTemplate{
		ObjectName:   &inBody.ObjectName,
		ObjectTypeID: objectTypeID,
		ParentID:     parentID,
		OwnerID:      &userID,
		Asset2dID:    asset2dID,
		Asset3dID:    asset3dID,
		Transform:    transform,
		Options: &entry.ObjectOptions{
			Minimap: utils.GetPTR(inBody.Minimap),
		},
	}

	objectID, err := tree.AddObjectFromTemplate(&objectTemplate, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsCreateObject: failed to add object from template")
		api.AbortRequest(c, http.StatusInternalServerError, "add_object_failed", err, n.log)
		return
	}

	type Out struct {
		ObjectID string `json:"object_id"`
	}
	out := Out{
		ObjectID: objectID.String(),
	}

	c.JSON(http.StatusCreated, out)
}

// TODO: it was created only for tests, fix or remove
func (n *Node) apiObjectsCreateObjectFromTemplate(c *gin.Context) {
	var template tree.ObjectTemplate

	if err := c.ShouldBindJSON(&template); err != nil {
		c.AbortWithStatusJSON(
			http.StatusBadRequest, gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	objectID, err := tree.AddObjectFromTemplate(&template, true)
	if err != nil {
		c.AbortWithStatusJSON(
			http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			},
		)
		return
	}

	c.JSON(
		http.StatusCreated, gin.H{
			"object_id": objectID,
		},
	)
}
