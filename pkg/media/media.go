package media

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/pkg/media/processor"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
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

func (m *Media) Load() error {
	universe.GetNode().AddAPIRegister(m)

	return nil
}
