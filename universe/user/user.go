package user

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/config"
	"github.com/momentum-xyz/ubercontroller/database"
	"go.uber.org/zap"
	"sync"
)

type User struct {
	id     uuid.UUID
	cfg    *config.Config
	ctx    context.Context
	log    *zap.SugaredLogger
	db     database.DB
	router *gin.Engine
	mu     sync.RWMutex
}

func (u *User) GetID() uuid.UUID {
	return u.id
}
