package world

import "github.com/gin-gonic/gin"

func (w *World) RegisterAPI(r *gin.Engine) {
	w.log.Infof("Registering api for world: %s...", w.GetID())

	w.registerWorldAPI(r)
	w.registerUsersAPI(r)
}
