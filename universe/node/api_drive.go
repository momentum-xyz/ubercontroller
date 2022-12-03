package node

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/logger"
	"github.com/momentum-xyz/ubercontroller/types/generic"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
)

type NodeJSOut struct {
	Data  any      `json:"data"`
	Logs  []string `json:"logs"`
	Error *string  `json:"error"`
}

type StoreItem struct {
	Status       string
	NodeJSResult *NodeJSOut
}

const StatusInProgress = "in progress"
const StatusDone = "done"

var log = logger.L()
var store = generic.NewSyncMap[uuid.UUID, StoreItem](0)

// @Summary Mint Odyssey for given wallet
// @Schemes
// @Description Returns job_id
// @Tags drive
// @Accept json
// @Produce json
// @Param body body node.apiDriveMintOdyssey.Body true "body params"
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

	store.Store(jobID, StoreItem{
		Status:       StatusInProgress,
		NodeJSResult: nil,
	})

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
	//time.Sleep(time.Second * 30)
	bob := "5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty"

	output, err := exec.Command("node", "mint.js", bob, "//Alice").Output()
	var nodeJSOut NodeJSOut
	err = json.Unmarshal(output, &nodeJSOut)
	if err != nil {
		log.Error(errors.WithMessage(err, "failed to json.Unmarshal nodejs out"))
	}

	store.Store(jobID, StoreItem{
		Status:       StatusDone,
		NodeJSResult: &nodeJSOut,
	})

	fmt.Println(string(output))
	fmt.Println("***")
}

// @Summary Get Mint Odyssey Job ID
// @Schemes
// @Description Returns Mint Odyssey Job ID status
// @Tags drive
// @Accept json
// @Produce json
// @Param job_id path string true "Job ID"
// @Success 200 {object} node.apiDriveMintOdysseyCheckJob.Out
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/spaces/{job_id} [get]
func (n *Node) apiDriveMintOdysseyCheckJob(c *gin.Context) {
	type Out struct {
		NodeJSOut *NodeJSOut `json:"nodeJSOut"`
		Status    string     `json:"status"`
		JobID     uuid.UUID  `json:"job_id"`
	}

	jobID, err := uuid.Parse(c.Param("jobID"))
	if err != nil {
		err = errors.WithMessage(err, "Node: apiDriveMintOdysseyCheckJob: failed to parse uuid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_param", err, n.log)
		return
	}

	item, ok := store.Load(jobID)
	if !ok {
		item.Status = "job not found"
	}

	out := Out{
		JobID:     jobID,
		Status:    item.Status,
		NodeJSOut: item.NodeJSResult,
	}

	c.JSON(http.StatusOK, out)
}
