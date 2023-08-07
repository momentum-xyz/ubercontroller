package processor

import (
	"crypto/md5"
	"encoding/hex"
	"go.uber.org/zap"
	"image"
	"os"
	"strings"
	"sync"

	"github.com/getsentry/sentry-go"
	lru "github.com/hashicorp/golang-lru"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
)

type Processor struct {
	ctx types.NodeContext
	log *zap.SugaredLogger
	cfg *config.Config

	Fontpath  string
	Imagepath string
	Audiopath string
	Videopath string
	Assetpath string

	ImPathF string
	ImPathS map[string]string

	ImageMapF *lru.Cache
	ImageMapS map[string]*lru.Cache

	PresentMutex sync.Mutex

	framesinprogress map[string]bool

	RenderQueue chan *FrameRenderRequest
	RenderDone  chan *FrameRenderRequest

	ImagesRescaleInProgress sync.Map
}

type MetaDef struct {
	H, W int
	mime string
}

const defaultCacheSize = 1024

func NewProcessor() *Processor {
	processor := &Processor{}

	return processor
}

func (p *Processor) Initialize(ctx types.NodeContext) *Processor {
	p.ctx = ctx
	p.log = ctx.Logger()
	p.cfg = ctx.Config()
	p.Imagepath = strings.TrimSuffix(p.cfg.Media.Imagepath, "/") + "/"
	p.Videopath = strings.TrimSuffix(p.cfg.Media.Videopath, "/") + "/"
	p.Audiopath = strings.TrimSuffix(p.cfg.Media.Audiopath, "/") + "/"
	p.Assetpath = strings.TrimSuffix(p.cfg.Media.Assetpath, "/") + "/"
	p.Fontpath = strings.TrimSuffix(p.cfg.Media.Fontpath, "/") + "/"
	p.framesinprogress = make(map[string]bool)
	p.RenderQueue = make(chan *FrameRenderRequest, 512)
	p.RenderDone = make(chan *FrameRenderRequest, 512)
	os.MkdirAll(p.Imagepath, os.ModePerm)
	os.MkdirAll(p.Videopath, os.ModePerm)
	os.MkdirAll(p.Audiopath, os.ModePerm)
	os.MkdirAll(p.Assetpath, os.ModePerm)

	os.MkdirAll(p.Imagepath+"F", os.ModePerm)
	p.ImPathF = p.Imagepath + "F/"

	p.ImageMapF, _ = lru.New(defaultCacheSize)
	p.ImageMapS = make(map[string]*lru.Cache)
	p.ImPathS = make(map[string]string)
	for rs := range Tsizes {
		os.MkdirAll(p.Imagepath+rs, os.ModePerm)
		p.ImPathS[rs] = p.Imagepath + rs + "/"
		p.ImageMapS[rs], _ = lru.New(defaultCacheSize)
	}

	os.MkdirAll(strings.TrimSuffix(p.ImPathF, "/"), 0775)
	go p.run()
	return p
}

func (p *Processor) run() {
	p.log.Info("Processor Runner...")
	for {
		select {
		case req := <-p.RenderQueue:
			if !p.framesinprogress[*req.ID] {
				p.framesinprogress[*req.ID] = true
				go p.RenderFrame(req)
			}
			n := len(p.RenderQueue)
			for i := 0; i < n; i++ {
				req := <-p.RenderQueue
				if !p.framesinprogress[*req.ID] {
					p.framesinprogress[*req.ID] = true
					go p.RenderFrame(req)
				}
			}
		case req := <-p.RenderDone:
			delete(p.framesinprogress, *req.ID)
			n := len(p.RenderDone)
			for i := 0; i < n; i++ {
				req := <-p.RenderDone
				delete(p.framesinprogress, *req.ID)
			}
		}
	}
}

func (p *Processor) HandleError(err error) bool {
	if err != nil {
		p.log.Error(err)
		sentry.CaptureException(err)
		return true
	}
	return false
}

func (p *Processor) FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func (p *Processor) GetMD5HashByte(text []byte) string {
	hash := md5.Sum(text)
	return hex.EncodeToString(hash[:])
}

func (p *Processor) present(ID *string) (*MetaDef, *string) {
	fpath := p.ImPathF + *ID

	res, ok := p.ImageMapF.Get(*ID)
	if ok {
		p.log.Debug(*ID + " is already in the map")
		return res.(*MetaDef), &fpath
	}

	reader, err := os.Open(fpath)
	p.log.Debug(fpath)
	if err != nil {
		return nil, nil
	}

	defer reader.Close()
	im, format, err1 := image.DecodeConfig(reader)
	meta := new(MetaDef)
	if err1 != nil {
		p.log.Debugf("%s: %v\n", *ID, err1)
		meta.mime = "image/png"
	} else {
		meta.H = im.Height
		meta.W = im.Width
		meta.mime = "image/" + format
		p.log.Debugf("%s %d %d\n", *ID, im.Width, im.Height)
	}
	p.log.Debug("Mime:", meta.mime)
	p.ImageMapF.Add(*ID, meta)

	return meta, &fpath
}
