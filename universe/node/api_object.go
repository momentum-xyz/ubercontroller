package node

import (
	"fmt"
	"github.com/momentum-xyz/ubercontroller/universe/common/helper"
	"net/http"
	"time"

	"github.com/AgoraIO-Community/go-tokenbuilder/rtctokenbuilder"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Generate Agora token
// @Schemes
// @Description Returns an Agora token
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param body body node.apiGenAgoraToken.Body false "body params"
// @Success 200 {object} node.apiGenAgoraToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/agora/token [post]
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

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	if _, ok := n.GetObjectFromAllObjects(spaceID); !ok {
		err := errors.Errorf("Node: apiGenAgoraToken: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err = errors.WithMessage(err, "Node: apiGenAgoraToken: failed to get user id")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	// 1 day in seconds
	expireSeconds := uint32(1 * 24 * 60 * 60)
	currentTimestamp := uint32(time.Now().UTC().Unix())
	expire := currentTimestamp + expireSeconds
	var channel string
	if inBody.ScreenShare {
		channel = fmt.Sprintf("ss|%s", spaceID)
	} else {
		channel = spaceID.String()
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

// @Summary Get space by ID
// @Schemes
// @Description Returns a space info based on ID and query
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param query query node.apiGetSpace.InQuery false "query params"
// @Success 202 {object} dto.Space
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id} [get]
func (n *Node) apiGetSpace(c *gin.Context) {
	type InQuery struct {
		Effective bool `form:"effective"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpace: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSpace: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiGetSpace: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	out := dto.Space{
		OwnerID: space.GetOwnerID().String(),
	}
	parent := space.GetParent()
	position := space.GetActualPosition()
	spaceType := space.GetObjectType()
	if parent != nil {
		out.ParentID = parent.GetID().String()
	}
	if position != nil {
		out.Position = *position
	}
	if spaceType != nil {
		out.SpaceTypeID = spaceType.GetID().String()
	}

	asset2d := space.GetAsset2D()
	asset3d := space.GetAsset3D()
	if inQuery.Effective {
		if asset2d == nil && spaceType != nil {
			asset2d = spaceType.GetAsset2d()
		}
		if asset3d == nil && spaceType != nil {
			asset3d = spaceType.GetAsset3d()
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

// @Summary Delete a space by ID
// @Schemes
// @Description Deletes a space by ID
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id} [delete]
func (n *Node) apiRemoveSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpace: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSpace: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	removed, err := helper.RemoveObjectFromParent(space.GetParent(), space, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSpace: failed to remove space from parent")
		api.AbortRequest(c, http.StatusInternalServerError, "remove_failed", err, n.log)
		return
	}

	if !removed {
		err := errors.Errorf("Node: apiRemoveSpace: space not found in parent")
		api.AbortRequest(c, http.StatusNotFound, "space_not_found_in_parent", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Update a space by ID
// @Description Updates a space by ID, 're-parenting' not supported, returns updated space ID.
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param body body node.apiUpdateSpace.InBody true "body params"
// @Success 200 {object} node.apiUpdateSpace.Out
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id} [patch]
func (n *Node) apiUpdateSpace(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiUpdateSpace: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	// TODO: ask @cnaize about it
	// not supporting 're-parenting' and changing type'. Have to delete and recreate for that.
	// Update/edit the positioning is done through unity edit mode.
	type InBody struct {
		SpaceName string `json:"space_name"`
		Asset2dID string `json:"asset_2d_id"`
		Asset3dID string `json:"asset_3d_id"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiUpdateSpace: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiUpdateSpace: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	var asset2d universe.Asset2d
	if inBody.Asset2dID != "" {
		asset2dID, err := uuid.Parse(inBody.Asset2dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiUpdateSpace: failed to parse asset 2d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_2d_id", err, n.log)
			return
		}
		asset2d, ok = n.GetAssets2d().GetAsset2d(asset2dID)
		if !ok {
			err := errors.Errorf("Node: apiUpdateSpace: 2D asset not found: %s", asset2dID)
			api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
			return
		}
	}
	if asset2d != nil {
		if err := space.SetAsset2D(asset2d, true); err != nil {
			err := errors.Errorf("Node: apiUpdateSpace: failed to update 2d asset: %s", asset2d.GetID())
			api.AbortRequest(c, http.StatusInternalServerError, "space_asset_2d", err, n.log)
			return
		}
	}

	var asset3d universe.Asset3d
	if inBody.Asset3dID != "" {
		asset3dID, err := uuid.Parse(inBody.Asset3dID)
		if err != nil {
			err := errors.WithMessage(err, "Node: apiUpdateSpace: failed to parse asset 3d id")
			api.AbortRequest(c, http.StatusBadRequest, "invalid_asset_3d_id", err, n.log)
			return
		}
		asset3d, ok = n.GetAssets3d().GetAsset3d(asset3dID)
		if !ok {
			err := errors.Errorf("Node: apiUpdateSpace: 3D asset not found: %s", asset3dID)
			api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
			return
		}
	}
	if asset3d != nil {
		if err := space.SetAsset3D(asset3d, true); err != nil {
			err := errors.Errorf("Node: apiUpdateSpace: failed to update 3d asset: %s", asset3d.GetID())
			api.AbortRequest(c, http.StatusInternalServerError, "space_asset_3d", err, n.log)
			return
		}
	}

	if inBody.SpaceName != "" {
		if err := space.SetName(inBody.SpaceName, true); err != nil {
			err := errors.WithMessagef(err, "Node: apiUpdateSpace: failed to set space name: %s", inBody.SpaceName)
			api.AbortRequest(c, http.StatusInternalServerError, "space_name", err, n.log)
			return
		}
	}

	if err := space.Update(false); err != nil {
		err = errors.WithMessage(err, "Node: apiUpdateSpace: failed to update space")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_space_update", err, n.log)
		return
	}

	// TODO: output full space data
	type Out struct {
		SpaceID string `json:"space_id"`
	}
	out := Out{
		SpaceID: spaceID.String(),
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Set space sub option by space ID
// @Schemes
// @Description Sets a space sub option by space ID, returns appended space option
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param body body node.apiSpacesSetSpaceSubOption.Body true "body params"
// @Success 202 {object} dto.SpaceSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/options/sub [post]
func (n *Node) apiSpacesSetSpaceSubOption(c *gin.Context) {
	type Body struct {
		SubOptionKey   string `json:"sub_option_key" binding:"required"`
		SubOptionValue any    `json:"sub_option_value" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSpacesSetSpaceSubOption: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesSetSpaceSubOption: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesSetSpaceSubOption: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
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

	if _, err := space.SetOptions(modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesSetSpaceSubOption: failed to set options")
		api.AbortRequest(c, http.StatusInternalServerError, "set_options_failed", err, n.log)
		return
	}

	out := dto.SpaceSubOptions{
		inBody.SubOptionKey: inBody.SubOptionValue,
	}

	c.JSON(http.StatusAccepted, out)
}

// @Summary Delete space sub option by space ID
// @Schemes
// @Description Deletes a space sub option by space ID
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param body body node.apiSpacesRemoveSpaceSubOption.Body true "body params"
// @Success 200 {object} nil
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/options/sub [delete]
func (n *Node) apiSpacesRemoveSpaceSubOption(c *gin.Context) {
	type Body struct {
		SubOptionKey string `json:"sub_option_key" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiSpacesRemoveSpaceSubOption: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesRemoveSpaceSubOption: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesRemoveSpaceSubOption: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.ObjectOptions) (*entry.ObjectOptions, error) {
		if current == nil || current.Subs == nil {
			return current, nil
		}

		delete(current.Subs, inBody.SubOptionKey)

		return current, nil
	}

	if _, err := space.SetOptions(modifyFn, true); err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesRemoveSpaceSubOption: failed to set options")
		api.AbortRequest(c, http.StatusInternalServerError, "set_options_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, nil)
}

// @Summary Get space options by space ID
// @Schemes
// @Description Returns a space options based on space ID and query
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param query query node.apiSpacesGetSpaceOptions.InQuery false "query params"
// @Success 200 {object} dto.SpaceOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/options [get]
func (n *Node) apiSpacesGetSpaceOptions(c *gin.Context) {
	type InQuery struct {
		Effective bool `form:"effective"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceOptions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceOptions: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesGetSpaceOptions: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	var out dto.SpaceOptions
	if inQuery.Effective {
		out = space.GetEffectiveOptions()
	} else {
		out = space.GetOptions()
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Get space sub options
// @Schemes
// @Description Returns a space sub options based on query
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Object ID"
// @Param query query node.apiSpacesGetSpaceSubOptions.InQuery true "query params"
// @Success 200 {object} dto.SpaceSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/options/sub [get]
func (n *Node) apiSpacesGetSpaceSubOptions(c *gin.Context) {
	type InQuery struct {
		Effective    bool   `form:"effective"`
		SubOptionKey string `form:"sub_option_key" binding:"required"`
	}
	inQuery := InQuery{
		Effective: true,
	}

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceSubOptions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceSubOptions: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetObjectFromAllObjects(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesGetSpaceSubOptions: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	var options *entry.ObjectOptions
	if inQuery.Effective {
		options = space.GetEffectiveOptions()
	} else {
		options = space.GetOptions()
	}

	if options == nil {
		err := errors.Errorf("Node: apiSpacesGetSpaceSubOptions: empty options")
		api.AbortRequest(c, http.StatusNotFound, "empty_options", err, n.log)
		return
	}

	if options.Subs == nil {
		err := errors.Errorf("Node: apiSpacesGetSpaceSubOptions: empty sub options")
		api.AbortRequest(c, http.StatusNotFound, "empty_sub_options", err, n.log)
		return
	}

	out := dto.SpaceSubOptions{
		inQuery.SubOptionKey: options.Subs[inQuery.SubOptionKey],
	}

	c.JSON(http.StatusOK, out)
}
