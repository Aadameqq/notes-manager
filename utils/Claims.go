package utils

import "github.com/golang-jwt/jwt/v4"

type Claims struct {
	Id       string `json:"id"`
	Username string `json:"username"`
	*jwt.RegisteredClaims
}
