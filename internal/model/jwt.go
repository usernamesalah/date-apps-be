package model

import (
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserUID string `json:"user_uid"`
	jwt.RegisteredClaims
}

func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {

	return c.Audience, nil

}
