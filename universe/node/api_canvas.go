package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
)

type LeonardoResponse struct {
	SdGenerationJob struct {
		GenerationId string `json:"generationId"`
	} `json:"sdGenerationJob"`
}

type GeneratedImage struct {
	URL  string `json:"url"`
	NSFW bool   `json:"nsfw"`
	ID   string `json:"id"`
}

type GenerationResponse struct {
	GenerationsByPK struct {
		GeneratedImages []GeneratedImage `json:"generated_images"`
		Prompt          string           `json:"prompt"`
		Status          string           `json:"status"`
		ID              string           `json:"id"`
		CreatedAt       string           `json:"createdAt"`
	} `json:"generations_by_pk"`
}

// @Summary Gets user contributions by object id
// @Description Returns an object with a nested items array of contributions
// @Tags canvas
// @Security Bearer
// @Param object_id path string true "ObjectID string"
// @Success 200 {object} node.apiCanvasGetUserContributions.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/canvas/{object_id}/user-contributions [get]
func (n *Node) apiCanvasGetUserContributions(c *gin.Context) {
	type InQuery struct {
		OrderBy string `form:"order" json:"order"`
		Limit   uint   `form:"limit,default=10" json:"limit"`
		Offset  uint   `form:"offset" json:"offset"`
		Search  string `form:"q" json:"q"`
	}
	var q InQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_query", err, n.log)
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiCanvasGetUserContributions: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, out)
}
