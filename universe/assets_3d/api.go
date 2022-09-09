package assets_3d

import "github.com/gin-gonic/gin"

func (a *Assets3d) RegisterAPI(r *gin.Engine) {
	a.log.Info("Registering api for assets 3d...")
}
