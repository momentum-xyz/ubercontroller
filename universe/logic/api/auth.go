package api

import (
	"bytes"
	"fmt"
	"strings"
	"time"

	"github.com/ChainSafe/go-schnorrkel"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/momentum-xyz/ubercontroller/database"
	"github.com/momentum-xyz/ubercontroller/types/entry"
	"github.com/momentum-xyz/ubercontroller/universe"
	"github.com/momentum-xyz/ubercontroller/utils"
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

func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

func VerifyEthereumSignature(address, challenge, signature string) (bool, error) {
	if !strings.HasPrefix(signature, "0x") {
		signature = "0x" + signature
	}
	sigBytes := hexutil.MustDecode(signature)
	if len(sigBytes) <= 64 || (sigBytes[64] != 27 && sigBytes[64] != 28) {
		return false, errors.New("unsupported signature format")
	}

	sigBytes[64] -= 27

	msgBytes := signHash([]byte(challenge))

	pubKey, err := crypto.SigToPub(msgBytes, sigBytes)
	if err != nil || pubKey == nil {
		return false, errors.Wrap(err, "failed to recover public key")
	}

	recoveredAddr := crypto.PubkeyToAddress(*pubKey).Hex()

	if recoveredAddr != address {
		return false, errors.New("the challenge was not signed by the correct address")
	}

	return true, nil
}

// CreateJWTToken saves a jwt token with the given userID as subject
// and signed with the given secret
func CreateJWTToken(userID uuid.UUID) (string, error) {
	claims := jwt.StandardClaims{
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour).Unix(),
		Issuer:    "ubercontroller",
		Subject:   userID.String(),
	}

	newJwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret, err := GetJWTSecret()
	if err != nil {
		return "", errors.WithMessage(err, "failed to get jwt secret")
	}

	signedString, err := newJwt.SignedString(secret)
	if err != nil {
		return "", errors.WithMessage(err, "failed to sign token")
	}

	return signedString, nil
}

func GetJWTSecret() ([]byte, error) {
	jwtSecret, ok := universe.GetNode().GetNodeAttributes().GetValue(
		entry.NewAttributeID(universe.GetSystemPluginID(), universe.ReservedAttributes.Node.JWTKey.Name),
	)
	if !ok || jwtSecret == nil {
		return nil, errors.New("failed to get jwt secret")
	}
	secret := utils.GetFromAnyMap(*jwtSecret, universe.ReservedAttributes.Node.JWTKey.Key, "")

	return []byte(secret), nil
}

func ValidateJWTWithSecret(signedString string, secret []byte) (*jwt.Token, error) {
	return jwt.Parse(signedString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Invalid token %v", token.Header["alg"])
		}
		return secret, nil
	})
}

func GenerateGuestName(c *gin.Context, db database.DB) (string, error) {
	visitorNameTemplate := "Visitor_"

	visitorSuffix, err := gonanoid.Generate("0123456789", 7)
	if err != nil {
		return "", errors.WithMessage(err, "failed to generate visitor name")
	}

	visitorName := visitorNameTemplate + visitorSuffix

	exists, err := db.GetUsersDB().CheckIsUserExistsByName(c, visitorName)
	if err != nil {
		return "", errors.WithMessage(err, "failed to check for duplicates")
	}
	if exists {
		// Row exists -> regenerate
		GenerateGuestName(c, db)
	}

	return visitorName, nil
}
