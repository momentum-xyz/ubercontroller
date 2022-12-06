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
	"github.com/momentum-xyz/ubercontroller/utils"
)

type NodeJSOut struct {
	Data  any      `json:"data"`
	Logs  []string `json:"logs"`
	Error *string  `json:"error"`
}

type NodeJSOutData struct {
	UserID uuid.UUID `json:"userID"`
	Name   string    `json:"name"`
	Image  string    `json:"image"`
}

type StoreItem struct {
	Status    string
	NodeJSOut *NodeJSOut
	Error     error
}

type NFTMeta struct {
	Name  string `json:"name" binding:"required"`
	Image string `json:"image" binding:"required"`
}

type WalletMeta struct {
	Wallet   string
	UserID   uuid.UUID
	Username string
	Avatar   string
}

const StatusInProgress = "in progress"
const StatusDone = "done"
const StatusFailed = "failed"

var log = logger.L()
var store = generic.NewSyncMap[uuid.UUID, StoreItem](0)

// @Summary Get wallet metadata
// @Schemes
// @Description Returns a metadata related to wallet
// @Tags drive
// @Accept json
// @Produce json
// @Param query query node.apiGetWalletMeta.InQuery true "query params"
// @Success 200 {object} node.WalletMeta
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/drive/wallet-meta [get]
func (n *Node) apiGetWalletMeta(c *gin.Context) {
	type InQuery struct {
		Wallet string `form:"wallet" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGetWalletMeta: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	meta, err := n.getWalletMetadata(inQuery.Wallet)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetWalletMeta: failed to get wallet meta")
		api.AbortRequest(c, http.StatusBadRequest, "get_meta_failed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, meta)
}

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
		BlockHash string  `json:"block_hash" binding:"required"`
		Wallet    string  `json:"wallet" binding:"required"`
		Meta      NFTMeta `json:"meta" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiDriveMintOdyssey: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	jobID := uuid.New()

	store.Store(jobID, StoreItem{
		Status:    StatusInProgress,
		NodeJSOut: nil,
	})

	go n.mint(jobID, inBody.Wallet, inBody.Meta, inBody.BlockHash)

	type Out struct {
		JobID uuid.UUID `json:"job_id"`
	}
	out := Out{
		JobID: jobID,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) mint(jobID uuid.UUID, wallet string, meta NFTMeta, blockHash string) {
	// node src/mint.js 5FHneW46xGXgs5mUiveU4sbTyGBzmstUspZC92UhjJM694ty //Alice '{"name":"Test Name", "image":"link"}'

	item := StoreItem{
		Status:    "",
		NodeJSOut: nil,
		Error:     nil,
	}

	b, err := json.Marshal(meta)
	if err != nil {
		err = errors.WithMessage(err, "failed to json.Marshal meta to nodejs in")
		{
			item.Status = StatusFailed
			item.Error = err
			store.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	output, err := exec.Command("node", "mint.js", wallet, "//Alice", string(b), blockHash).Output()
	if err != nil {
		err = errors.WithMessage(err, "failed to exec node mint.js")
		{
			item.Status = StatusFailed
			item.Error = err
			store.Store(jobID, item)
		}
		log.Error(err)
		return
	}
	fmt.Println(string(output))

	var nodeJSOut NodeJSOut
	err = json.Unmarshal(output, &nodeJSOut)
	if err != nil {
		err = errors.WithMessage(err, "failed to json.Unmarshal nodejs out")
		{
			item.Status = StatusFailed
			item.Error = err
			store.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	var data NodeJSOutData
	err = utils.MapDecode(nodeJSOut.Data, &data)
	if err != nil {
		err = errors.WithMessage(err, "failed to utils.MapDecode data to NodeJSOutData")
		{
			item.Status = StatusFailed
			item.Error = err
			store.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	item.NodeJSOut = &nodeJSOut

	if nodeJSOut.Error != nil {
		{
			item.Status = StatusFailed
			item.Error = errors.New("NodeJS mint script returned logic error:" + *nodeJSOut.Error)
			item.NodeJSOut = &nodeJSOut
			store.Store(jobID, item)
		}
		log.Error(err)
		return
	}

	//world, err := n.GetWorlds().CreateWorld(nodeJSOut.Data.UserID)
	//if err != nil {
	//	err = errors.WithMessage(err, "failed to CreateWorld")
	//	{
	//		item.Status = StatusFailed
	//		item.Error = err
	//		store.Store(jobID, item)
	//	}
	//	log.Error(err)
	//	return
	//}
	//
	//err = n.GetWorlds().AddWorld(world, true)
	//if err != nil {
	//	err = errors.WithMessage(err, "failed to AddWorld")
	//	{
	//		item.Status = StatusFailed
	//		item.Error = err
	//		store.Store(jobID, item)
	//	}
	//	log.Error(err)
	//	return
	//}

	item.Status = StatusDone
	item.NodeJSOut = &nodeJSOut
	store.Store(jobID, item)

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
		Error     *string    `json:"error"`
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

	var message *string
	if item.Error != nil {
		e := item.Error.Error()
		message = &e
	}

	out := Out{
		JobID:     jobID,
		Status:    item.Status,
		NodeJSOut: item.NodeJSOut,
		Error:     message,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) getWalletMetadata(wallet string) (*WalletMeta, error) {
	output, _ := exec.Command("node", "./nodejs/check-nft/check-nft.js", wallet).Output()
	var nodeJSOut NodeJSOut
	if err := json.Unmarshal(output, &nodeJSOut); err != nil {
		return nil, errors.WithMessage(err, "failed to unmarshal output")
	}

	if nodeJSOut.Error != nil {
		return nil, errors.New(*nodeJSOut.Error)
	}

	data := utils.GetFromAny(nodeJSOut.Data, []any{})
	if len(data) != 4 {
		return nil, errors.Errorf("invalid data: len %d != 4", len(data))
	}

	userID, err := uuid.Parse(utils.GetFromAny(data[0], ""))
	if err != nil {
		return nil, errors.WithMessage(err, "failed to parse user id")
	}

	meta := &WalletMeta{
		Wallet:   wallet,
		UserID:   userID,
		Username: utils.GetFromAny(data[1], ""),
		Avatar:   utils.GetFromAny(data[3], ""),
	}

	return meta, nil
}
