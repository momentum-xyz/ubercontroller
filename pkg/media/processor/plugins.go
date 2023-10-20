package processor

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func dirMD5(dirPath string) (string, error) {
	fileSystem := os.DirFS(dirPath)
	hasher := md5.New()
	err := fs.WalkDir(
		fileSystem, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.Type().IsRegular() {
				fmt.Println(filepath.Join(dirPath, path))
				f, err := os.Open(filepath.Join(dirPath, path))
				if err != nil {
					return err
				}
				defer f.Close()
				_, err = io.Copy(hasher, f)
				if err != nil {
					return err
				}
			}
			return nil
		},
	)
	if err != nil {
		return "", err
	}
	binHash := hasher.Sum(nil)
	hash := hex.EncodeToString(binHash[:])
	return hash, nil
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

	tempDir, err := os.MkdirTemp(p.Pluginpath, ".upload-")
	if err != nil {
		err := errors.WithMessage(err, "Error creating temporary directory")
		return "", err
	}
	defer os.RemoveAll(tempDir)

	file, err := os.Open(archivePath)
	if err != nil {
		err := errors.WithMessage(err, "Error opening file")
		return "", err
	}
	defer file.Close()

	if err := extract(file, tempDir); err != nil {
		err := errors.WithMessage(err, "Error extracting archive")
		return "", err
	}

	hash, err := dirMD5(tempDir)
	if err != nil {
		err := errors.WithMessage(err, "Error calculating md5 hash of content")
		return "", err
	}

	targetDir := p.Pluginpath + "/" + hash
	tStat, err := os.Stat(targetDir)
	if err == nil {
		if tStat.IsDir() {
			return hash, nil
		} else {
			return "", errors.New("There is already regular file with such name")
		}
	}

	err = os.Rename(tempDir, targetDir)

	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			tStat, err = os.Stat(targetDir)
			if err == nil {
				if tStat.IsDir() {
					return hash, nil
				}
			}
		}

		err := errors.WithMessage(err, "Error renaming temp dir to target name")
		return "", err
	}

	os.Chmod(targetDir, 0755)

	return hash, nil
}
