package node

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/common/api"
	"github.com/momentum-xyz/ubercontroller/universe/common/api/dto"
)

// @Summary Uploads an image to the media manager
// @Schemes
// @Description Sends an image file to the media manager and returns a hash
// @Tags media
// @Accept json
// @Produce json
// @Success 200 {object} dto.HashResponse
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/upload-image [post]
func (n *Node) apiMediaUploadImage(c *gin.Context) {
	imageFile, err := c.FormFile("file")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to read file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, n.log)
		return
	}

	openedFile, err := imageFile.Open()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to open file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_open", err, n.log)
		return
	}

	defer openedFile.Close()

	req, err := http.NewRequest("POST", n.cfg.Common.RenderInternalURL+"/render/addimage", openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to create post request")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_create_request", err, n.log)
		return
	}

	req.Header.Set("Content-Type", imageFile.Header.Get("Content-Type"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to post data to media-manager")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_post_request", err, n.log)
		return
	}

	defer resp.Body.Close()

	response := dto.HashResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to decode json into response")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_decode", err, n.log)
		return
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Renders name as image, returns media manager hash
// @Schemes
// @Description Sends name to the media manager and returns a hash
// @Tags media
// @Accept json
// @Produce json
// @Param body body node.apiMediaRenderName.inBody true "body params"
// @Success 200 {object} dto.HashResponse
// @Failure 500 {object} api.HTTPError
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/render-name [post]
func (n *Node) apiMediaRenderName(c *gin.Context) {
	type InBody struct {
		Text string `json:"text" binding:"required"`
	}

	inBody := InBody{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		err = errors.WithMessage(err, "Node: apiMediaRenderText: failed to bind json")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_request_body", err, n.log)
		return
	}

	renderRecipe := dto.BackgroundRecipe{
		Background: []int{0, 0, 0, 0},
		Color:      []int{0, 255, 0, 0},
		Thickness:  0,
		Width:      1024,
		Height:     64,
		X:          0,
		Y:          0,
		Text: dto.TextRecipe{
			String:    inBody.Text,
			FontFile:  "",
			FontSize:  0,
			FontColor: []int{220, 220, 200, 255},
			Wrap:      false,
			PadX:      0,
			PadY:      1,
			AlignH:    "center",
			AlignV:    "center",
		},
	}

	buffer := new(bytes.Buffer)
	if err := json.NewEncoder(buffer).Encode(&renderRecipe); err != nil {
		err := errors.WithMessage(err, "Node: apiMediaRenderText: failed to decode json into response")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_decode", err, n.log)
		return
	}

	req, err := http.NewRequest("POST", n.cfg.Common.RenderInternalURL+"/render/addframe", buffer)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaRenderText: failed to create post request")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_create_request", err, n.log)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaRenderText: failed to post data to media-manager")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_post_request", err, n.log)
		return
	}

	defer resp.Body.Close()

	response := dto.HashResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		err := errors.WithMessage(err, "Node: apiMediaRenderText: failed to decode json into response")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_decode", err, n.log)
		return
	}

	c.JSON(http.StatusOK, response)
}
