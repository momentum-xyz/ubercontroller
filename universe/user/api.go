package user

import (
	"github.com/gin-gonic/gin"
)

func (u *User) RegisterAPI(r *gin.Engine) {
	u.log.Infof("Registering api for user: %s...", u.GetID())
	//v1 := r.Group(u.cfg.Common.APIPrefix)
	//{
	//	v1.GET("/user/check", u.Check)
	//}
}

//func (u *User) Check(c *gin.Context) {
//	provider, err := oidc.NewProvider(u.ctx, u.cfg.Auth.OIDCWeb3URL)
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to connect to oidc server"})
//		return
//	}
//
//	verifier := provider.Verifier(&oidc.Config{ClientID: u.cfg.Auth.OIDCWeb3ID})
//	idToken, err := verifier.Verify(u.ctx, "")
//	if err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to verify idToken"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": idToken})
//}
