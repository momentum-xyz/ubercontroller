package seed

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/momentum-xyz/ubercontroller/universe"
)

func SeedMedia(node universe.Node) error {
	cfg := node.GetConfig()
	basePath := cfg.Settings.SeedDataFiles

	return seedPath(node, basePath)
}

func seedPath(node universe.Node, basePath string) error {
	log := node.GetLogger()
	return filepath.WalkDir(basePath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("seed data files: %w", err)
		}
		if !entry.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return fmt.Errorf("reading seed data file: %w", err)
			}

			defer file.Close()

			fHash, err := uploadSeedFile(node, file)
			if err != nil {
				return fmt.Errorf("upload seed data file: %w", err)
			}
			rPath, err := filepath.Rel(basePath, path)
			if err != nil {
				rPath = path
			}
			log.Infof("Seed %s = %s", rPath, fHash)
		}
		return nil
	})
}

func uploadSeedFile(node universe.Node, file *os.File) (string, error) {
	log := node.GetLogger()
	media := node.GetMedia()
	ext := filepath.Ext(file.Name())

	var hash string
	var err error

	switch ext {
	case ".png":
		hash, err = media.AddImage(file)
		if err != nil {
			return "", fmt.Errorf("error seeding image: %w", err)
		}
	case ".glb":
		hash, err = media.AddAsset(file)
		if err != nil {
			return "", fmt.Errorf("error seeding asset: %w", err)
		}
	default:
		log.Warnf("Unhandled seed file type %s", ext)
		return "", nil
	}

	return hash, nil
}
