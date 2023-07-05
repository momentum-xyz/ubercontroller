package node

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type LeonardoResponse struct {
	SdGenerationJob struct {
		GenerationId string `json:"generationId"`
	} `json:"sdGenerationJob"`
}

// @Summary Create a Generation of Images
// @Schemes
// @Description Returns a generation id with which the images can be fetched
// @Tags ai-images
// @Accept json
// @Produce json
// @Param body body node.apiPostGetImageGenerationID.Body true "body params"
// @Success 200 {object} node.apiPostSkyboxGenerate.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/ai-image/generate [post]
func (n *Node) apiPostGetImageGenerationID(c *gin.Context) {
	type Body struct {
		Prompt string `json:"prompt" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if len(inBody.Prompt) > 550 {
		err := errors.New("Node: apiPostGetImageGenerationID: prompt length must be less than 550")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	apiKey, err := n.getApiLeonardoKeyAndSecret()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to getApiKeyAndSecret")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	apiUrl := "https://cloud.leonardo.ai/api/rest/v1/generations"

	jsonData := map[string]string{"prompt": inBody.Prompt}
	reqBody := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(reqBody).Encode(jsonData); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_encode", errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to encode JSON"), n.log)
		return
	}

	req, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to create new request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	req.Header.Set("Authorization", "Bearer "+*apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to send post request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		err := errors.New("Node: apiPostGetImageGenerationID: failed to read leonardo API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	response := LeonardoResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		err := errors.New("Node: apiPostGetImageGenerationID: failed to Unmarshal leonardo API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	type Out struct {
		Success bool             `json:"success"`
		Data    LeonardoResponse `json:"data"`
	}
	out := Out{
		Success: true,
		Data:    response,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) getApiLeonardoKeyAndSecret() (*string, error) {
	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "leonardo")
	attr, ok := n.nodeAttributes.GetValue(attrID)
	if !ok {
		err := errors.New("'leonardo' node attribute not found")
		return nil, err
	}

	if attr == nil {
		err := errors.New("'leonardo' node attribute is nil")
		return nil, err
	}
	apiKey := utils.GetFromAnyMap(*attr, "api_key", "")

	return &apiKey, nil
}
