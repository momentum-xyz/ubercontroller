package node

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

var store = generic.NewSyncMap[uuid.UUID, string](0)

// @Summary Mint Odyssey for given wallet
// @Schemes
// @Description Returns job_id
// @Tags drive
// @Accept json
// @Produce json
// @Param body node.apiDriveMintOdyssey.Body false
// @Success 200 {object} node.apiDriveMintOdyssey.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/drive/mint-odyssey [post]
func (n *Node) apiDriveMintOdyssey(c *gin.Context) {

	type Body struct {
		BlockHash string `json:"block_hash" binding:"required"`
		Meta      any    `json:"meta" binding:"required"`
	}
	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiDriveMintOdyssey: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	jobID := uuid.New()

	store.Store(jobID, "in progress")

	go mint(jobID)

	type Out struct {
		JobID uuid.UUID `json:"job_id"`
	}
	out := Out{
		JobID: jobID,
	}

	c.JSON(http.StatusOK, out)
}

func mint(jobID uuid.UUID) {
	time.Sleep(time.Second * 30)
	store.Store(jobID, "done")
}

// @Summary Get Mint Odyssey Job ID
// @Schemes
// @Description Returns Mint Odyssey Job ID status
// @Tags drive
// @Accept json
// @Produce json
// @Param job_id path string true "Job ID"
// @Param query query node.apiGetSpace.InQuery false "query params"
// @Success 202 {object} dto.Space
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{job_id} [get]
func (n *Node) apiDriveMintOdysseyCheckJob(c *gin.Context) {
	type Out struct {
		Status string    `json:"status"`
		JobID  uuid.UUID `json:"job_id"`
	}

	jobID, err := uuid.Parse(c.Param("jobID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiDriveMintOdysseyCheckJob: failed to parse uuid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_param", err, n.log)
		return
	}

	status, ok := store.Load(jobID)
	if !ok {
		status = "job not found"
	}

	out := Out{
		JobID:  jobID,
		Status: status,
	}

	c.JSON(http.StatusOK, out)
}
