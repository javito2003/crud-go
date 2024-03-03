package middlewares

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/javito2003/crud-go/models"
	"github.com/javito2003/crud-go/server"
)

func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JwtSecret), nil
			})

			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
