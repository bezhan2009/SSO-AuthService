package utils

import (
	"SSO/internal/domain/models"
	"SSO/pkg/errs"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"os"
	"time"
)

// CustomClaims определяет кастомные поля токена
type CustomClaims struct {
	UserID   int    `json:"user_id"`
	Username string `json:"username"`
	AppID    int    `json:"app_id"`
	jwt.StandardClaims
}

// GenerateToken генерирует JWT токен с кастомными полями
func GenerateToken(userID uint,
	username string,
	app models.App,
	duration time.Duration,
	refreshDuration time.Duration) (string, string, error) {
	// Access token
	claims := &CustomClaims{
		UserID:   int(userID),
		Username: username,
		AppID:    int(app.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(duration).Unix(), // токен истекает через 1 час
			Issuer:    app.Name,
		},
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessTokenString, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	// Refresh token
	refreshClaims := &CustomClaims{
		UserID:   int(userID),
		Username: username,
		AppID:    int(app.ID),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(refreshDuration).Unix(), // токен истекает через 72 часа
			Issuer:    app.Name,
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil {
		return "", "", err
	}

	return accessTokenString, refreshTokenString, nil
}

// ParseToken парсит JWT токен и возвращает кастомные поля
func ParseToken(tokenString string, tokenSecret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Проверяем метод подписи токена
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errs.ErrInvalidToken
}
