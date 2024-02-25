package security

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
)

type JwtToken struct {
	ID      uuid.UUID
	Name    string
	Email   string
	IsAdmin bool
	jwt.StandardClaims
}

type ContextKey struct{}

var secretKeyJWT = []byte(os.Getenv("JWT_SECRET"))

func SignJWT(claims JwtToken) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Minute * 10).Unix()

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tokenObj.SignedString(secretKeyJWT)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Verify(token string) (*JwtToken, error) {
	var userClaim = &JwtToken{}

	tokenObj, err := jwt.ParseWithClaims(token, userClaim, func(t *jwt.Token) (interface{}, error) {
		return secretKeyJWT, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return nil, errors.New("unauthorized access")
		}

		return nil, errors.New("status bad request")
	}

	if !tokenObj.Valid {
		return nil, errors.New("token invalid")
	}

	return userClaim, nil
}
