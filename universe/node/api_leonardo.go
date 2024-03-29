package node

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
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

// @Summary Get images by generation id
// @Description Returns an array of images by generation id
// @Tags leonardo
// @Security Bearer
// @Param leonardo_id path string true "LeonardoID string"
// @Success 200 {object} node.apiGetImageGeneration.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/leonardo/generate/{leonardo_id} [get]
func (n *Node) apiGetImageGeneration(c *gin.Context) {
	leonardoID := c.Param("leonardoID")

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetImageGeneration: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	apiKey, err := n.getApiLeonardoKeyAndSecret(nil, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetImageGeneration: failed to getApiKeyAndSecret")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	apiUrl := "https://cloud.leonardo.ai/api/rest/v1/generations/" + leonardoID

	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetImageGeneration: failed to create new request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	req.Header.Set("Authorization", "Bearer "+*apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetImageGeneration: failed to send post request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	response := GenerationResponse{}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", errors.New("Node: apiGetImageGeneration: failed to decode leonardo API response"), n.log)
		return
	}

	r.Body.Close()

	type Out struct {
		Success bool               `json:"success"`
		Data    GenerationResponse `json:"data"`
	}
	out := Out{
		Success: true,
		Data:    response,
	}

	c.JSON(http.StatusOK, out)
}

// @Summary Create a Generation of Images
// @Description Returns a generation id with which the images can be fetched
// @Tags leonardo
// @Security Bearer
// @Param body body node.apiPostImageGenerationID.Body true "body params"
// @Success 200 {object} node.apiPostImageGenerationID.Out
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/leonardo/generate [post]
func (n *Node) apiPostImageGenerationID(c *gin.Context) {
	type Body struct {
		Prompt string `json:"prompt" binding:"required"`
		Model  string `json:"model" binding:"required"`
		//WorldID umid.UMID `json:"world_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostImageGenerationID: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if len(inBody.Prompt) > 550 {
		err := errors.New("Node: apiPostImageGenerationID: prompt length must be less than 550")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostImageGenerationID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	apiKey, err := n.getApiLeonardoKeyAndSecret(nil, userID)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostImageGenerationID: failed to getApiKeyAndSecret")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	apiUrl := "https://cloud.leonardo.ai/api/rest/v1/generations"

	jsonData := map[string]string{
		"prompt":  inBody.Prompt,
		"modelId": inBody.Model,
	}
	reqBody := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(reqBody).Encode(jsonData); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_encode", errors.WithMessage(err, "Node: apiPostGetImageGenerationID: failed to encode JSON"), n.log)
		return
	}

	req, err := http.NewRequest("POST", apiUrl, reqBody)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostImageGenerationID: failed to create new request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	req.Header.Set("Authorization", "Bearer "+*apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	r, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostImageGenerationID: failed to send post request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	response := LeonardoResponse{}
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", errors.New("Node: apiPostGetImageGenerationID: failed to decode leonardo API response"), n.log)
		return
	}

	r.Body.Close()

	type Out struct {
		Success bool             `json:"success"`
		Data    LeonardoResponse `json:"data"`
	}
	out := Out{
		Success: true,
		Data:    response,
	}

	n.trackAIUsage(c, "leonardo", userID)

	c.JSON(http.StatusOK, out)
}

func (n *Node) getApiLeonardoKeyAndSecret(objectID *umid.UMID, userID umid.UMID) (*string, error) {
	var objectAttr, nodeAttr, userAttr *entry.AttributeValue
	var apiKey string

	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "leonardo")

	if objectID != nil {
		object, ok := n.GetObjectFromAllObjects(*objectID)
		if ok {
			objectAttr, _ = object.GetObjectAttributes().GetValue(attrID)
		}
	}
	nodeAttr, _ = n.nodeAttributes.GetValue(attrID)
	userAttributeID := entry.NewUserAttributeID(attrID, userID)
	userAttr, _ = n.GetUserAttributes().GetValue(userAttributeID)

	list := []*entry.AttributeValue{objectAttr, nodeAttr, userAttr}

	for _, attr := range list {
		if attr == nil {
			continue
		}
		if value := utils.GetFromAnyMap(*attr, "api_key", ""); value != "" && apiKey == "" {
			apiKey = value
		}
	}

	if objectAttr == nil && nodeAttr == nil && userAttr == nil {
		err := errors.New("'leonardo' node attribute not found")
		return nil, err
	}

	return &apiKey, nil

}
