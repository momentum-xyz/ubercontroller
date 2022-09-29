package world

import (
	"github.com/gin-gonic/gin"
)

func (w *World) registerUsersAPI(r *gin.Engine) {
	w.log.Infof("Registering api for users: %s...", w.GetID())
}
