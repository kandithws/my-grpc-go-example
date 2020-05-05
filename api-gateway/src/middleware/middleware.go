package middleware

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/kandithws/sharespace-api/api-gateway/src/utils"
	"github.com/labstack/echo/v4"
	echoMw "github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

type User struct {
	Username  string `json:"username"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"email"`
}

type JWTClaims struct {
	User
	jwt.StandardClaims
}

// https://echo.labstack.com/cookbook/jwt

func NewJWTAuthMiddleware() echo.MiddlewareFunc {
	config := echoMw.JWTConfig{
		Claims:     &JWTClaims{},
		SigningKey: []byte(viper.GetString("app.key")),
	}

	return echoMw.JWTWithConfig(config)
}

func GetAuthUser(c echo.Context) (*User, error) {
	userPayload := &User{}
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*JWTClaims)
	if err := utils.BindJSON(claims, userPayload); err != nil {
		return nil, err
	}

	return userPayload, nil
}
