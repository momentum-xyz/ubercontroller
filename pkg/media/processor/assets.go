package processor

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/h2non/filetype"
	"github.com/h2non/filetype/types"
)

var glbType types.Type
var AllowedAssetsTypes map[types.Type]bool

func init() {
	glbType = filetype.NewType("glb", "model/gltf-binary")
	filetype.AddMatcher(glbType, glbMatcher)
	AllowedAssetsTypes = make(map[types.Type]bool)
	AllowedAssetsTypes[glbType] = true
}

func glbMatcher(buf []byte) bool {
	return len(buf) > 3 && buf[0] == 'g' && buf[1] == 'l' && buf[2] == 'T' && buf[3] == 'F'
}

func (p *Processor) ProcessAsset(body io.ReadCloser) (string, error) {
	hasher := md5.New()
	bodyReader := io.TeeReader(body, hasher)
	buf := make([]byte, 265)
	n, err := bodyReader.Read(buf)
	if err != nil {
		return "", err
	}

	t, err := filetype.Get(buf[:n])
	if err != nil {
		return "", err
	}

	if _, ok := AllowedAssetsTypes[t]; !ok {
		return "", fmt.Errorf("Not accepted asset type: %s", t.MIME.Value)
	}

	file, err := os.CreateTemp(p.Assetpath, "tmp")
	if err != nil {
		return "", err
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
		return "", err
	}

	file.Close()

	bhash := hasher.Sum(nil)
	hash := hex.EncodeToString(bhash[:])
	os.Rename(file.Name(), p.Assetpath+hash)

	return hash, err
}
