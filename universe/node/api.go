package node

import (
	"encoding/json"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	u "github.com/momentum-xyz/ubercontroller/utils"
	"net/http"
	"strings"
)

type CheckRequestBody struct {
	IdToken string `json:"idToken"`
}

type audience []string

type idToken struct {
	Guest struct {
		IsGuest bool `json:"1"`
	}
	Web3Address string   `json:"web3_address"`
	Web3Type    string   `json:"web3_type"`
	Issuer      string   `json:"iss"`
	Subject     string   `json:"sub"`
	Audience    audience `json:"aud"`
	Expiry      int      `json:"exp"`
	IssuedAt    int      `json:"iat"`
}

func (n *Node) RegisterAPI(r *gin.Engine) {
	n.log.Infof("Registering api for node: %s...", n.GetID())
	v1 := r.Group(n.cfg.Common.APIPrefix)
	{
		v1.POST("/user/check", n.Check)
	}
}

func (n *Node) Check(c *gin.Context) {
	provider, err := oidc.NewProvider(n.ctx, n.cfg.Auth.OIDCURL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to connect to oidc server"})
		return
	}

	var requestBody CheckRequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no idToken received"})
		return
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: n.cfg.Auth.OIDCWeb3ID})

	parsedAccessToken := c.Request.Header["Authorization"][0]
	splitToken := strings.Split(parsedAccessToken, "Bearer ")
	parsedAccessToken = splitToken[1]

	accessToken, err := verifier.Verify(n.ctx, parsedAccessToken)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "failed to verify accessToken"})
		return
	}

	if accessToken != nil {
		// Check existence of user
		var idT idToken
		jwt, err := u.ParseJWT(requestBody.IdToken)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "malformed jwt"})
			return
		}
		if err := json.Unmarshal(jwt, &idT); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "failed to unmarshal claims"})
			return
		}
		if accessToken.Subject != idT.Subject {
			c.JSON(http.StatusBadRequest, gin.H{"message": "accessToken and idToken do not belong to same user"})
			return
		}

		if !idT.Guest.IsGuest {
			// Set usertype to user?
			// Create user
		} else {
			// Set usertype to temporary user?
			// Create user
		}

		// Check for invitation
		// Assign rights based on invitation?

		c.JSON(http.StatusOK, gin.H{"message": accessToken})
	}
}
