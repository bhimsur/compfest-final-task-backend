package utils

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func EncodeAuthToken(uid uint, email string, role string, status string) (string, error) {
	claims := jwt.MapClaims{}
	claims["UserID"] = uid
	claims["Email"] = email
	claims["Role"] = role
	claims["Status"] = status
	claims["IssuedAt"] = time.Now().Unix()
	claims["ExpiresAt"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS512"), claims)
	return token.SignedString([]byte(os.Getenv("SECRET")))
}
