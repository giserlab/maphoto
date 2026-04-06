package jwt

import (
	"errors"
	"maphoto/internal/app/model"
	"maphoto/internal/app/store"

	"github.com/golang-jwt/jwt/v5"
)

func ParseJWT(data interface{}) (*model.User, error) {
	if Env.Debug {
		user := model.User{}
		err := store.DB.Preload("Config").Preload("Places").Where("admin = ?", true).First(&user, 1).Error
		return &user, err
	}
	if data == nil {
		return nil, errors.New("token data is none")
	}

	token, ok := data.(*jwt.Token)
	if !ok {
		return nil, errors.New("JWT token missing or invalid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("failed to cast claims as jwt.MapClaims")
	}
	user := model.User{}
	idVal, ok := claims["id"]
	if !ok {
		return nil, errors.New("id claim missing")
	}
	id, ok := idVal.(float64)
	if !ok {
		return nil, errors.New("id claim is not a number")
	}
	err := store.DB.Preload("Config").Preload("Places").First(&user, int64(id)).Error
	if err != nil {
		return nil, err
	}
	if user.Username == "" {
		return nil, errors.New("user not exist")
	}
	if *user.Token != token.Raw {
		return nil, errors.New("token not match")
	}
	return &user, nil
}
