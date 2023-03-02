package seed

import (
	"context"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types"
	"github.com/momentum-xyz/ubercontroller/utils"
)

func seedMedia(ctx context.Context) error {
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.New("failed to get config from context")
	}

	basePath := cfg.Settings.SeedDataFiles
	client := &http.Client{}

	return filepath.WalkDir(basePath, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return errors.WithMessagef(err, "seed data files %s", path)
		}
		if !entry.IsDir() {
			file, err := os.Open(path)
			if err != nil {
				return errors.WithMessagef(err, "reading seed data file %s", path)
			}

			defer file.Close()

			if err := uploadSeedFile(ctx, client, file); err != nil {
				return errors.WithMessagef(err, "upload seed data file %s", path)
			}
		}
		return nil
	})
}

func uploadSeedFile(ctx context.Context, client *http.Client, f *os.File) error {
	cfg := utils.GetFromAny(ctx.Value(types.ConfigContextKey), (*config.Config)(nil))
	if cfg == nil {
		return errors.New("failed to get config from context")
	}

	req, err := http.NewRequest("POST", cfg.Common.RenderInternalURL+"/render/addimage", f)
	if err != nil {
		return errors.WithMessage(err, "failed to media manager url")
	}

	req.Header.Set("Content-Type", "image/png")
	resp, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, "failed to post seed data to media-manager")
	}

	defer resp.Body.Close()

	return nil
}
