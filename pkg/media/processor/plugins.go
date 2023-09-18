package processor

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func createMd5Hash(input string) string {
	hash := md5.New()
	hash.Write([]byte(input))
	hashInBytes := hash.Sum(nil)[:16]
	return hex.EncodeToString(hashInBytes)
}

func extract(gzipStream io.Reader, targetDir string) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)
	manifestExists := false

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		target := filepath.Join(targetDir, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if header.Name == "manifest.json" {
				manifestExists = true
			}
			file, err := os.Create(target)
			if err != nil {
				return err
			}
			if _, err := io.Copy(file, tarReader); err != nil {
				return err
			}
			file.Close()
		}
	}

	if !manifestExists {
		return errors.New("manifest.json not found")
	}

	return nil
}

func (p *Processor) ProcessPlugin(archivePath string) (string, error) {
	hash := createMd5Hash(filepath.Base(archivePath))

	targetDir := p.Pluginpath + "/" + hash
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		err := errors.WithMessage(err, "Error creating directory")
		return "", err
	}

	file, err := os.Open(archivePath)
	if err != nil {
		err := errors.WithMessage(err, "Error opening file")
		return "", err
	}
	defer file.Close()

	if err := extract(file, targetDir); err != nil {
		err := errors.WithMessage(err, "Error extracting archive")
		return "", err
	}

	return hash, nil
}
