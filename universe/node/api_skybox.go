package node

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/momentum-xyz/ubercontroller"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/modify"
	"github.com/momentum-xyz/ubercontroller/utils/umid"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/utils"
)

type StylesCache struct {
	updated time.Time
	value   []StyleItem
}

type StyleItem struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	MaxChar             int    `json:"max-char"`
	NegativeTextMaxChar int    `json:"negative-text-max-char"`
	Image               any    `json:"image"`
	SortOrder           int    `json:"sort_order"`
	Premium             int    `json:"premium"`
	SkyboxStyleFamilies []any  `json:"skybox_style_families"`
}

type CreateSkyboxResponse struct {
	Id                int         `json:"id"`
	Message           *string     `json:"message"`
	ObfuscatedId      string      `json:"obfuscated_id"`
	UserId            int         `json:"user_id"`
	Title             string      `json:"title"`
	Prompt            string      `json:"prompt"`
	Seed              int         `json:"seed"`
	NegativeText      interface{} `json:"negative_text"`
	Username          string      `json:"username"`
	Status            string      `json:"status"`
	QueuePosition     int         `json:"queue_position"`
	FileUrl           string      `json:"file_url"`
	ThumbUrl          string      `json:"thumb_url"`
	DepthMapUrl       string      `json:"depth_map_url"`
	RemixImagineId    interface{} `json:"remix_imagine_id"`
	RemixObfuscatedId interface{} `json:"remix_obfuscated_id"`
	IsMyFavorite      bool        `json:"isMyFavorite"`
	CreatedAt         time.Time   `json:"created_at"`
	UpdatedAt         time.Time   `json:"updated_at"`
	ErrorMessage      interface{} `json:"error_message"`
	PusherChannel     string      `json:"pusher_channel"`
	PusherEvent       string      `json:"pusher_event"`
	Type              string      `json:"type"`
	SkyboxStyleId     int         `json:"skybox_style_id"`
	SkyboxId          int         `json:"skybox_id"`
	SkyboxStyleName   string      `json:"skybox_style_name"`
	SkyboxName        string      `json:"skybox_name"`
}

var stylesCache = StylesCache{
	updated: time.Time{},
	value:   nil,
}

var skyboxIDToUserID = make(map[string]umid.UMID)
var skyboxIDToWorldID = make(map[string]umid.UMID)

// @Summary Get lists the known blockadelabs art styles
// @Schemes
// @Description Return blockadelabs art styles
// @Tags skybox
// @Accept json
// @Produce json
// @Success 200 {object} []node.StyleItem
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/skybox/styles [get]
func (n *Node) apiGetSkyboxStyles(c *gin.Context) {
	agoTime := time.Now().Add(-time.Minute * 10)
	if stylesCache.updated.Before(agoTime) {
		attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "blockadelabs")
		attr, ok := n.nodeAttributes.GetValue(attrID)
		if !ok {
			err := errors.New("Node: apiGetSkyboxStyles: 'blockadelabs' node attribute not found")
			api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
			return
		}

		if attr == nil {
			err := errors.New("Node: apiGetSkyboxStyles: 'blockadelabs' node attribute is nul")
			api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
			return
		}
		apiKey := utils.GetFromAnyMap(*attr, "api_key", "")
		secret := utils.GetFromAnyMap(*attr, "secret", "")
		_ = secret

		url := "https://backend.blockadelabs.com/api/v1/skybox/styles"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			err := errors.New("Node: apiGetSkyboxStyles: failed to create request to blockadelabs API")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		req.Header.Set("x-api-key", apiKey)
		client := http.Client{
			Timeout: 20 * time.Second,
		}

		res, err := client.Do(req)
		if err != nil {
			err := errors.New("Node: apiGetSkyboxStyles: failed to send request to blockadelabs API")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			err := errors.New("Node: apiGetSkyboxStyles: failed to read blockadelabs API response")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		items := make([]StyleItem, 0)

		err = json.Unmarshal(resBody, &items)
		if err != nil {
			err := errors.New("Node: apiGetSkyboxStyles: failed to Unmarshal blockadelabs API response")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		stylesCache = StylesCache{
			updated: time.Now(),
			value:   items,
		}
	}

	c.JSON(http.StatusOK, stylesCache.value)
}

// @Summary Start skybox generation
// @Schemes
// @Description Start skybox generation
// @Tags skybox
// @Accept json
// @Produce json
// @Param body body node.apiPostSkyboxGenerate.Body true "body params"
// @Success 200 {object} node.apiPostSkyboxGenerate.Out
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/skybox/generate [post]
func (n *Node) apiPostSkyboxGenerate(c *gin.Context) {

	type Body struct {
		SkyboxStyleID int       `json:"skybox_style_id" binding:"required"`
		Prompt        string    `json:"prompt" binding:"required"`
		WorldID       umid.UMID `json:"world_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if len(inBody.Prompt) > 550 {
		err := errors.New("Node: apiPostSkyboxGenerate: prompt length must be less than 550")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	//TODO check that world exists and user has permissions

	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "blockadelabs")
	attr, ok := n.nodeAttributes.GetValue(attrID)
	if !ok {
		err := errors.New("Node: apiPostSkyboxGenerate: 'blockadelabs' node attribute not found")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	if attr == nil {
		err := errors.New("Node: apiPostSkyboxGenerate: 'blockadelabs' node attribute is nul")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}
	apiKey := utils.GetFromAnyMap(*attr, "api_key", "")
	secret := utils.GetFromAnyMap(*attr, "secret", "")
	_ = secret

	apiUrl := "https://backend.blockadelabs.com/api/v1/skybox?api_key=" + apiKey

	form := url.Values{}
	form.Add("skybox_style_id", strconv.Itoa(inBody.SkyboxStyleID))
	form.Add("prompt", inBody.Prompt)
	form.Add("webhook_url", fmt.Sprintf("%s/api/v%d/webhook/skybox-blockadelabs", n.cfg.Settings.FrontendURL, ubercontroller.APIMajorVersion))

	n.log.Info("apiPostSkyboxGenerate: form: ", form)

	r, err := http.PostForm(apiUrl, form)
	if err != nil {
		fmt.Println(err)
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		err := errors.New("Node: apiPostSkyboxGenerate: failed to read blockadelabs API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	response := CreateSkyboxResponse{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		err := errors.New("Node: apiGetSkyboxStyles: failed to Unmarshal blockadelabs API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	if response.Message != nil {
		// if blockade labs server has internal error it return only 'message' field with error string
		err := errors.New("Node: apiGetSkyboxStyles: blockadelabs API response error: " + *response.Message)
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	skyboxIDToUserID[response.ObfuscatedId] = userID
	skyboxIDToWorldID[response.ObfuscatedId] = inBody.WorldID

	attrID = entry.NewAttributeID(universe.GetSystemPluginID(), "skybox_ai")
	objectUserAttributeID := entry.ObjectUserAttributeID{
		AttributeID: attrID,
		ObjectID:    inBody.WorldID,
		UserID:      userID,
	}

	m := make(entry.AttributeValue)
	m[response.ObfuscatedId] = response

	var p = entry.AttributePayload{
		Value:   &m,
		Options: nil,
	}

	_, err = n.objectUserAttributes.Upsert(objectUserAttributeID, modify.MergeWith(&p), true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSkyboxStyles: failed to upsert attribute skybox_ai")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	type Out struct {
		Success bool                 `json:"success"`
		Data    CreateSkyboxResponse `json:"data"`
	}
	out := Out{
		Success: true,
		Data:    response,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiPostSkyboxWebHook(c *gin.Context) {
	var inBody CreateSkyboxResponse
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostSkyboxWebHook: failed to bind json")
		n.log.Error(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	fmt.Println(inBody)

	if inBody.ObfuscatedId == "" {
		err := errors.New("Node: apiPostSkyboxWebHook: body missing required 'obfuscated_id' field")
		n.log.Error(err)
	}

	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "skybox_ai")
	objectUserAttributeID := entry.ObjectUserAttributeID{
		AttributeID: attrID,
		ObjectID:    skyboxIDToWorldID[inBody.ObfuscatedId],
		UserID:      skyboxIDToUserID[inBody.ObfuscatedId],
	}

	m := make(entry.AttributeValue)
	m[inBody.ObfuscatedId] = inBody

	var p = entry.AttributePayload{
		Value:   &m,
		Options: nil,
	}

	// TODO Start here merge function fail for some reasons (same way create attribute without error on line 266)
	_, err := n.objectUserAttributes.Upsert(objectUserAttributeID, modify.MergeWith(&p), true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSkyboxStyles: failed to upsert attribute skybox_ai")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, "ok")
}
