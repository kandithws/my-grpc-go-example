package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"

	"github.com/kandithws/sharespace-api/api-gateway/src/auth"
)

func initEcho() *echo.Echo {
	e := echo.New()

	e.Pre(middleware.RemoveTrailingSlash())

	e.Use(middleware.Logger())

	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	return e
}

func initConfig() {
	appdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(appdir)
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("GO_ENV", "development")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config file: %s", err))
	}

	viper.Set("app.root", appdir)
}

func main() {
	initConfig()
	e := initEcho()
	aGroup := e.Group("/auth")
	aStore := auth.NewAuthStore(viper.GetString("auth_service.url"))
	auth.NewHttpHandler(aGroup, aStore)
	p := fmt.Sprintf(":%s", viper.GetString("PORT"))
	e.Logger.Fatal(e.Start(p))
}
