package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
	"github.com/momentum-xyz/ubercontroller/universe/api/dto"
)

// @Summary Returns space effective options
// @Schemes
// @Description Returns space effective options
// @Tags spaces
// @Accept json
// @Produce json
// @Param space_id path string true "Space ID"
// @Success 200 {object} dto.SpaceEffectiveOptions
// @Success 500 {object} api.HTTPError
// @Success 400 {object} api.HTTPError
// @Success 404 {object} api.HTTPError
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

	out := dto.SpaceEffectiveOptions{
		spaceID: space.GetEffectiveOptions(),
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiSpacesGetSpaceEffectiveSubOptions(c *gin.Context) {
	panic("how we can get any subfield if our SpaceOptions is a struct?")

	//inQuery := struct {
	//	SubOptionKey string `form:"sub_option_key" binding:"required"`
	//}{}
	//
	//if err := c.ShouldBindQuery(&inQuery); err != nil {
	//	err := errors.WithMessage(err, "Node: apiSpacesGetSpaceEffectiveSubOptions: failed to bind query")
	//	api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
	//	return
	//}
	//
	//spaceID, err := uuid.Parse(c.Param("spaceID"))
	//if err != nil {
	//	err := errors.WithMessage(err, "Node: apiSpacesGetSpaceEffectiveSubOptions: failed to parse space id")
	//	api.AbortRequest(c, http.StatusBadRequest, "invalid_space_id", err, n.log)
	//	return
	//}
	//
	//space, ok := n.GetSpaceFromAllSpaces(spaceID)
	//if !ok {
	//	err := errors.Errorf("Node: apiSpacesGetSpaceEffectiveSubOptions: space not found: %s", spaceID)
	//	api.AbortRequest(c, http.StatusNotFound, "space_not_found", err, n.log)
	//	return
	//}
}
