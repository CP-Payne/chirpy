package main

import (
	"encoding/base64"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func (cfg *apiConfig) CreateToken(userID int, expiresInSeconds int64) (string, error) {
	now := time.Now().UTC()

	const maxExpirationInSeconds = 24 * 60 * 60
	if expiresInSeconds <= 0 || expiresInSeconds > maxExpirationInSeconds {
		expiresInSeconds = maxExpirationInSeconds
	}

	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expiresInSeconds) * time.Second)),
		Subject:   fmt.Sprintf("%d", userID),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret, err := base64.StdEncoding.DecodeString(cfg.jwtSecret)
	if err != nil {
		fmt.Println("Error decoding secret key:", err)
		return "", err
	}
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Println(cfg.jwtSecret)
		return "", err
	}

	return signedToken, nil
}
