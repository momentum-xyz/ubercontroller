package node

import (
	"context"

	"github.com/momentum-xyz/ubercontroller/config"
)

type Auth struct {
	ctx context.Context
	cfg *config.Config
}
