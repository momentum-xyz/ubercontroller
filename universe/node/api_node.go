package node

import (
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

// @Summary Checks if an odyssey can be registered with a node
// @Description Checks if an odyssey can be registered with a node
// @Tags node
// @Security Bearer
// @Param body body node.apiNodeGetChallenge.Body true "body params"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// /api/v4/node/get-challenge [post]
func (n *Node) apiNodeGetChallenge(c *gin.Context) {
	if n.cfg.Common.HostingAllowAll {
		c.JSON(http.StatusOK, nil)
	}

	type Body struct {
		OdysseyID string `json:"odyssey_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	hostingAllowID := entry.NewAttributeID(universe.GetSystemPluginID(), "hosting_allow_list")
	nodeKeyID := entry.NewAttributeID(universe.GetSystemPluginID(), "node_key")
	hostingAllowValue, ok := n.GetNodeAttributes().GetValue(hostingAllowID)
	if !ok || hostingAllowValue == nil {
		err := errors.New("Node: apiNodeGetChallenge: node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}
	nodeKeyValue, ok := n.GetNodeAttributes().GetValue(nodeKeyID)
	if !ok || nodeKeyValue == nil {
		err := errors.New("Node: apiNodeGetChallenge: node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	allowedUserIDs := utils.GetFromAnyMap(*hostingAllowValue, "users", []string{})
	nodeKey := utils.GetFromAny(*nodeKeyValue, "")

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	if !utils.Contains(allowedUserIDs, userID.String()) {
		err := errors.New("Node: apiNodeGetChallenge: allow list does not contain user id")
		api.AbortRequest(c, http.StatusBadRequest, "user_not_allowed", err, n.log)
		return
	}

	parts := strings.Split(nodeKey, ":")
	if len(parts) != 2 {
		err := errors.New("Node: apiNodeGetChallenge: key pair not found or invalid")
		api.AbortRequest(c, http.StatusBadRequest, "pair_not_found", err, n.log)
		return
	}

	privateKeyBytes, err := hexutil.Decode("0x" + parts[1])
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to decode private key")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_decode", err, n.log)
		return
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to parse private key")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_parse", err, n.log)
		return
	}

	nodeID := n.GetID().String()

	message := []byte(nodeID + ":" + inBody.OdysseyID)
	messageHash := crypto.Keccak256Hash(message)

	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	type ChallengeResponse struct {
		Challenge string `json:"challenge"`
	}

	c.JSON(http.StatusOK, ChallengeResponse{Challenge: hexutil.Encode(signature)})
}
