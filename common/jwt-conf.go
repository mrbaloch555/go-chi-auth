package common

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID int `json:"id"`
	jwt.StandardClaims
}

type JWTOutput struct {
	Token   string    `json:"token"`
	Expires time.Time `json:"expires"`
}
