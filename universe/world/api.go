package world

import "github.com/gin-gonic/gin"

func (w *World) RegisterAPI(r *gin.Engine) {
	w.Space.RegisterAPI(r)
}
