package auth

import (
	"context"
	"net/http"

	"github.com/kandithws/sharespace-api/api-gateway/src/genproto"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseError struct {
	Code    string `json:"rpc_error_code"`
	Message string `json:"message"`
}

type AuthHandler struct {
	authStore *AuthStore
}

func NewHttpHandler(e *echo.Group, s *AuthStore) {
	h := &AuthHandler{
		authStore: s,
	}

	e.POST("/login", h.Login)
	e.POST("/register", h.Register)
}

func GetHttpErrorFromGRPCError(err error, c echo.Context) error {
	s, _ := status.FromError(err)

	if s.Code() == codes.Unauthenticated {
		return c.JSON(http.StatusUnauthorized, &ResponseError{Code: s.Code().String(), Message: s.Message()})
	}

	if s.Code() == codes.NotFound {
		return c.JSON(http.StatusNotFound, &ResponseError{Code: s.Code().String(), Message: s.Message()})
	}

	return c.JSON(http.StatusInternalServerError, &ResponseError{Code: s.Code().String(), Message: s.Message()})
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := &genproto.RegisterRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	res, err := h.authStore.ServiceClient().Register(context.Background(), req)
	if err != nil {
		return GetHttpErrorFromGRPCError(err, c)
	}

	return c.JSON(http.StatusOK, res)
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := &genproto.LoginRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	res, err := h.authStore.ServiceClient().Login(context.Background(), req)

	if err != nil {
		return GetHttpErrorFromGRPCError(err, c)
	}

	return c.JSON(http.StatusOK, res)
}
