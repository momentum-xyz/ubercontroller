package media

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/pkg/media/processor"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
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
	fmt.Println("Endpoint Hit: addImage")

	body, err := io.ReadAll(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error during reading body: %v")
	}
	hash, err := m.processor.ProcessImage(body)
	if err != nil {
		return "", errors.WithMessagef(err, "error during writing image: %v")
	}
	return hash, err
}
