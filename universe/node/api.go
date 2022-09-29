package node

import (
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CheckRequestBody struct {
	IdToken string `json:"idToken"`
}

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())
	v1 := r.Group(n.cfg.Common.APIPrefix)
	{
		v1.POST("/user/check", n.Check)
	}
}

func (n *Node) Check(c *gin.Context) {
	provider, err := oidc.NewProvider(n.ctx, n.cfg.Auth.OIDCWeb3URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to connect to oidc server"})
		return
	}

	var requestBody CheckRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no idToken received"})
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: n.cfg.Auth.OIDCWeb3ID})
	idToken, err := verifier.Verify(n.ctx, requestBody.IdToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to verify idToken"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": idToken})
}
