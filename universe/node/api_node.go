package node

import (
	"net/http"

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
// @Success 200 {object} node.apiNodeGetChallenge.ChallengeResponse
// @Failure 400 {object} api.HTTPError
// /api/v4/node/get-challenge [post]
func (n *Node) apiNodeGetChallenge(c *gin.Context) {
	type Body struct {
		OdysseyID string `json:"odyssey_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	type ChallengeResponse struct {
		Challenge string `json:"challenge"`
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

	var nodeKeyPair universe.NodeKeyPair
	allowedUserIDs := utils.GetFromAnyMap(*hostingAllowValue, "users", []string{})
	if err := utils.MapDecode(*nodeKeyValue, nodeKeyPair); err != nil {
		err = errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to decode node key pair")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_decode", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	nodeID := n.GetID().String()
	signature, err := GetSignature(nodeKeyPair, nodeID, inBody.OdysseyID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get signature")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_signature", err, n.log)
		return
	}

	if !n.cfg.Common.HostingAllowAll || !utils.Contains(allowedUserIDs, userID.String()) {
		err := errors.New("Node: apiNodeGetChallenge: allow list does not contain user id")
		api.AbortRequest(c, http.StatusBadRequest, "user_not_allowed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, ChallengeResponse{Challenge: hexutil.Encode(signature)})
}

func GetSignature(nodeKeyPair universe.NodeKeyPair, nodeID string, odysseyID string) ([]byte, error) {
	privateKeyBytes, err := hexutil.Decode(nodeKeyPair.PrivateKey)
	if err != nil {
		return nil, err
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}

	message := []byte(nodeID + ":" + odysseyID)
	messageHash := crypto.Keccak256Hash(message)

	signature, err := crypto.Sign(messageHash.Bytes(), privateKey)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
