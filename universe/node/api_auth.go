package node

import (
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
	"net/http"
)

// @Summary Generate auth challenge
// @Schemes
// @Description Returns a new generated challenge based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param query query node.apiGenChallenge.InQuery true "query params"
// @Success 200 {object} node.apiGenChallenge.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/challenge [get]
func (n *Node) apiGenChallenge(c *gin.Context) {
	type InQuery struct {
		Wallet string `form:"wallet" binding:"required"`
	}
	var inQuery InQuery

	if err := c.ShouldBindQuery(&inQuery); err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to bind query")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_query", err, n.log)
		return
	}

	challenge, err := api.GenerateChallenge(inQuery.Wallet)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to generate challenge")
		api.AbortRequest(c, http.StatusInternalServerError, "challenge_generation_failed", err, n.log)
		return
	}

	challengesKey := universe.Attributes.Kusama.ChallengeStore.Key
	modifyFn := func(current *entry.AttributePayload) (*entry.AttributePayload, error) {
		if current == nil {
			current = entry.NewAttributePayload(nil, nil)
		}
		if current.Value == nil {
			current.Value = entry.NewAttributeValue()
		}

		challenges := utils.GetFromAnyMap(*current.Value, challengesKey, make(map[string]any))
		challenges[inQuery.Wallet] = challenge

		// store challenges because we don't know where we got it from
		(*current.Value)[challengesKey] = challenges

		return current, nil
	}

	if _, err := n.UpsertNodeAttribute(
		entry.NewAttributeID(universe.GetKusamaPluginID(), universe.Attributes.Kusama.ChallengeStore.Name),
		modifyFn, true,
	); err != nil {
		err := errors.WithMessage(err, "Node: apiGenChallenge: failed to upsert node attribute")
		api.AbortRequest(c, http.StatusInternalServerError, "node_attribute_upsert_failed", err, n.log)
		return
	}

	type Out struct {
		Challenge string `json:"challenge"`
	}
	out := Out{
		Challenge: challenge,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Generate auth token
// @Schemes
// @Description Returns a new generated token based on params
// @Tags auth
// @Accept json
// @Produce json
// @Param body body node.apiGenToken.InBody true "body params"
// @Success 200 {object} node.apiGenToken.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/auth/token [post]
func (n *Node) apiGenToken(c *gin.Context) {
	type InBody struct {
		Wallet          string `json:"wallet" binding:"required"`
		SignedChallenge string `json:"signedChallenge" binding:"required"`
	}
	var inBody InBody

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	value, ok := n.GetNodeAttributeValue(
		entry.NewAttributeID(universe.GetKusamaPluginID(), universe.Attributes.Kusama.ChallengeStore.Name),
	)
	if !ok {
		err := errors.Errorf("Node: apiGenToken: node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	var challenge string
	if value != nil {
		if store, ok := (*value)[universe.Attributes.Kusama.ChallengeStore.Key]; ok {
			store := utils.GetFromAny(store, (map[string]any)(nil))
			challenge = utils.GetFromAnyMap(store, inBody.Wallet, "")
		}
	}
	if challenge == "" {
		err := errors.Errorf("Node: apiGenToken: challenge not found")
		api.AbortRequest(c, http.StatusNotFound, "challenge_not_found", err, n.log)
		return
	}

	valid, err := api.VerifyPolkadotSignature(inBody.Wallet, challenge, inBody.SignedChallenge)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGenToken: failed to verify polkadot signature")
		api.AbortRequest(c, http.StatusBadRequest, "signature_validation_failed", err, n.log)
		return
	}

	if !valid {
		err := errors.Errorf("Node: apiGenToken: invalid signature")
		api.AbortRequest(c, http.StatusForbidden, "invalid_signature", err, n.log)
		return
	}

	type Out struct {
		Token string `json:"token"`
	}

	out := Out{
		Token: "my super secret token",
	}

	c.JSON(http.StatusOK, out)
}
