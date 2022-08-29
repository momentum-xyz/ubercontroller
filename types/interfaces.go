package types

import (
	"context"

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
