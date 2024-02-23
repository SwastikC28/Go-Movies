package auth

import (
	"errors"
	"net/http"
	"os"
	"shared/utils/web"
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	ID      uint
	Name    string
	Email   string
	IsAdmin bool
	jwt.StandardClaims
}

type contextKey string

var secretKeyJWT = []byte(os.Getenv("JWT_SECRET"))

const UserIDKey contextKey = "userid"

func SignJWT(claims Claims) (string, error) {
	claims.StandardClaims.ExpiresAt = time.Now().Add(time.Minute * 10).Unix()

	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := tokenObj.SignedString(secretKeyJWT)

	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func Verify(token string) (*Claims, error) {
	var userClaim = &Claims{}

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

func IsAdmin(h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")

		user := r.Context().Value(UserIDKey).(*Claims)

		if !user.IsAdmin {
			web.RespondJSON(w, http.StatusUnauthorized, "User Unauthorized to access this Route")
			return
		}

		h.ServeHTTP(w, r)
	})
}
