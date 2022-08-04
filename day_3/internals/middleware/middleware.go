package middleware

import (
	"elastic_study/internals/view"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"net/http"
	"strings"
	"time"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.RequestURI, time.Since(start))
	})
}

func CheckAuth(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := checkAuthorization([]byte(secret), r); err != nil {
			view.HandleErrorJSON(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func NewCheckAuthHandler(secret string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return CheckAuth(secret, h)
	}
}

func checkAuthorization(signingKey []byte, r *http.Request) error {
	var headerArr []string
	headerVal, ok := r.Header["Authorization"]
	if ok {
		headerArr = strings.Split(headerVal[0], " ")
	}
	if len(headerArr) == 2 && headerArr[0] == "Bearer" {
		err := parseToken(headerArr[1], signingKey)
		if err != nil {
			return err
		}
		return nil
	}
	return errors.New("invalid auth token")
}

func parseToken(accessToken string, signingKey []byte) error {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return signingKey, nil
	})

	if err != nil {
		return err
	}

	if token.Valid {
		return nil
	}

	return errors.New("invalid auth token")
}
