package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func GenerateDevToken(c echo.Context) string {
	claims := jwt.MapClaims{
		"id":       1,
		"username": "admin",
		"exp":      time.Now().Add(time.Hour * 72).Unix(),
	}
	config := c.Get("config").(Config)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(config.JWTScrect))

	return "Bearer " + signedToken
}

// 开发环境专用中间件：自动注入 JWT Token
func DevJWTInjector(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		// 仅限开发环境启用
		if !Env.Debug {
			return next(c)
		}
		authHeader := c.Request().Header.Get(echo.HeaderAuthorization)
		if authHeader == "" {
			c.Request().Header.Set(echo.HeaderAuthorization, GenerateDevToken(c))
		}

		return next(c)
	}
}
