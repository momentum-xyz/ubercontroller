package api

import (
	"context"

	"github.com/zitadel/oidc/pkg/client/rs"

	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/types/generic"
)

var api = struct {
	ctx           context.Context
	cfg           *config.Config
	oidcProviders *generic.SyncMap[string, rs.ResourceServer]
}{}

func Initialize(ctx context.Context, cfg *config.Config) error {
	api.ctx = ctx
	api.cfg = cfg
	api.oidcProviders = generic.NewSyncMap[string, rs.ResourceServer]()

	return nil
}
