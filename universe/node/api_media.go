package node

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/logic/api"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
)

func makeTimestamp() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

// @Summary Gets an image from the (internal) media-manager
// @Description Serves a generic image from the (internal) media-manager
// @Tags media
// @Security Bearer
// @Accept json
// @Produce json
// @Param file path string true "image file"
// @Success 200 {object} dto.HashResponse
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/render/get/{file} [get]
func (n *Node) apiMediaGetImage(c *gin.Context) {
	tm1 := makeTimestamp()

	filename := c.Param("file")
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaGetImage: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetImage: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	meta, filepath, err := n.media.GetImage(filename)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetImage: failed to get image")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_image", err, n.log)
		return
	}

	c.Header("Content-Type", meta.Mime)
	c.Header("x-height", strconv.Itoa(meta.H))
	c.Header("x-width", strconv.Itoa(meta.W))
	c.File(*filepath)

	n.log.Infof("Endpoint Hit: Image served: %s %d", filename, makeTimestamp()-tm1)
}

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

// @Summary Gets a video from the media manager
// @Description Serves a video file from the media manager
// @Tags media
// @Security Bearer
// @Accept json
// @Produce json
// @Param file path string true "video file"
// @Success 200 {file} byte "video file"
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/media/render/video/{file} [get]
func (n *Node) apiMediaGetVideo(c *gin.Context) {
	filename := c.Param("file")
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaGetVideo: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetVideo: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	file, fileInfo, contentType, err := n.media.GetVideo(filename)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetVideo: failed to get video")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_video", err, n.log)
		return
	}
	defer file.Close()

	c.Header("Content-Type", contentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))

	if _, err := io.Copy(c.Writer, file); err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetVideo: failed to copy video")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_copy", err, n.log)
		return
	}
	n.log.Infof("Endpoint Hit: Video served: %s", filename)
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

// @Summary Gets a video from the media manager
// @Description Serves a video file from the media manager
// @Tags media
// @Security Bearer
// @Accept json
// @Produce json
// @Param file path string true "video file"
// @Success 200 {file} byte "video file"
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/media/render/track/{file} [get]
func (n *Node) apiMediaGetAudio(c *gin.Context) {
	filename := c.Param("file")
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaGetAudio: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetAudio: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	fileType, filepath, err := n.media.GetAudio(filename)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetAudio: failed to get audio")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_audio", err, n.log)
		return
	}

	c.Header("Content-Type", fileType.MIME.Value)

	c.File(filepath)
	n.log.Infof("Endpoint Hit: Audio served: %s", filename)
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

	hash, err := n.media.AddAudio(openedFile)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaUploadAudio: failed to add audio")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_add_audio", err, n.log)
		return
	}

	response := dto.HashResponse{
		Hash: hash,
	}

	c.JSON(http.StatusOK, response)
}

// @Summary Deletes an audio file from the media manager
// @Description Deletes an audio file based on the provided filename from the media manager
// @Tags media
// @Security Bearer
// @Param file path string true "audio filename"
// @Success 200
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/media/deltrack/{file} [delete]
func (n *Node) apiMediaDeleteAudio(c *gin.Context) {
	filename := c.Param("file")
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaDeleteAudio: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaDeleteAudio: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	if err := n.media.DeleteAudio(filename); err != nil {
		err = errors.WithMessage(err, "Node: apiMediaDeleteAudio: failed to delete audio")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_delete", err, n.log)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary Gets a texture from the (internal) media-manager
// @Description Serves a generic texture from the (internal) media-manager
// @Tags media
// @Security Bearer
// @Accept application/octet-stream
// @Produce application/octet-stream
// @Param rsize path string true "Rendering size parameter. Should be in the format 's' followed by a single digit."
// @Param file path string true "Texture file identifier"
// @Failure 400 {object} api.HTTPError
// @Failure 500 {object} api.HTTPError
// @Router /api/v4/media/render/texture/{rsize}/{file} [get]
func (n *Node) apiMediaGetTexture(c *gin.Context) {
	tm1 := makeTimestamp()

	rsize := c.Param("rsize")
	filename := c.Param("file")
	match, err := regexp.MatchString(`^s[0-9]$`, rsize)
	if !match {
		err := errors.New("Node: apiMediaGetTexture: invalid rsize format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetTexture: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	match, err = regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaGetTexture: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetTexture: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	meta, filepath, err := n.media.GetTexture(rsize, filename)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetTexture: failed to get texture")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_texture", err, n.log)
		return
	}

	c.Header("Content-Type", meta.Mime)
	c.Header("x-height", strconv.Itoa(meta.H))
	c.Header("x-width", strconv.Itoa(meta.W))
	c.File(*filepath)
	n.log.Infof("Endpoint Hit: Texture served: %s %d", filename, makeTimestamp()-tm1)
}

// @Summary Gets an image from the (internal) media-manager
// @Description Serves a generic image from the (internal) media-manager
// @Tags media
// @Security Bearer
// @Accept json
// @Produce json
// @Param file path string true "image file"
// @Success 200 {object} dto.HashResponse
// @Failure 400 {object} api.HTTPError
// @Router /api/v4/media/render/asset/{file} [get]
func (n *Node) apiMediaGetAsset(c *gin.Context) {
	filename := c.Param("file")
	match, err := regexp.MatchString(`^[a-zA-Z0-9]+$`, filename)
	if !match {
		err := errors.New("Node: apiMediaGetAsset: invalid filename format")
		api.AbortRequest(c, http.StatusBadRequest, "invalid_format", err, n.log)
		return
	}
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetAsset: failed to match regexp string")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_validate", err, n.log)
		return
	}

	fileType, filepath, err := n.media.GetAsset(filename)
	if err != nil {
		err := errors.WithMessage(err, "Node: apiMediaGetAsset: failed to get asset")
		api.AbortRequest(c, http.StatusInternalServerError, "failed_to_get_asset", err, n.log)
		return
	}

	c.Header("Content-Type", fileType.MIME.Value)

	c.File(filepath)
	n.log.Infof("Endpoint Hit: Asset served: %s", filename)
}
