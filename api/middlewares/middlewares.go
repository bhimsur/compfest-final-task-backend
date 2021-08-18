package middlewares

import (
	"context"
	"net/http"
	"os"
	"restgo/api/responses"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
)

func SetResponsesMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Requested-With, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Content-Type", "application/json")
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080, http://localhost:3000, https://pentapeduli.hexalogi.cyou/")
			w.Header().Set("Access-Control-Allow-Methods", "GET, PUT, POST, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, X-Requested-With, Content-Type")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Content-Type", "application/json")
		}
		next.ServeHTTP(w, r)
	})
}

func AuthJwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "message": "Missing authorization token"}

		var header = r.Header.Get("Authorization")
		header = strings.TrimSpace(header)

		if header == "" {
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(header, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET")), nil
		})
		if err != nil {
			resp["status"] = false
			resp["message"] = "Invalid token, please login"
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "Role", claims["Role"])
		ctx = context.WithValue(ctx, "UserID", claims["UserID"])
		ctx = context.WithValue(ctx, "Status", claims["Status"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
