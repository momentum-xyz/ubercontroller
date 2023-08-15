package processor

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
	"github.com/pkg/errors"
)

var AllowedVideoTypes = map[types.Type]bool{
	matchers.TypeMp4:  true,
	matchers.TypeMov:  true,
	matchers.TypeWmv:  true,
	matchers.TypeMpeg: true,
	matchers.TypeWebm: true,
	matchers.TypeMkv:  true,
}

func (p *Processor) ProcessVideo(body io.ReadCloser) (string, error) {
	hasher := md5.New()

	buf := make([]byte, 265)
	n, err := body.Read(buf)
	if err != nil {
		return "", errors.Wrap(err, "Error reading from body")
	}

	t, err := filetype.Get(buf[:n])
	if err != nil {
		return "", errors.Wrap(err, "Error getting filetype")
	}

	if _, ok := AllowedVideoTypes[t]; !ok {
		return "", errors.New("Video type not supported")
	}

	file, err := os.CreateTemp(p.Videopath, "tmp")
	if err != nil {
		return "", errors.Wrap(err, "Error creating temp file")
	}

	defer func() {
		if err != nil {
			os.Remove(file.Name())
		}
	}()

	mw := io.MultiWriter(file, hasher)

	if _, err := mw.Write(buf[:n]); err != nil {
		return "", errors.Wrap(err, "Error writing to temp file")
	}

	if _, err = io.Copy(mw, body); err != nil {
		return "", errors.Wrap(err, "Error copying to temp file")
	}

	if err = file.Close(); err != nil {
		return "", errors.Wrap(err, "Error closing the file")
	}

	hash := hex.EncodeToString(hasher.Sum(nil))
	if err = os.Rename(file.Name(), filepath.Join(p.Videopath, hash)); err != nil {
		return "", errors.Wrap(err, "Error renaming the file")
	}

	return hash, nil
}
