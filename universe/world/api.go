package world

import "github.com/gin-gonic/gin"

func (w *World) RegisterAPI(r *gin.Engine) {
	w.registerWorldAPI(r)
	w.registerUsersAPI(r)
}
