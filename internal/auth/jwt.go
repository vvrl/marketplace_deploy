package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	authHeader   = "Authorization"
	bearerPrefix = "Bearer "
)

type (
	JWTProvider interface {
		GenerateToken(userID int) (string, error)
		ParseIdFromToken(tokenString string) (int, error)
		ExtractToken(r *http.Request) (string, error)
	}

	jwtProvider struct {
		secretKey  string
		expiration time.Duration
	}
)

func NewJWTProvider(secret string, expirationHours int) JWTProvider {
	return &jwtProvider{
		secretKey:  secret,
		expiration: time.Hour * time.Duration(expirationHours),
	}
}

func (j *jwtProvider) GenerateToken(userID int) (string, error) {

	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(j.expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(j.secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (j *jwtProvider) ParseIdFromToken(tokenString string) (int, error) {
	// Парсинг токена с проверкой подписи
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secretKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, fmt.Errorf("expired token: %v", err) // Отдельная ошибка для просроченного токена
		}
		return 0, fmt.Errorf("invalid token: %v", err)
	}

	// Получение ID
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if userID, ok := claims["sub"].(float64); ok {
			return int(userID), nil
		}
	}

	return 0, errors.New("invalid token claims")
}

func (j *jwtProvider) ExtractToken(r *http.Request) (string, error) {
	authRequest := r.Header.Get(authHeader)
	if authRequest == "" {
		return "", errors.New("authorization header is needed")
	}

	if !strings.HasPrefix(authRequest, bearerPrefix) {
		return "", errors.New("invalid authorization header")
	}

	return authRequest[len(bearerPrefix):], nil
}
