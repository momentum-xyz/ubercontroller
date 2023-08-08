package node

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
)

// @Summary Uploads an image to the media manager
// @Description Sends an image file to the media manager and returns a hash
// @Tags media
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "image file"
// @Success 200 {object} dto.HashResponse
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

	hash, err := n.media.AddImage(openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadImage: failed to add image")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_image", err, n.log)
		return
	}

	response := dto.HashResponse{
		Hash: hash,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Uploads a video to the media manager
// @Description Sends a video file to the media manager and returns a hash
// @Tags media
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "image file"
// @Success 200 {object} dto.HashResponse
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/upload/video [post]
func (n *Node) apiMediaUploadVideo(c *gin.Context) {
	videoFile, err := c.FormFile("file")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadVideo: failed to read file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, n.log)
		return
	}

	openedFile, err := videoFile.Open()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadVideo: failed to open file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_open", err, n.log)
		return
	}

	defer openedFile.Close()

	hash, err := n.media.AddVideo(openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadVideo: failed to add image")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_image", err, n.log)
		return
	}

	response := dto.HashResponse{
		Hash: hash,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Uploads an audio file to the media manager
// @Description Sends an audio file to the media manager and returns its hash
// @Tags media
// @Security Bearer
// @Accept multipart/form-data
// @Param file formData file true "image file"
// @Success 200 {object} dto.HashResponse
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/upload/audio [post]
func (n *Node) apiMediaUploadAudio(c *gin.Context) {
	audioFile, err := c.FormFile("file")
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadAudio: failed to read file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_read", err, n.log)
		return
	}

	openedFile, err := audioFile.Open()
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadAudio: failed to open file")
		api.AbortRequest(c, http.StatusBadRequest, "failed_to_open", err, n.log)
		return
	}

	defer openedFile.Close()

	hash, err := n.media.AddTrack(openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadAudio: failed to add image")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_image", err, n.log)
		return
	}

	response := dto.HashResponse{
		Hash: hash,
	}

	c.JSON(http.StatusOK, response)
}
