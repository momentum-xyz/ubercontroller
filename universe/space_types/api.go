package space_types

import "github.com/gin-gonic/gin"

func (s *SpaceTypes) RegisterAPI(r *gin.Engine) {
	s.log.Info("Registering api for space types...")
}
