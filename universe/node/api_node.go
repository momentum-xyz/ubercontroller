package node

import (
	"fmt"
	"math/big"
	"net/http"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
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
	nodePrivateKeyID := entry.NewAttributeID(universe.GetSystemPluginID(), "node_private_key")
	hostingAllowValue, ok := n.GetNodeAttributes().GetValue(hostingAllowID)
	if !ok || hostingAllowValue == nil {
		err := errors.New("Node: apiNodeGetChallenge: node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}
	nodePrivateKeyValue, ok := n.GetNodeAttributes().GetValue(nodePrivateKeyID)
	if !ok || nodePrivateKeyValue == nil {
		err := errors.New("Node: apiNodeGetChallenge: node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}
	fmt.Println("hostingAllowValue:", hostingAllowValue)
	fmt.Println("nodePrivateKeyValue:", nodePrivateKeyValue)

	allowedUserIDs := utils.GetFromAnyMap(*hostingAllowValue, "users", []string{})
	privateKey := utils.GetFromAnyMap(*nodePrivateKeyValue, "node_private_key", "")
	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}

	nodeID := n.GetID().String()
	fmt.Println("nodeID:", nodeID)
	signature, err := GetSignature(privateKey, nodeID, inBody.OdysseyID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiNodeGetChallenge: failed to get signature")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_signature", err, n.log)
		return
	}

	if !n.cfg.Common.HostingAllowAll && !utils.Contains(allowedUserIDs, userID.String()) {
		fmt.Println("allowedUserIDs:", allowedUserIDs)
		err := errors.New("Node: apiNodeGetChallenge: allow list does not contain user id: " + userID.String())
		api.AbortRequest(c, http.StatusBadRequest, "user_not_allowed", err, n.log)
		return
	}

	c.JSON(http.StatusOK, ChallengeResponse{Challenge: hexutil.Encode(signature)})
}

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	fmt.Println("msg:", msg, "data len:", len(data))
	return crypto.Keccak256([]byte(msg))
}

func padTo256bits(data []byte) []byte {
	if len(data) >= 32 {
		return data
	}
	pad := make([]byte, 32-len(data))
	return append(pad, data...)
}

func GetSignature(privateKey string, nodeID string, odysseyID string) ([]byte, error) {
	privateKeyBytes, err := hexutil.Decode(privateKey)
	if err != nil {
		return nil, err
	}
	fmt.Printf("privateKeyBytes: %x\n", privateKeyBytes)

	privateECDSAKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	fmt.Println("privateECDSAKey:", privateECDSAKey)

	fmt.Println("nodeID:", nodeID)
	fmt.Println("odysseyID:", odysseyID)

	uuidNodeID, err := umid.Parse(nodeID)
	// uuidObj, err := uuid.Parse(nodeID)
	if err != nil {
		return nil, err
	}
	uuidOdysseyID, err := umid.Parse(odysseyID)
	if err != nil {
		return nil, err
	}

	nodeIDBigInt := new(big.Int)
	nodeIDBigInt.SetBytes(uuidNodeID[:])
	// nodeIDBigInt.SetBytes(uuidObj.Bytes())

	odysseyIDBigInt := new(big.Int)
	// odysseyIDBigInt.SetString(odysseyID, 10)
	odysseyIDBigInt.SetBytes(uuidOdysseyID[:])

	// nodeIDDecStr := nodeIDBigInt.String()
	// nodeIDDecStr := strconv.FormatUint(utils.UMIDToSEQ((nodeID)), 10)
	// fmt.Println("Converted nodeID ", nodeID, "to decimal:", nodeIDDecStr)

	// message := []byte(nodeIDDecStr + ":" + odysseyID)
	// message := []byte(nodeIDDecStr + odysseyID)
	// message := []byte(nodeID + odysseyID)
	// messageHash := crypto.Keccak256Hash(message)

	message := append(padTo256bits(nodeIDBigInt.Bytes()), padTo256bits(odysseyIDBigInt.Bytes())...)
	fmt.Printf("message: %x\n", message)
	hashedMessage := crypto.Keccak256(message)
	fmt.Printf("hashedMessage: %x\n", hashedMessage)
	// signature, err := crypto.Sign(hashedMessage, privateECDSAKey)

	// messageHash := crypto.Keccak256Hash(message)
	// messageHash := signHash(message)

	prefixedMessageHash := signHash(hashedMessage)
	fmt.Printf("prefixedMessageHash: %x\n", prefixedMessageHash)
	signature, err := crypto.Sign(prefixedMessageHash, privateECDSAKey)

	// signature, err := crypto.Sign(messageHash, privateECDSAKey)
	// signature, err := crypto.Sign(messageHash.Bytes(), privateECDSAKey)
	if err != nil {
		return nil, err
	}

	signature[64] += 27
	fmt.Printf("signature: %x\n", signature)
	fmt.Println("signature len:", len(signature), "signature[64]", signature[64])

	return signature, nil
}
