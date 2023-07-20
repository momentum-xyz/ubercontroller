package node

import (
	"fmt"
	"net/http"
	"time"

	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/universe/logic/tree"
)

// @Summary Generate Agora token
// @Schemes
// @Description Returns an Agora token
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiGenAgoraToken.Body false "body params"
// @Success 200 {object} node.apiGenAgoraToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/agora/token [post]
func (n *Node) apiGenAgoraToken(c *gin.Context) {
	type Body struct {
		ScreenShare bool `json:"screenshare"`
	}
	var inBody Body

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	if _, ok := n.GetObjectFromAllObjects(objectID); !ok {
		err := errors.Errorf("Node: apiGenAgoraToken: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	// 1 day in seconds
	expireSeconds := uint32(1 * 24 * 60 * 60)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expire := currentTimestamp + expireSeconds
	var channel string
	if inBody.ScreenShare {
		channel = fmt.Sprintf("ss|%s", objectID)
	} else {
		channel = objectID.String()
	}
	token, err := rtctokenbuilder.BuildTokenWithUserAccount(
		n.CFG.UIClient.AgoraAppID,
		n.CFG.Common.AgoraAppCertificate,
		channel,
		userID.String(),
		rtctokenbuilder.RolePublisher,
		expire,
	)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to get token")
		api.AbortRequest(c, http.StatusInternalServerError, "get_token_failed", err, n.log)
		return
	}

	type Out struct {
		Token   string `json:"token"`
		Channel string `json:"channel"`
	}
	out := Out{
		Token:   token,
		Channel: channel,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get object by UMID
// @Schemes
// @Description Returns a object info based on UMID and query
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.apiGetObject.InQuery false "query params"
// @Success 202 {object} dto.Object
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id} [get]
func (n *Node) apiGetObject(c *gin.Context) {
	type InQuery struct {
		Effective bool `form:"effective"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetObject: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiGetObject: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	out := dto.Object{
		OwnerID: object.GetOwnerID().String(),
	}
	parent := object.GetParent()
	position := object.GetActualTransform()
	objectType := object.GetObjectType()
	if parent != nil {
		out.ParentID = parent.GetID().String()
	}
	if position != nil {
		out.Transform = *position
	}
	if objectType != nil {
		out.ObjectTypeID = objectType.GetID().String()
	}

	asset2d := object.GetAsset2D()
	asset3d := object.GetAsset3D()
	if inQuery.Effective {
		if asset2d == nil && objectType != nil {
			asset2d = objectType.GetAsset2d()
		}
		if asset3d == nil && objectType != nil {
			asset3d = objectType.GetAsset3d()
		}
	}
	if asset2d != nil {
		out.Asset2dID = asset2d.GetID().String()
	}
	if asset3d != nil {
		out.Asset3dID = asset3d.GetID().String()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Delete a object by UMID
// @Schemes
// @Description Deletes a object by UMID
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id} [delete]
func (n *Node) apiRemoveObject(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveObject: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObject: failed to get user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	isAdmin, err := n.db.GetUserObjectsDB().CheckIsIndirectAdminByID(c, entry.NewUserObjectID(userID, objectID))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObject: failed to check object indirect admin")
		api.AbortRequest(c, http.StatusBadRequest, "admin_check_failed", err, n.log)
		return
	}

	if !isAdmin {
		err := errors.New("Node: apiRemoveObject: operation is not permitted for user")
		api.AbortRequest(c, http.StatusForbidden, "object_remove_not_permitted", err, n.log)
		return
	}

	objectCategory := object.GetObjectType().GetCategoryName()
	if objectCategory == "Worlds" {
		err := errors.New("Node: apiRemoveObject: worlds are not removable")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	removed, err := tree.RemoveObjectFromParent(object.GetParent(), object, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveObject: failed to remove object from parent")
		api.AbortRequest(c, http.StatusInternalServerError, "remove_failed", err, n.log)
		return
	}

	if !removed {
		err := errors.Errorf("Node: apiRemoveObject: object not found in parent")
		api.AbortRequest(c, http.StatusNotFound, "object_not_found_in_parent", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Update a object by UMID
// @Description Updates a object by UMID, 're-parenting' not supported, returns updated object UMID.
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiUpdateObject.InBody true "body params"
// @Success 200 {object} node.apiUpdateObject.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id} [patch]
func (n *Node) apiUpdateObject(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUpdateObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	// TODO: ask @cnaize about it
	// not supporting 're-parenting' and changing type'. Have to delete and recreate for that.
	// Update/edit the positioning is done through unity edit mode.
	type InBody struct {
		ObjectName string  `json:"object_name"`
		Asset2dID  *string `json:"asset_2d_id"`
		Asset3dID  string  `json:"asset_3d_id"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiUpdateObject: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiUpdateObject: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	var asset2d universe.Asset2d
	if inBody.Asset2dID != nil {
		if *inBody.Asset2dID != "" {
			asset2dID, err := umid.Parse(*inBody.Asset2dID)
			if err != nil {
				err := errors.WithMessage(err, "Node: apiUpdateObject: failed to parse asset 2d umid")
				api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
				return
			}
			asset2d, ok = n.GetAssets2d().GetAsset2d(asset2dID)
			if !ok {
				err := errors.Errorf("Node: apiUpdateObject: 2D asset not found: %s", asset2dID)
				api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
				return
			}
			if asset2d != nil {
				if err := object.SetAsset2D(asset2d, true); err != nil {
					err := errors.Errorf("Node: apiUpdateObject: failed to update 2d asset: %s", asset2d.GetID())
					api.AbortRequest(c, http.StatusInternalServerError, "object_asset_2d", err, n.log)
					return
				}
			}
		} else { // empty string asset 2d, allow resetting it
			if err := object.SetAsset2D(nil, true); err != nil {
				err := errors.Errorf("Node: apiUpdateObject: failed to clear 2d asset: %s", asset2d.GetID())
				api.AbortRequest(c, http.StatusInternalServerError, "object_asset_2d", err, n.log)
				return
			}

		}
	}

	var asset3d universe.Asset3d
	if inBody.Asset3dID != "" {
		asset3dID, err := umid.Parse(inBody.Asset3dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiUpdateObject: failed to parse asset 3d umid")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3d, ok = n.GetAssets3d().GetAsset3d(asset3dID)
		if !ok {
			err := errors.Errorf("Node: apiUpdateObject: 3D asset not found: %s", asset3dID)
			api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
			return
		}
	}
	if asset3d != nil {
		if err := object.SetAsset3D(asset3d, true); err != nil {
			err := errors.Errorf("Node: apiUpdateObject: failed to update 3d asset: %s", asset3d.GetID())
			api.AbortRequest(c, http.StatusInternalServerError, "object_asset_3d", err, n.log)
			return
		}
	}

	if inBody.ObjectName != "" {
		if err := object.SetName(inBody.ObjectName, true); err != nil {
			err := errors.WithMessagef(err, "Node: apiUpdateObject: failed to set object name: %s", inBody.ObjectName)
			api.AbortRequest(c, http.StatusInternalServerError, "object_name", err, n.log)
			return
		}
	}

	if err := object.Update(false); err != nil {
		err = errors.WithMessage(err, "Node: apiUpdateObject: failed to update object")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_object_update", err, n.log)
		return
	}

	// TODO: output full object data
	type Out struct {
		ObjectID string `json:"object_id"`
	}
	out := Out{
		ObjectID: objectID.String(),
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Clones an object by UMID
// @Description Clones an object by UMID, 're-parenting' not supported, returns cloned object UMID.
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiCloneObject.InBody true "body params"
// @Success 200 {object} node.apiCloneObject.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/clone [post]
func (n *Node) apiCloneObject(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCloneObject: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiCloneObject: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCloneObject: failed to get user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	parentID := object.GetParent().GetID()
	transform, err := tree.CalcObjectSpawnPosition(parentID, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCloneObject: failed to calc object spawn position")
		api.AbortRequest(c, http.StatusBadRequest, "calc_spawn_position_failed", err, n.log)
		return
	}

	objectName := object.GetName()
	objectOwnerID := object.GetOwnerID()
	asset2dID := object.GetAsset2D().GetID()
	asset3dID := object.GetAsset3D().GetID()

	objectTemplate := tree.ObjectTemplate{
		ObjectName:   &objectName,
		ObjectTypeID: object.GetObjectType().GetID(),
		ParentID:     parentID,
		OwnerID:      &objectOwnerID,
		Asset2dID:    &asset2dID,
		Asset3dID:    &asset3dID,
		Transform:    transform,
		Options:      object.GetOptions(),
	}

	clonedObjectID, err := tree.AddObjectFromTemplate(&objectTemplate, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCloneObject: failed to clone object")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_clone", err, n.log)
		return
	}

	type CloneableObjectAttribute struct {
		AttributeID entry.AttributeID
		Value       *entry.AttributeValue
	}

	var cloneableAttributes []CloneableObjectAttribute
	objectAttributes := object.GetObjectAttributes().GetAll()
	for attributeID := range objectAttributes {
		effectiveOptions, _ := object.GetObjectAttributes().GetEffectiveOptions(attributeID)
		if effectiveOptions != nil {
			cloneable := utils.GetFromAnyMap(*effectiveOptions, "cloneable", map[string]any(nil))
			if cloneable != nil {
				defaultValue := utils.GetFromAnyMap(*effectiveOptions, "use_default", entry.AttributeValue(nil))
				if defaultValue != nil {
					cloneableObjectAttribute := CloneableObjectAttribute{
						AttributeID: attributeID,
						Value:       &defaultValue,
					}

					cloneableAttributes = append(cloneableAttributes, cloneableObjectAttribute)
				}
			}
		}
	}

	clonedObject, ok := n.GetObjectFromAllObjects(clonedObjectID)
	if !ok {
		err := errors.Errorf("Node: apiCloneObject: object not found: %s", clonedObjectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	for _, clAtr := range cloneableAttributes {
		var attributeModifyFunc modify.Fn[entry.AttributePayload]
		attributeModifyFunc = func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
			if current == nil {
				current = entry.NewAttributePayload(nil, nil)
			}

			if current.Value == nil {
				current.Value = entry.NewAttributeValue()
			}

			current.Value = clAtr.Value

			return current, nil
		}

		_, err = clonedObject.GetObjectAttributes().Upsert(clAtr.AttributeID, attributeModifyFunc, true)
		if err != nil {
			err = errors.WithMessage(err, "Node: apiCloneObject: failed to upsert object attribute")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}
	}

	type Out struct {
		ObjectID string `json:"object_id"`
	}
	out := Out{
		ObjectID: clonedObject.GetID().String(),
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set object sub option by object UMID
// @Schemes
// @Description Sets a object sub option by object UMID, returns appended object option
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiObjectsSetObjectSubOption.Body true "body params"
// @Success 202 {object} dto.ObjectSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/options/sub [post]
func (n *Node) apiObjectsSetObjectSubOption(c *gin.Context) {
	type Body struct {
		SubOptionKey   string `json:"sub_option_key" binding:"required"`
		SubOptionValue any    `json:"sub_option_value" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiObjectsSetObjectSubOption: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsSetObjectSubOption: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiObjectsSetObjectSubOption: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.ObjectOptions) (*entry.ObjectOptions, error) {
		if current == nil {
			current = &entry.ObjectOptions{}
		}
		if current.Subs == nil {
			current.Subs = make(map[string]any)
		}

		current.Subs[inBody.SubOptionKey] = inBody.SubOptionValue

		return current, nil
	}

	if _, err := object.SetOptions(modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsSetObjectSubOption: failed to set options")
		api.AbortRequest(c, http.StatusInternalServerError, "set_options_failed", err, n.log)
		return
	}

	out := dto.ObjectSubOptions{
		inBody.SubOptionKey: inBody.SubOptionValue,
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete object sub option by object UMID
// @Schemes
// @Description Deletes a object sub option by object UMID
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiObjectsRemoveObjectSubOption.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/options/sub [delete]
func (n *Node) apiObjectsRemoveObjectSubOption(c *gin.Context) {
	type Body struct {
		SubOptionKey string `json:"sub_option_key" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiObjectsRemoveObjectSubOption: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsRemoveObjectSubOption: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiObjectsRemoveObjectSubOption: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.ObjectOptions) (*entry.ObjectOptions, error) {
		if current == nil || current.Subs == nil {
			return current, nil
		}

		delete(current.Subs, inBody.SubOptionKey)

		return current, nil
	}

	if _, err := object.SetOptions(modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsRemoveObjectSubOption: failed to set options")
		api.AbortRequest(c, http.StatusInternalServerError, "set_options_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Get object options by object UMID
// @Schemes
// @Description Returns a object options based on object UMID and query
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.apiObjectsGetObjectOptions.InQuery false "query params"
// @Success 200 {object} dto.ObjectOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/options [get]
func (n *Node) apiObjectsGetObjectOptions(c *gin.Context) {
	type InQuery struct {
		Effective bool `form:"effective"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsGetObjectOptions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsGetObjectOptions: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiObjectsGetObjectOptions: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	var out dto.ObjectOptions
	if inQuery.Effective {
		out = object.GetEffectiveOptions()
	} else {
		out = object.GetOptions()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get object sub options
// @Schemes
// @Description Returns a object sub options based on query
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param query query node.apiObjectsGetObjectSubOptions.InQuery true "query params"
// @Success 200 {object} dto.ObjectSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/options/sub [get]
func (n *Node) apiObjectsGetObjectSubOptions(c *gin.Context) {
	type InQuery struct {
		Effective    bool   `form:"effective"`
		SubOptionKey string `form:"sub_option_key" binding:"required"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsGetObjectSubOptions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiObjectsGetObjectSubOptions: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiObjectsGetObjectSubOptions: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	var options *entry.ObjectOptions
	if inQuery.Effective {
		options = object.GetEffectiveOptions()
	} else {
		options = object.GetOptions()
	}

	if options == nil {
		err := errors.Errorf("Node: apiObjectsGetObjectSubOptions: empty options")
		api.AbortRequest(c, http.StatusNotFound, "empty_options", err, n.log)
		return
	}

	if options.Subs == nil {
		err := errors.Errorf("Node: apiObjectsGetObjectSubOptions: empty sub options")
		api.AbortRequest(c, http.StatusNotFound, "empty_sub_options", err, n.log)
		return
	}

	out := dto.ObjectSubOptions{
		inQuery.SubOptionKey: options.Subs[inQuery.SubOptionKey],
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Claim and customise object
// @Schemes
// @Description Claim and customise object
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Param body body node.apiClaimAndCustomise.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/claim-and-customise [post]
func (n *Node) apiClaimAndCustomise(c *gin.Context) {
	type Body struct {
		Title     string `json:"title" binding:"required"`
		Text      string `json:"text" binding:"required"`
		ImageHash string `json:"image_hash" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiClaimAndCustomise: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	if object.GetObjectType().GetID() != universe.CustomisableObjectTypeID {
		err = errors.New("Node: apiClaimAndCustomise: object not customisable")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	userObjects, err := n.userObjects.GetUserObjectsByObjectID(object.GetID())
	if err != nil {
		err = errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to GetUserObjectsByObjectID")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	if userObjects != nil {
		for _, uo := range userObjects {
			if uo.Value != nil && (*uo.Value)["role"] == "admin" && uo.UserID != userID {
				err = errors.New("Node: apiClaimAndCustomise: object already claimed")
				api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
				return
			}
		}
	}

	userObjectID := entry.NewUserObjectID(userID, object.GetID())
	var modifyFunc modify.Fn[entry.UserObjectValue]
	modifyFunc = func(payload *entry.UserObjectValue) (*entry.UserObjectValue, error) {
		if payload == nil {
			payload = entry.NewUserObjectValue()
		}
		(*payload)["role"] = "admin"

		return payload, nil
	}

	_, err = n.userObjects.Upsert(userObjectID, modifyFunc, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to upsert user_object")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), "user_customisable_data")
	var attributeModifyFunc modify.Fn[entry.AttributePayload]
	attributeModifyFunc = func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		if current == nil {
			current = entry.NewAttributePayload(nil, nil)
		}

		if current.Value == nil {
			current.Value = entry.NewAttributeValue()
		}

		value := current.Value
		(*value)["title"] = inBody.Title
		(*value)["text"] = inBody.Text
		(*value)["image_hash"] = inBody.ImageHash
		(*value)["claimed_by"] = userID
		layout := "2006-01-02T15:04:05.999Z"
		(*value)["created_at"] = time.Now().UTC().Format(layout)

		return current, nil
	}

	_, err = object.GetObjectAttributes().Upsert(attributeID, attributeModifyFunc, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiClaimAndCustomise: failed to upsert object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, true)
}

// @Summary Unclaim and clear object customisation
// @Schemes
// @Description Unclaim and clear object customisation
// @Tags objects
// @Accept json
// @Produce json
// @Param object_id path string true "Object UMID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 403 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/objects/{object_id}/unclaim-and-clear-customisation [post]
func (n *Node) apiUnclaimAndClearCustomisation(c *gin.Context) {
	objectID, err := umid.Parse(c.Param("objectID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUnclaimAndClearCustomisation: failed to parse object umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_object_id", err, n.log)
		return
	}

	object, ok := n.GetObjectFromAllObjects(objectID)
	if !ok {
		err := errors.Errorf("Node: apiUnclaimAndClearCustomisation: object not found: %s", objectID)
		api.AbortRequest(c, http.StatusNotFound, "object_not_found", err, n.log)
		return
	}

	if object.GetObjectType().GetID() != universe.CustomisableObjectTypeID {
		err = errors.New("Node: apiUnclaimAndClearCustomisation: object not customisable")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUnclaimAndClearCustomisation: failed to get user umid")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	userObjectID := entry.NewUserObjectID(userID, objectID)
	value, ok := n.userObjects.GetValue(userObjectID)
	if !ok {
		err = errors.New("Node: apiUnclaimAndClearCustomisation: object not claimed by this user")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	if value == nil {
		err = errors.New("Node: apiUnclaimAndClearCustomisation: object not claimed by this user")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	if (*value)["role"] != "admin" {
		err = errors.New("Node: apiUnclaimAndClearCustomisation: object not claimed by this user")
		api.AbortRequest(c, http.StatusForbidden, "forbidden", err, n.log)
		return
	}

	_, err = n.userObjects.Remove(userObjectID, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUnclaimAndClearCustomisation: failed to remove user_object")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	attributeID := entry.NewAttributeID(universe.GetSystemPluginID(), "user_customisable_data")
	_, err = object.GetObjectAttributes().Remove(attributeID, true)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiUnclaimAndClearCustomisation: failed to remove object attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusAccepted, true)
}
