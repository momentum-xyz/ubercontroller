package api

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/momentum-xyz/ubercontroller/utils"
	"github.com/pkg/errors"
)

func GetTokenFromRequest(c *gin.Context) string {
	authHeader := c.GetHeader("Authorization")
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func GetTokenFromContext(c *gin.Context) (jwt.Token, error) {
	value, ok := c.Get(TokenContextKey)
	if !ok {
		return jwt.Token{}, errors.Errorf("failed to get token value from context")
	}
	return utils.GetFromAny(value, jwt.Token{}), nil
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
	token, err := GetTokenFromContext(c)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get token from context")
	}
	userID, err := GetUserIDFromToken(token)
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to get user id from token")
	}
	return userID, nil
}

func GetUserIDFromToken(token jwt.Token) (uuid.UUID, error) {
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return uuid.Nil, errors.New("failed to get token claims")
	}
	userID, err := uuid.Parse(utils.GetFromAnyMap(claims, "sub", "")) // TODO! proper jwt parsing
	if err != nil {
		return uuid.Nil, errors.WithMessage(err, "failed to parse user id")
	}
	return userID, nil
}

func GenerateChallenge(wallet string) (string, error) {
	return fmt.Sprintf(
		"Please sign this message with the private key for address %s to prove that you own it. %s",
		wallet, uuid.New().String(),
	), nil
}

func VerifyPolkadotSignature(wallet, challenge, signature string) (bool, error) {
	pub, err := schnorrkel.NewPublicKeyFromHex(wallet)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get public key")
	}
	sig, err := schnorrkel.NewSignatureFromHex(signature)
	if err != nil {
		return false, errors.WithMessage(err, "failed to get signature")
	}

	trContext := []byte("substrate")
	trMessage := bytes.NewBufferString("<Bytes>")
	trMessage.WriteString(challenge)
	trMessage.WriteString("</Bytes>")

	transcript := schnorrkel.NewSigningContext(trContext, trMessage.Bytes())

	return pub.Verify(sig, transcript)
}

// CreateJWTToken saves a jwt token with the given userID as subject
// and signed with the given secret
func CreateJWTToken(userID uuid.UUID, secret []byte) (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(4 * time.Hour).Unix(),
		Issuer:    "ubercontroller",
		Subject:   userID.String(),
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedString, err := jwt.SignedString(secret)
	if err != nil {
		return "", errors.WithMessage(err, "failed to sign token")
	}

	return signedString, nil
}

func ValidateJWT(signedString string, secret []byte) (*jwt.Token, error) {
	parser := new(jwt.Parser)
	token, _, err := parser.ParseUnverified(signedString, jwt.MapClaims{})
	return token, err
	/* TODO:
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}
		return secret, nil
	})
	*/
}
