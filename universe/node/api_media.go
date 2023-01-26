package node

import (
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
// @Router /api/v4/media/upload/image [post]
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

	req, err := http.NewRequest("POST", n.CFG.Common.RenderInternalURL+"/render/addimage", openedFile)
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
