package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type ResponseError struct {
	Message string `json:"message"`
}

type AuthHandler struct {
	authStore *AuthStore
}

func NewHttpHandler(e *echo.Group) {
	h := &AuthHandler{
		authStore: NewAuthStore(),
	}

	e.POST("/login", h.Login)
	e.GET("/hello_world", h.HelloWorld)
}

func (h *AuthHandler) HelloWorld(c echo.Context) error {
	return c.JSON(http.StatusOK, &ResponseError{Message: "Hello World!"})
}

func (h *AuthHandler) Login(c echo.Context) error {
	return c.JSON(http.StatusOK, &ResponseError{Message: "Logged In"})
}
