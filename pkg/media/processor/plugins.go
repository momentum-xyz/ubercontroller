package processor

import (
	"archive/tar"
	"compress/gzip"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"

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

func extract(gzipStream io.Reader, targetDir string, stripFirstLevel bool) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if stripFirstLevel {
			// Strip the first level of directory structure
			parts := strings.SplitN(filepath.Clean(header.Name), "/", 2)
			if len(parts) > 1 {
				header.Name = parts[1]
			} else {
				continue // skip the top-level directory itself
			}
		}

		target := filepath.Join(targetDir, header.Name)
		fmt.Println("Extracting", target)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
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

	return nil
}

func storeToFile(body io.ReadCloser) (string, error) {
	buf := make([]byte, 265)

	n, err := body.Read(buf)
	if err != nil {
		return "", err
	}

	file, err := os.CreateTemp("", "tmp")
	if err != nil {
		return "", err
	}
	fmt.Println("Created tmp file", file.Name())

	defer func() {
		// _, err := os.Stat(file.Name())
		// if err == nil {
		if err != nil {
			fmt.Println("Removing tmp file after error", file.Name())
			os.Remove(file.Name())
		}
	}()

	file.Write(buf[:n])
	_, err = io.Copy(file, body)
	if err != nil {
		return "", err
	}

	file.Close()

	return file.Name(), nil
}

type AttributeTypeDescription struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Sync        *string `json:"sync"`
}

type Manifest struct {
	Name           string                      `json:"name"`
	Description    string                      `json:"description"`
	Version        string                      `json:"version"`
	AttributeTypes *[]AttributeTypeDescription `json:"attribute_types"`
	Scopes         map[string][]string         `json:"scopes"`
}

func (p *Processor) LoadPluginManifest(pluginHash string) (*Manifest, error) {
	manifestFn := filepath.Join(p.Pluginpath, pluginHash, "manifest.json")

	_, err := os.Stat(manifestFn)
	if err != nil {
		err := errors.WithMessage(err, "manifest.json not found")
		return nil, err
	}

	file, err := os.Open(manifestFn)
	if err != nil {
		err := errors.WithMessage(err, "Error opening file")
		return nil, err
	}

	var manifest Manifest
	err = json.NewDecoder(file).Decode(&manifest)
	if err != nil {
		err := errors.WithMessage(err, "Error decoding manifest")
		return nil, err
	}

	return &manifest, nil
}

func (p *Processor) ProcessPlugin(body io.ReadCloser) (string, error) {
	archivePath, err := storeToFile(body)
	if err != nil {
		return "", err
	}
	// defer os.Remove(archivePath)
	defer func() {
		fmt.Println("Removing tmp file", archivePath)
		os.Remove(archivePath)
	}()

	tempDir, err := os.MkdirTemp(p.Pluginpath, ".upload-")
	if err != nil {
		err := errors.WithMessage(err, "Error creating temporary directory")
		return "", err
	}
	fmt.Println("Created tmp dir", tempDir)
	// defer os.RemoveAll(tempDir)
	defer func() {
		fmt.Println("Removing tmp dir", tempDir)
		os.RemoveAll(tempDir)
	}()

	file, err := os.Open(archivePath)
	if err != nil {
		err := errors.WithMessage(err, "Error opening file")
		return "", err
	}
	defer file.Close()

	if err := extract(file, tempDir, true); err != nil {
		err := errors.WithMessage(err, "Error extracting archive")
		return "", err
	}

	if _, err := p.LoadPluginManifest(path.Base(tempDir)); err != nil {
		err := errors.WithMessage(err, "Error loading manifest")
		return "", err
	}

	hash, err := dirMD5(tempDir)
	if err != nil {
		err := errors.WithMessage(err, "Error calculating md5 hash of content")
		return "", err
	}

	targetDir := filepath.Join(p.Pluginpath, hash)
	tStat, err := os.Stat(targetDir)
	if err == nil {
		if tStat.IsDir() {
			return hash, nil
		} else {
			return "", errors.New("There is already regular file with such name")
		}
	}

	fmt.Println("Renaming", tempDir, "to", targetDir)
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
