package node

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"

	"github.com/momentum-xyz/ubercontroller/universe/api"
)

func (n *Node) apiUsersCheck(c *gin.Context) {
	inBody := struct {
		IDToken string `json:"idToken"`
	}{}

	if err := c.ShouldBindJSON(&inBody); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errors.WithMessage(err, "failed to bind json"),
		})
		return
	}

	token, err := api.VerifyToken(c, api.GetTokenFromRequest(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"message": errors.WithMessage(err, "failed to verify token"),
		})
		return
	}

	idToken, err := api.ParseToken(inBody.IDToken)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": errors.WithMessage(err, "failed to parse id token"),
		})
		return
	}

	if token.GetSubject() != idToken.Subject {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("tokens subject mismatch: %s != %s", token.GetSubject(), idToken.Subject),
		})
		return
	}

	//userID, err := uuid.Parse(idToken.Subject)
	//if err != nil {
	//	c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
	//		"message": errors.WithMessage(err, "failed to parse user id"),
	//	})
	//	return
	//}

	//userEntry, err := n.db.UsersGetUserByID(n.ctx, userID)
	//if err != nil {
	//	userEntry = &entry.User{
	//		UserID: &userID,
	//		Profile: &entry.UserProfile{
	//			Name:  &idToken.Name,
	//			Email: &idToken.Email,
	//		},
	//	}
	//}

	if idToken.Guest.IsGuest {

	}
}
