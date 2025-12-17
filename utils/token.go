package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateJWT(userID string, secret string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "expense",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(5 * time.Minute)),
		Subject:   userID,
	})

	token, err := claims.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return token, nil
}

func ValidateJWT(tokenString, secret string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok {
		userID, err := claims.GetSubject()
		if err != nil {
			return "", err
		}
		return userID, nil
	}
	return "", errors.New("invalid token claims")
}

func GenerateRefresh() (string, error) {
	key := make([]byte, 32)
	rand.Read(key)
	token := hex.EncodeToString(key)
	if token == "" {
		return "", errors.New("error generating refresh token")
	}
	return token, nil
}

func GetAuthorizationToken(headers http.Header) (string, error) {
	authorization := headers.Get("Authorization")
	if authorization == "" {
		return "", errors.New("missing authorization header")
	}
	prefix := "Bearer "
	if !strings.HasPrefix(authorization, prefix) {
		return "", errors.New("invalid authorization header format")
	}
	token := strings.TrimPrefix(authorization, prefix)
	if token == "" {
		return "", errors.New("invalid token")
	}
	return token, nil
}
