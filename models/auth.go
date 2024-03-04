package models

import "github.com/golang-jwt/jwt"

type UserClaims struct {
	UserID string `json:"adminId"`

	jwt.StandardClaims
}
