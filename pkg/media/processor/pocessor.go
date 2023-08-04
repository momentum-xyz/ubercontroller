package processor

import (
	"crypto/md5"
	"encoding/hex"
	"go.uber.org/zap"
	"image"
	"os"
	"sync"

	"github.com/getsentry/sentry-go"
	lru "github.com/hashicorp/golang-lru"
)

type Processor struct {
	log *zap.SugaredLogger

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
