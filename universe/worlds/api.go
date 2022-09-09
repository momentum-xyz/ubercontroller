package worlds

import "github.com/gin-gonic/gin"

func (w *Worlds) RegisterAPI(r *gin.Engine) {
	w.log.Info("Registering api for worlds...")
}
