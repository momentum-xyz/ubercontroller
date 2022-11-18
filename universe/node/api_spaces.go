package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Set space sub option
// @Schemes
// @Description Sets a space sub option
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
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

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesSetSpaceSubOption: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.SpaceOptions) (*entry.SpaceOptions, error) {
		if current == nil {
			current = &entry.SpaceOptions{}
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

// @Summary Delete space sub option
// @Schemes
// @Description Deletes a space sub option
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
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

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesRemoveSpaceSubOption: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	modifyFn := func(current *entry.SpaceOptions) (*entry.SpaceOptions, error) {
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

// @Summary Get space effective options
// @Schemes
// @Description Returns a space effective options
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Success 200 {object} dto.SpaceEffectiveOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/effective-options [get]
func (n *Node) apiSpacesGetSpaceEffectiveOptions(c *gin.Context) {
	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceEffectiveOptions: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesGetSpaceEffectiveOptions: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	out := dto.SpaceEffectiveOptions(space.GetEffectiveOptions())

	c.JSON(http.StatusOK, out)
}

// @Summary Get space effective sub option
// @Schemes
// @Description Returns a space effective sub option
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Param query query node.apiSpacesGetSpaceEffectiveSubOption.Query true "query params"
// @Success 200 {object} dto.SpaceEffectiveSubOptions
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{space_id}/effective-options/sub [get]
func (n *Node) apiSpacesGetSpaceEffectiveSubOption(c *gin.Context) {
	type Query struct {
		SubOptionKey string `form:"sub_option_key" binding:"required"`
	}

	var inQuery Query
	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceEffectiveSubOption: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	spaceID, err := uuid.Parse(c.Param("spaceID"))
	if err != nil {
		err := errors.WithMessage(err, "Node: apiSpacesGetSpaceEffectiveSubOption: failed to parse space id")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
		return
	}

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesGetSpaceEffectiveSubOption: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	effectiveOptions := space.GetEffectiveOptions()
	if effectiveOptions == nil {
		err := errors.Errorf("Node: apiSpacesGetSpaceEffectiveSubOption: empty effective options")
		api.AbortRequest(c, http.StatusNotFound, "empty_effective_options", err, n.log)
		return
	}

	if effectiveOptions.Subs == nil {
		err := errors.Errorf("Node: apiSpacesGetSpaceEffectiveSubOption: empty effective sub options")
		api.AbortRequest(c, http.StatusNotFound, "empty_effective_sub_options", err, n.log)
		return
	}

	out := dto.SpaceEffectiveSubOptions{
		inQuery.SubOptionKey: effectiveOptions.Subs[inQuery.SubOptionKey],
	}

	c.JSON(http.StatusOK, out)
}
