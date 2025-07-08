package jwt

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt"
)

const defaultExpirationTimeToken = 3600

type JWTGenerator struct{}

func (j JWTGenerator) GenerateToken(email string) (string, string, error) {
	key := os.Getenv("SECRET_KEY")
	secretKey := []byte(key)

	expirationTimeStr := os.Getenv("TIME_TOKEN")
	expirationTimeDuration, err := strconv.Atoi(expirationTimeStr)
	if err != nil {
		expirationTimeDuration = defaultExpirationTimeToken
	}

	refreshExpirationTimeStr := os.Getenv("TIME_REFRESH_TOKEN")
	refreshExpirationTimeDuration, err := strconv.Atoi(refreshExpirationTimeStr)
	if err != nil {
		refreshExpirationTimeDuration = defaultExpirationTimeToken
	}

	expirationTime := time.Now().Add(time.Duration(expirationTimeDuration) * time.Minute)
	refreshExpirationTime := time.Now().Add(time.Duration(refreshExpirationTimeDuration) * time.Minute)

	claims := &jwt.StandardClaims{
		ExpiresAt: expirationTime.Unix(),
		Subject:   email,
	}
	refreshClaims := &jwt.StandardClaims{
		ExpiresAt: refreshExpirationTime.Unix(),
		Subject:   email,
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(secretKey)
	if err != nil {
		return "", "", err
	}
	return token, refreshToken, nil
}

func (j JWTGenerator) ValidateToken(token string) (bool, error) {
	key := os.Getenv("SECRET_KEY")
	secretKey := []byte(key)

	parsedToken, err := jwt.ParseWithClaims(token, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return secretKey, nil
	})

	if err != nil || !parsedToken.Valid {
		return false, err
	}

	claims, ok := parsedToken.Claims.(*jwt.StandardClaims)
	if !ok {
		return false, err
	}

	if claims.ExpiresAt < time.Now().Unix() {
		return false, errors.New("token has expired")
	}

	return true, nil
}
