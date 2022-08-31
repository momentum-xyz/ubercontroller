package types

import (
	"context"
	"github.com/gin-gonic/gin"

	"github.com/google/uuid"
)

type IDer interface {
	GetID() uuid.UUID
}

type Initializer interface {
	Initialize(ctx context.Context) error
}

type Runner interface {
	Run(ctx context.Context) error
}

type Stopper interface {
	Stop() error
}

type RunStopper interface {
	Runner
	Stopper
}

type Loader interface {
	Load(ctx context.Context) error
}

type Saver interface {
	Save(ctx context.Context) error
}

type LoadSaver interface {
	Loader
	Saver
}

type APIRegister interface {
	RegisterAPI(r *gin.Engine)
}
