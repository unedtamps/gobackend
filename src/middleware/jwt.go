package middleware

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
	"github.com/unedtamps/gobackend/utils"
)

type Login struct {
	Id    uuid.UUID `json:"id"`
	Email string    `json:"email"`
}

var TokenAuth *jwtauth.JWTAuth

func NewJwt() *jwtauth.JWTAuth {
	key := fmt.Sprintf("%s", os.Getenv("JWT_SECRET"))
	if key == "" {
		panic("JWT_SECRET is not set")
	}
	return jwtauth.New(
		"HS256",
		[]byte(key),
		nil,
	)
}

func SetJwt() {
	TokenAuth = NewJwt()
}

func getTokenString(r *http.Request) string {
	tokenString := jwtauth.TokenFromHeader(r)
	if tokenString == "" {
		tokenString = jwtauth.TokenFromQuery(r)
	}
	if tokenString == "" {
		tokenString = jwtauth.TokenFromCookie(r)
	}
	return tokenString
}

func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := getTokenString(r)
		if tokenString == "" {
			utils.ResponseError(w, 400, jwtauth.ErrNoTokenFound)
			return
		}
		_, err := jwtauth.VerifyToken(TokenAuth, tokenString)
		if err != nil {
			utils.ResponseError(w, 400, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}
