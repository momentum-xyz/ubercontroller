package media

import (
	"fmt"
	"go.uber.org/zap"
	"io"
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/pkg/media/processor"
	"github.com/momentum-xyz/ubercontroller/types"
)

type Media struct {
	ctx    types.NodeContext
	cfg    *config.Config
	log    *zap.SugaredLogger
	router *gin.Engine

	processor *processor.Processor
}

func NewMedia() *Media {
	media := &Media{}

	return media
}

func (m *Media) Initialize(ctx types.NodeContext) error {
	m.ctx = ctx
	m.log = ctx.Logger()
	m.cfg = ctx.Config()

	p := processor.NewProcessor()
	p.Initialize(ctx)
	m.processor = p

	return nil
}

func (m *Media) AddImage(file multipart.File) (string, error) {
	fmt.Println("Endpoint Hit: AddImage")

	body, err := io.ReadAll(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error reading file: %v")
	}
	hash, err := m.processor.ProcessImage(body)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing image: %v")
	}
	return hash, err
}

func (m *Media) AddFrame(file multipart.File) (string, error) {
	m.log.Debug("Endpoint Hit: AddFrame")

	body, err := io.ReadAll(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error reading file: %v")
	}
	hash, err := m.processor.ProcessFrame(body)
	if err != nil {
		return "", errors.WithMessagef(err, "error processing frame: %v")
	}

	return hash, err
}

func (m *Media) AddTube(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddTube")
	body, err := io.ReadAll(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error reading file: %v")
	}

	hash, err := m.processor.ProcessTube(body)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing image: %v")
	}

	return hash, err
}

func (m *Media) AddVideo(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddVideo")

	hash, err := m.processor.ProcessVideo(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing video: %v")
	}

	return hash, err
}

func (m *Media) AddTrack(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddTrack")

	hash, err := m.processor.ProcessTrack(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing audio: %v")
	}

	return hash, err
}

func (m *Media) AddAsset(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddAsset")

	hash, err := m.processor.ProcessAsset(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing asset: %v")
	}

	return hash, err
}
