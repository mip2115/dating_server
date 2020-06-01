package types

import (
	jwt "github.com/dgrijalva/jwt-go"
)

type Token struct {
	UserUUID string `json:"userUUID"`
	jwt.StandardClaims
}
