package media

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/pkg/media/processor"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils/umid"
)

type Media struct {
	node universe.Node

	ctx    types.NodeContext
	cfg    *config.Config
	log    *zap.SugaredLogger
	db     database.DB
	router *gin.Engine

	p *processor.Processor
}

func NewMedia(
	id umid.UMID,
	db database.DB,
) *Media {
	media := &Media{
		db: db,
	}

	return media
}

func (m *Media) Initialize(ctx types.NodeContext) error {
	m.ctx = ctx
	m.log = ctx.Logger()
	m.cfg = ctx.Config()

	m.node = universe.GetNode()
	return nil
}

func (m *Media) Load() error {
	universe.GetNode().AddAPIRegister(m)

	return nil
}
