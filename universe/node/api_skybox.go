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
	"github.com/momentum-xyz/ubercontroller/universe/auth"
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

type SkyboxStatus struct {
	Id            int         `json:"id"`
	Message       *string     `json:"message"`
	Status        string      `json:"status"`
	FileUrl       string      `json:"file_url"`
	ThumbUrl      string      `json:"thumb_url"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	ErrorMessage  interface{} `json:"error_message"`
	SkyboxStyleId int         `json:"skybox_style_id"`
}

var stylesCache = StylesCache{
	updated: time.Time{},
	value:   nil,
}

var skyboxIDToUserID = make(map[int]umid.UMID)
var skyboxIDToWorldID = make(map[int]umid.UMID)

func (s *SkyboxStatus) ToMap() map[string]any {
	m := make(map[string]any)
	m["id"] = s.Id
	m["message"] = s.Message
	m["status"] = s.Status
	m["file_url"] = s.FileUrl
	m["thumb_url"] = s.ThumbUrl
	layout := "2006-01-02T15:04:05Z0700"
	m["created_at"] = s.CreatedAt.Format(layout)
	m["updated_at"] = s.UpdatedAt.Format(layout)
	m["error_message"] = s.ErrorMessage
	m["skybox_style_id"] = s.SkyboxStyleId

	return m
}

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
		apiKey, _, err := n.getApiKeyAndSecret()
		if err != nil {
			err := errors.WithMessage(err, "Node: apiGetSkyboxStyles: failed to getApiKeyAndSecret")
			api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
			return
		}

		url := "https://backend.blockadelabs.com/api/v1/skybox/styles"
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			err := errors.New("Node: apiGetSkyboxStyles: failed to create request to blockadelabs API")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		req.Header.Set("x-api-key", *apiKey)
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

	attrType, ok := n.GetAttributeTypes().GetAttributeType(
		entry.AttributeTypeID{PluginID: universe.GetSystemPluginID(), Name: "skybox_ai"})
	if !ok {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attributeID := entry.AttributeID{
		PluginID: universe.GetSystemPluginID(),
		Name:     "skybox_ai",
	}

	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, inBody.WorldID, userID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	apiKey, _, err := n.getApiKeyAndSecret()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to getApiKeyAndSecret")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	apiUrl := "https://backend.blockadelabs.com/api/v1/skybox?api_key=" + *apiKey

	form := url.Values{}
	form.Add("skybox_style_id", strconv.Itoa(inBody.SkyboxStyleID))
	form.Add("prompt", inBody.Prompt)
	form.Add("webhook_url", fmt.Sprintf("%s/api/v%d/webhook/skybox-blockadelabs", n.cfg.Settings.FrontendURL, ubercontroller.APIMajorVersion))

	n.log.Info("apiPostSkyboxGenerate: form: ", form)

	r, err := http.PostForm(apiUrl, form)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to send post request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)

	if err != nil {
		err := errors.New("Node: apiPostSkyboxGenerate: failed to read blockadelabs API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	response := SkyboxStatus{}

	err = json.Unmarshal(body, &response)
	if err != nil {
		err := errors.New("Node: apiPostSkyboxGenerate: failed to Unmarshal blockadelabs API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	if response.Message != nil {
		// if blockade labs server has internal error it return only 'message' field with error string
		err := errors.New("Node: apiPostSkyboxGenerate: blockadelabs API response error: " + *response.Message)
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	skyboxIDToUserID[response.Id] = userID
	skyboxIDToWorldID[response.Id] = inBody.WorldID

	var modifyFunc modify.Fn[entry.AttributePayload]
	modifyFunc = func(payload *entry.AttributePayload) (*entry.AttributePayload, error) {
		if payload == nil {
			payload = &entry.AttributePayload{
				Value:   &entry.AttributeValue{},
				Options: nil,
			}
		}
		val := *payload.Value

		val[strconv.Itoa(response.Id)] = response.ToMap()

		return payload, nil
	}

	_, err = n.objectUserAttributes.Upsert(objectUserAttributeID, modifyFunc, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiPostSkyboxGenerate: failed to upsert attribute skybox_ai")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	type Out struct {
		Success bool         `json:"success"`
		Data    SkyboxStatus `json:"data"`
	}
	out := Out{
		Success: true,
		Data:    response,
	}

	c.JSON(http.StatusOK, out)
}

func (n *Node) apiPostSkyboxWebHook(c *gin.Context) {
	var inBody SkyboxStatus
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiPostSkyboxWebHook: failed to bind json")
		n.log.Error(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	if inBody.Id == 0 {
		err := errors.New("Node: apiPostSkyboxWebHook: body missing required 'id' field")
		n.log.Error(err)
	}

	_, ok := skyboxIDToWorldID[inBody.Id]
	if !ok {
		err := errors.New("Node: apiPostSkyboxWebHook: no world_id for given 'id'")
		n.log.Error(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	_, ok = skyboxIDToUserID[inBody.Id]
	if !ok {
		err := errors.New("Node: apiPostSkyboxWebHook: no user_id for given 'id'")
		n.log.Error(err)
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "skybox_ai")
	objectUserAttributeID := entry.ObjectUserAttributeID{
		AttributeID: attrID,
		ObjectID:    skyboxIDToWorldID[inBody.Id],
		UserID:      skyboxIDToUserID[inBody.Id],
	}

	var modifyFunc modify.Fn[entry.AttributePayload]
	modifyFunc = func(payload *entry.AttributePayload) (*entry.AttributePayload, error) {
		val := *payload.Value

		val[strconv.Itoa(inBody.Id)] = inBody.ToMap()

		return payload, nil
	}
	_, err := n.objectUserAttributes.Upsert(objectUserAttributeID, modifyFunc, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSkyboxStyles: failed to upsert attribute skybox_ai")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, "ok")
}

// @Summary Delete skybox by ID
// @Schemes
// @Description Delete skybox by ID
// @Tags skybox
// @Accept json
// @Produce json
// @Param body body node.apiRemoveSkyboxByID.Body true "body params"
// @Param skyboxID path string true "SkyboxID int"
// @Success 200 {object} int
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/skybox/{skyboxID} [delete]
func (n *Node) apiRemoveSkyboxByID(c *gin.Context) {

	id, err := strconv.ParseInt(c.Param("skyboxID"), 10, 32)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to parse skyboxID")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, n.log)
		return
	}

	type Body struct {
		WorldID umid.UMID `json:"world_id" binding:"required"`
	}

	var inBody Body
	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	attrType, ok := n.GetAttributeTypes().GetAttributeType(
		entry.AttributeTypeID{PluginID: universe.GetSystemPluginID(), Name: "skybox_ai"})
	if !ok {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to parse user umid")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_user_id", err, n.log)
		return
	}

	attributeID := entry.AttributeID{
		PluginID: universe.GetSystemPluginID(),
		Name:     "skybox_ai",
	}

	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, inBody.WorldID, userID)

	allowed, err := auth.CheckAttributePermissions(
		c, *attrType.GetEntry(), n.GetObjectUserAttributes(), objectUserAttributeID, userID,
		auth.WriteOperation)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: permissions check")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_permissions_check", err, n.log)
		return
	} else if !allowed {
		err := fmt.Errorf("operation not permitted")
		api.AbortRequest(c, http.StatusForbidden, "operation_not_permitted", err, n.log)
		return
	}

	attr, ok := n.GetObjectUserAttributes().GetValue(objectUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiRemoveSkyboxByID: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	apiKey, _, err := n.getApiKeyAndSecret()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to getApiKeyAndSecret")
		api.AbortRequest(c, http.StatusNotFound, "node_attribute_not_found", err, n.log)
		return
	}

	skybox, ok := (*attr)[strconv.Itoa(int(id))]
	if !ok {
		err := errors.Errorf("Node: apiRemoveSkyboxByID: skybox with id not found: %d", id)
		api.AbortRequest(c, http.StatusNotFound, "skybox_not_found", err, n.log)
		return
	}
	s := skybox.(map[string]any)
	statusValue := s["status"].(string)

	if statusValue != "complete" {
		url := "https://backend.blockadelabs.com/api/v1/imagine/requests/" + strconv.Itoa(int(id))
		req, err := http.NewRequest(http.MethodDelete, url, nil)
		if err != nil {
			err := errors.New("Node: apiRemoveSkyboxByID: failed to create request to blockadelabs API")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		req.Header.Set("x-api-key", *apiKey)
		client := http.Client{
			Timeout: 20 * time.Second,
		}

		res, err := client.Do(req)
		if err != nil {
			err := errors.New("Node: apiRemoveSkyboxByID: failed to send request to blockadelabs API")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			err := errors.New("Node: apiRemoveSkyboxByID: failed to read blockadelabs API response")
			api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
			return
		}
		_ = resBody
	}

	url := "https://backend.blockadelabs.com/api/v1/imagine/deleteImagine/" + strconv.Itoa(int(id))
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		err := errors.New("Node: apiRemoveSkyboxByID: failed to create request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	req.Header.Set("x-api-key", *apiKey)
	client := http.Client{
		Timeout: 20 * time.Second,
	}

	res, err := client.Do(req)
	if err != nil {
		err := errors.New("Node: apiRemoveSkyboxByID: failed to send request to blockadelabs API")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	resBody, err := io.ReadAll(res.Body)
	if err != nil {
		err := errors.New("Node: apiRemoveSkyboxByID: failed to read blockadelabs API response")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
	_ = resBody

	var modifyFunc modify.Fn[entry.AttributePayload]
	modifyFunc = func(payload *entry.AttributePayload) (*entry.AttributePayload, error) {
		val := *payload.Value
		delete(val, strconv.Itoa(int(id)))

		return payload, nil
	}
	_, err = n.objectUserAttributes.Upsert(objectUserAttributeID, modifyFunc, true)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiRemoveSkyboxByID: failed to upsert attribute skybox_ai")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.JSON(http.StatusOK, id)
}

func (n *Node) getApiKeyAndSecret() (*string, *string, error) {
	attrID := entry.NewAttributeID(universe.GetSystemPluginID(), "blockadelabs")
	attr, ok := n.nodeAttributes.GetValue(attrID)
	if !ok {
		err := errors.New("'blockadelabs' node attribute not found")
		return nil, nil, err
	}

	if attr == nil {
		err := errors.New("'blockadelabs' node attribute is nul")
		return nil, nil, err
	}
	apiKey := utils.GetFromAnyMap(*attr, "api_key", "")
	secret := utils.GetFromAnyMap(*attr, "secret", "")

	return &apiKey, &secret, nil
}

// @Summary Get skybox image by ID
// @Schemes
// @Description Return Get skybox image by ID
// @Tags skybox
// @Accept json
// @Produce image/jpeg,json
// @Param skyboxID path string true "SkyboxID int"
// @Success 200 {object} nil
// @Failure 400 {object} api.HTTPError
// @Failure 404 {object} api.HTTPError
// @Router /api/v4/skybox/{skyboxID} [get]
func (n *Node) apiGetSkyboxByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("skyboxID"), 10, 32)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSkyboxByID: failed to parse skyboxID")
		api.AbortRequest(c, http.StatusInternalServerError, "invalid_uuid_parse", err, n.log)
		return
	}

	userID, err := api.GetUserIDFromContext(c)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiGetSkyboxByID: failed to get user umid from context")
		api.AbortRequest(c, http.StatusInternalServerError, "get_user_id_failed", err, n.log)
		return
	}
	_ = userID

	attributeID := entry.AttributeID{
		PluginID: universe.GetSystemPluginID(),
		Name:     "skybox_ai",
	}

	worldID, ok := skyboxIDToWorldID[int(id)]
	if !ok {
		err := errors.New("Node: apiGetSkyboxByID: no world_id for given 'id'")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	objectUserAttributeID := entry.NewObjectUserAttributeID(attributeID, worldID, userID)

	attr, ok := n.GetObjectUserAttributes().GetValue(objectUserAttributeID)
	if !ok {
		err := errors.Errorf("Node: apiGetSkyboxByID: object attribute value not found: %s", attributeID)
		api.AbortRequest(c, http.StatusNotFound, "attribute_not_found", err, n.log)
		return
	}

	skybox, ok := (*attr)[strconv.Itoa(int(id))]
	if !ok {
		err := errors.Errorf("Node: apiGetSkyboxByID: skybox with id not found: %d", id)
		api.AbortRequest(c, http.StatusNotFound, "skybox_not_found", err, n.log)
		return
	}

	s := skybox.(map[string]any)
	statusValue := s["status"].(string)
	if statusValue != "complete" {
		err := errors.New("Node: apiGetSkyboxByID: image not generated yet. Current status: " + statusValue)
		api.AbortRequest(c, http.StatusNotFound, "image_not_found", err, n.log)
		return
	}

	fileUrl := s["file_url"].(string)

	resp, err := http.Get(fileUrl)
	defer resp.Body.Close()

	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		err := errors.New("Node: apiGetSkyboxByID: failed to read response body")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}

	c.Writer.Header().Set("Content-Type", "image/jpeg")
	//c.Writer.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))
	_, err = c.Writer.Write(fileBytes)
	if err != nil {
		err := errors.New("Node: apiGetSkyboxByID: failed to write fileBytes")
		api.AbortRequest(c, http.StatusInternalServerError, "internal_error", err, n.log)
		return
	}
}
