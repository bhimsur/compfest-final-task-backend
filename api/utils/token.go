package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang-module/carbon"
)

const (
	rfcnano    = "2021-08-16 22:28:56.362787+07"
	carbonrfc  = "2021-12-11 17:00:00 +0000 UTC"
	layoutDate = "2021-08-12"
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

func RFCToDate(date time.Time) string {
	return date.String()[0:10]
}

func DateToRFC(date string) time.Time {
	fmt.Printf("\ndate : %s\n", carbon.Parse(date).ToRfc3339String())
	t, _ := time.Parse(time.RFC3339Nano, carbon.Parse(date).ToRfc3339String())
	fmt.Println(t)
	return t
}
