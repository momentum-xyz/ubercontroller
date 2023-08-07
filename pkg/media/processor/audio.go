package processor

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"errors"
	"io"
	"os"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/matchers"
	"github.com/h2non/filetype/types"
)

var AllowedAudioTypes = map[types.Type]bool{
	matchers.TypeMp3:  true,
	matchers.TypeOgg:  true,
	matchers.TypeAac:  true,
	matchers.TypeWebm: true,
	matchers.TypeFlac: true,
}

func (p *Processor) ProcessTrack(body io.ReadCloser) (error, string) {
	hasher := md5.New()
	bodyReader := io.TeeReader(body, hasher)
	buf := make([]byte, 265)

	n, err := bodyReader.Read(buf)
	if p.HandleError(err) {
		return err, ""
	}

	t, err := filetype.Get(buf[:n])

	if _, ok := AllowedAudioTypes[t]; !ok {
		return errors.New("audio type not accepted"), ""
	}

	file, err := os.CreateTemp(p.Audiopath, "tmp")
	if err != nil {
		return err, ""
	}

	defer func() {
		_, err := os.Stat(file.Name())
		if err == nil {
			os.Remove(file.Name())
		}
	}()

	file.Write(buf[:n])
	_, err = io.Copy(file, bodyReader)

	if err != nil {
		return err, ""
	}
	file.Close()

	bhash := hasher.Sum(nil)
	hash := hex.EncodeToString(bhash[:])
	os.Rename(file.Name(), p.Audiopath+hash)

	return err, hash
}
