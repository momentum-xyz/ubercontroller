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

// @Summary Get space options
// @Schemes
// @Description Returns a space options based on query
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
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

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
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
// @Param space_id path string true "Space ID"
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

	space, ok := n.GetSpaceFromAllSpaces(spaceID)
	if !ok {
		err := errors.Errorf("Node: apiSpacesGetSpaceSubOptions: space not found: %s", spaceID)
		api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
		return
	}

	var options *entry.SpaceOptions
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
