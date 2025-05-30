package middleware

import (
	"Praiseson6065/ocrolus-be/config"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

type JWTClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}

func GenerateToken(userId string) (string, error) {
	// Get JWT settings from config
	jwtExpiration := config.Config.JWT.Expire
	signingKey := []byte(config.Config.JWT.Secret)

	claims := JWTClaims{
		userId,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(jwtExpiration) * time.Hour).Unix(),
			IssuedAt:  jwt.TimeFunc().Unix(),
			Issuer:    "ocrolus",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(signingKey)

	return tokenString, err
}

func ValidateToken(encodedToken string) (string, error) {
	// Get signing key from config
	signingKey := []byte(config.Config.JWT.Secret)

	claims := &JWTClaims{}

	_, err := jwt.ParseWithClaims(encodedToken, claims, func(t *jwt.Token) (interface{}, error) {
		if _, isValid := t.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("invalid token %s", t.Header["alg"])
		}
		return signingKey, nil
	})
	if err != nil {
		return "", err
	}
	return claims.UserId, nil
}
