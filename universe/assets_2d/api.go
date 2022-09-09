package assets_2d

import "github.com/gin-gonic/gin"

func (a *Assets2d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 2d...")
}
