package media

import (
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"
	fileTypes "github.com/h2non/filetype/types"

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

var Tsizes = map[string]int{
	"s0": 1024,
	"s1": 4096,
	"s2": 9216,
	"s3": 25600,
	"s4": 65536,
	"s5": 193600,
	"s6": 577600,
	"s7": 1721344,
	"s8": 5062500,
	"s9": 14745600,
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

func (m *Media) GetImage(filename string) (*processor.MetaDef, *string, error) {
	m.log.Debug("Endpoint Hit: Image Get:", filename)

	meta, filepath := m.processor.Present(&(filename))
	if meta == nil {
		return nil, nil, errors.New("no meta for file")
	}

	return meta, filepath, nil
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

func (m *Media) GetTexture(rsize string, filename string) (*processor.MetaDef, *string, error) {
	if _, ok := Tsizes[rsize]; !ok {
		return nil, nil, errors.New("tsize not found for texture")
	}

	m.log.Debug("Endpoint Hit: Texture Get:", filename, rsize)

	meta, filepath := m.processor.PresentTexture(&(filename), rsize)
	if meta == nil {
		return nil, nil, errors.New("no meta for file")
	}

	return meta, filepath, nil
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

func (m *Media) GetVideo(filename string) (*os.File, os.FileInfo, string, error) {
	m.log.Debug("Endpoint Hit: Video Get:", filename)

	filepath := m.processor.Videopath + filename
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, "", errors.WithMessage(err, "error opening file")
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, nil, "", errors.WithMessage(err, "error getting file info")
	}

	buf := make([]byte, 512)

	_, err = file.Read(buf)
	if err != nil {
		return nil, nil, "", errors.WithMessage(err, "error reading file buffer")
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		return nil, nil, "", errors.WithMessage(err, "error seeking file")
	}

	contentType := http.DetectContentType(buf)
	return file, fileInfo, contentType, nil
}

func (m *Media) AddVideo(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddVideo")

	hash, err := m.processor.ProcessVideo(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing video: %v")
	}

	return hash, err
}

func (m *Media) GetAudio(filename string) (*fileTypes.Type, string, error) {
	m.log.Debug("Endpoint Hit: Audio Get:", filename)

	filepath := m.processor.Audiopath + filename
	buf := make([]byte, 264)
	file, err := os.Open(filepath)
	if err != nil {
		return nil, "", errors.WithMessage(err, "error opening audio")
	}
	defer file.Close()

	n, err := file.Read(buf)
	if err != nil {
		return nil, "", errors.WithMessage(err, "error reading buffer")
	}

	fileType, err := filetype.Get(buf[:n])
	if err != nil {
		return nil, "", errors.WithMessage(err, "error getting audio type")
	}

	return &fileType, filepath, nil
}

func (m *Media) AddAudio(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddAudio")

	hash, err := m.processor.ProcessTrack(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing audio: %v")
	}

	return hash, err
}

func (m *Media) GetAsset(c *gin.Context) {
	filename := c.Param("file")

	L().Debug("Endpoint Hit: Asset Get:", filename)

	// res, filepath := x.present(&(filename))
	filepath := x.Assetpath + filename

	buf := make([]byte, 264)
	f, err := os.Open(filepath)
	if check_error(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		sentry.CaptureException(err)
		L().Error(err)
		return
	}
	defer f.Close()

	n, err := f.Read(buf)
	if check_error(err) {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		sentry.CaptureException(err)
		L().Error(err)
		return
	}

	ftype, err := filetype.Get(buf[:n])

	c.Header("Content-Type", ftype.MIME.Value)

	c.File(filepath)
	L().Infof("Endpoint Hit: Asset served: %s", filename)
}

func (m *Media) AddAsset(file multipart.File) (string, error) {
	m.log.Info("Endpoint Hit: AddAsset")

	hash, err := m.processor.ProcessAsset(file)
	if err != nil {
		return "", errors.WithMessagef(err, "error writing asset: %v")
	}

	return hash, err
}
