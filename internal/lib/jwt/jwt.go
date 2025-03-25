package jwtauth

import (
	"SSO/internal/domain/models"
	"SSO/pkg/utils"
	"time"
)

func NewToken(user models.User,
	app models.App,
	duration time.Duration,
	refreshDuration time.Duration) (string, string, error) {
	accessToken, refreshToken, err := utils.GenerateToken(user.ID, user.Username, app, duration, refreshDuration)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

	//token := jwt.New(jwt.SigningMethodHS256)
	//
	//claims := token.Claims.(jwt.MapClaims)
	//claims["uid"] = user.ID
	//claims["email"] = user.Email
	//claims["exp"] = time.Now().Add(duration).Unix()
	//claims["app_id"] = app.ID
	//
	//tokenString, err := token.SignedString([]byte(app.Secret))
	//if err != nil {
	//	return "", err
	//}
	//
	//return tokenString, nil
}
