package harvester

import (
	"context"
	"sync"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"

	"github.com/momentum-xyz/ubercontroller/config"
)

var instance *Harvester
var mu sync.Mutex
var logger *zap.SugaredLogger

func Initialise(ctx context.Context, log *zap.SugaredLogger, cfg *config.Config, pool *pgxpool.Pool) {
	mu.Lock()
	defer mu.Unlock()

	logger = log

	if instance == nil {
		instance = NewHarvester(pool)
	}
}

func GetInstance() *Harvester {
	if instance == nil {
		logger.Error("Harvester must be initialised")
	}

	return instance
}
