package processor

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"image"
	"os"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/getsentry/sentry-go"
	lru "github.com/hashicorp/golang-lru"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
)

type Processor struct {
	ctx types.NodeContext
	log *zap.SugaredLogger
	cfg *config.Config

	Fontpath   string
	Imagepath  string
	Audiopath  string
	Videopath  string
	Assetpath  string
	Pluginpath string

	ImPathF string
	ImPathS map[string]string

	ImageMapF *lru.Cache
	ImageMapS map[string]*lru.Cache

	PresentMutex sync.Mutex

	framesinprogress map[string]bool

	RenderQueue chan *types.FrameRenderRequest
	RenderDone  chan *types.FrameRenderRequest

	ImagesRescaleInProgress sync.Map
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
	p.Pluginpath = strings.TrimSuffix(p.cfg.Media.Pluginpath, "/") + "/"
	p.Fontpath = strings.TrimSuffix(p.cfg.Media.Fontpath, "/") + "/"
	p.framesinprogress = make(map[string]bool)
	p.RenderQueue = make(chan *types.FrameRenderRequest, 512)
	p.RenderDone = make(chan *types.FrameRenderRequest, 512)
	os.MkdirAll(p.Imagepath, os.ModePerm)
	os.MkdirAll(p.Videopath, os.ModePerm)
	os.MkdirAll(p.Audiopath, os.ModePerm)
	os.MkdirAll(p.Assetpath, os.ModePerm)
	os.MkdirAll(p.Pluginpath, os.ModePerm)

	os.MkdirAll(p.Imagepath+"F", os.ModePerm)
	p.ImPathF = p.Imagepath + "F/"

	p.ImageMapF, _ = lru.New(defaultCacheSize)
	p.ImageMapS = make(map[string]*lru.Cache)
	p.ImPathS = make(map[string]string)
	for rs := range types.Tsizes {
		os.MkdirAll(p.Imagepath+rs, os.ModePerm)
		p.ImPathS[rs] = p.Imagepath + rs + "/"
		p.ImageMapS[rs], _ = lru.New(defaultCacheSize)
	}

	os.MkdirAll(strings.TrimSuffix(p.ImPathF, "/"), 0775)
	go p.run()
	return p
}

func (p *Processor) checkError(err error) bool {
	if err != nil {
		p.log.Error(err)
		return true
	}
	return false
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

func (p *Processor) ProcessFrame(body []byte) (string, error) {
	var payload types.FrameDesc
	if err := json.Unmarshal(body, &payload); err != nil {
		return "", errors.WithMessage(err, "error during unmarshal")
	}

	hash := p.GetMD5HashByte(body)
	//if IDi != "" && IDi != ID {
	//	return "", errors.WithMessage(err, "hash mismatch")
	//}

	// what does this check do?

	p.log.Debug(hash)
	str, err := json.MarshalIndent(payload, "", "    ")
	if err != nil {
		sentry.CaptureException(err)
	}

	p.log.Debug(string(str))
	res, _ := p.Present(hash)
	if res == nil {
		req := &types.FrameRenderRequest{ID: &hash, Frame: &payload}
		req.Wg.Add(1)
		p.RenderQueue <- req
		req.Wg.Wait()
	}

	return hash, nil
}

func (p *Processor) Present(imageID string) (*types.MetaDef, *string) {
	filePath := path.Join(p.ImPathF, imageID)
	res, ok := p.ImageMapF.Get(imageID)
	if ok {
		p.log.Debug(imageID + " is already in the map")
		return res.(*types.MetaDef), &filePath
	}

	reader, err := os.Open(filePath)
	p.log.Debug(filePath)
	if err != nil {
		return nil, nil
	}

	defer reader.Close()
	im, format, err1 := image.DecodeConfig(reader)
	meta := new(types.MetaDef)
	if err1 != nil {
		p.log.Debugf("%s: %v\n", imageID, err1)
		meta.Mime = "image/png"
	} else {
		meta.H = im.Height
		meta.W = im.Width
		meta.Mime = "image/" + format
		p.log.Debugf("%s %d %d\n", imageID, im.Width, im.Height)
	}
	p.log.Debug("Mime:", meta.Mime)
	p.ImageMapF.Add(imageID, meta)

	return meta, &filePath
}

func (p *Processor) PresentTexture(ID *string, rsize string) (*types.MetaDef, *string) {
	fpath := p.ImPathS[rsize] + *ID
	meta0, ok := p.ImageMapS[rsize].Get(*ID)
	if ok {
		p.log.Debug(*ID + " is already in the map")
		return meta0.(*types.MetaDef), &fpath
	}

	defer func() {
		p.ImagesRescaleInProgress.Delete(fpath)
	}()

	if _, ok := p.ImagesRescaleInProgress.Load(fpath); ok {
		maxRetries := 20
		retryInterval := time.Millisecond * 300
		for i := 1; i <= maxRetries; i++ {
			time.Sleep(retryInterval)
			meta0, ok := p.ImageMapS[rsize].Get(*ID)
			if ok {
				p.log.Debug(*ID + " found in map on retry round N: " + strconv.Itoa(i))
				return meta0.(*types.MetaDef), &fpath
			}
		}
		p.log.Error(*ID + " timeout reached while waiting image to rescale:" + strconv.Itoa(int(retryInterval.Milliseconds())*maxRetries) + " ms")
		return nil, nil
	}

	p.log.Debug(fpath)
	reader, err := os.Open(fpath)
	if err != nil {
		p.ImagesRescaleInProgress.Store(fpath, true)
		converted := false
		p.log.Debug(*ID + " : converting from full")
		if meta, filepath := p.Present(*ID); meta != nil {
			if reader, err = os.Open(*filepath); !p.checkError(err) {
				if img, _, errl := image.Decode(reader); !p.checkError(errl) {
					if err = p.WriteToScaled(*ID, img, rsize); !p.checkError(err) {
						if reader, err = os.Open(fpath); err == nil {
							converted = true
						}
					}
				}
			}
		}

		if !converted {
			return nil, nil
		}
	}

	defer reader.Close()
	im, format, err1 := image.DecodeConfig(reader)
	meta := new(types.MetaDef)
	if err1 != nil {
		p.log.Debugf("%s: %v\n", *ID, err1)
		meta.Mime = "image/png"
	} else {
		meta.H = im.Height
		meta.W = im.Width
		meta.Mime = "image/" + format
		p.log.Debugf("%s %d %d\n", *ID, im.Width, im.Height)
	}
	p.log.Debug("Mime:", meta.Mime)
	p.ImageMapS[rsize].Add(*ID, meta)
	return meta, &fpath
}
