package authservice

import (
	"crypto/rsa"
	"date-apps-be/infrastructure/config"
	"date-apps-be/internal/model"
	"date-apps-be/pkg/derrors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	AuthService interface {
		GenerateToken(uid string) (_ string, err error)
	}

	authService struct {
		privateKey *rsa.PrivateKey
		expiration int
	}
)

func NewAuthService(conf *config.Config) AuthService {
	return &authService{
		privateKey: conf.JWTRS256PrivateKey,
		expiration: conf.JWTExpiration,
	}
}

func (a *authService) newJWTClaims(uid string) *model.JWTClaims {
	fmt.Println(a.expiration)
	expirationTime := time.Now().Add(time.Duration(a.expiration) * time.Minute)
	fmt.Println("expirationTime ", expirationTime)
	claims := &model.JWTClaims{
		UserUID: uid,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "date-apps",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	return claims
}

func (a *authService) GenerateToken(uid string) (_ string, err error) {
	defer derrors.Wrap(&err, "GenerateToken(%q)", uid)

	claims := a.newJWTClaims(uid)
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(a.privateKey)
	if err != nil {
		return "", derrors.WrapStack(err, derrors.Unknown, "token.SignedString")
	}

	return tokenString, nil
}
