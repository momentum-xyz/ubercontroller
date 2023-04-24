package seed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/universe/logic/api/dto"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func seedMedia(ctx context.Context) error {
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.New("failed to get config from context")
	}

	basePath := cfg.Settings.SeedDataFiles
	client := &http.Client{}

	return seedPath(ctx, client, basePath)
}

func seedPath(ctx context.Context, client *http.Client, basePath string) error {
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
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

			fHash, err := uploadSeedFile(ctx, client, file)
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

func uploadSeedFile(ctx context.Context, client *http.Client, f *os.File) (string, error) {
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	log := utils.GetFromAny(ctx.Value(types.LoggerContextKey), (*zap.SugaredLogger)(nil))
	if cfg == nil {
		return "", errors.New("failed to get config from context")
	}

	ext := filepath.Ext(f.Name())

	var uploadURL string
	var mimeType string
	switch ext {
	case ".png":
		uploadURL = cfg.Common.RenderInternalURL + "/render/addimage"
		mimeType = "image/png"
	case ".glb":
		uploadURL = cfg.Common.RenderInternalURL + "/addasset"
		mimeType = "model/gltf-binary"
	default:
		log.Warnf("Unhandled seed file type %s", ext)
		return "", nil
	}
	return uploadFile(ctx, client, f, uploadURL, mimeType)
}

func uploadFile(ctx context.Context, client *http.Client, f *os.File, renderURL string, mimeType string) (string, error) {
	req, err := http.NewRequest("POST", renderURL, f)
	if err != nil {
		return "", fmt.Errorf("media manager request: %w", err)
	}

	req.Header.Set("Content-Type", mimeType)
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to post seed data to media-manager: %w", err)
	}

	defer resp.Body.Close()
	response := dto.HashResponse{}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return "", fmt.Errorf("decode json into response: %w", err)
	}

	return response.Hash, nil
}
