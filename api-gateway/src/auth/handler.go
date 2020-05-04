package auth

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	authproto "github.com/kandithws/sharespace-api/api-gateway/src/genproto/auth-service"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResponseError struct {
	Code    string `json:"rpc_error_code,omitempty"`
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

func MakeErrorResponse(c echo.Context, code int, err error) error {
	return c.JSON(code, &ResponseError{Message: err.Error()})
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := &authproto.RegisterRequest{}

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
	req := &authproto.LoginRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	res, err := h.authStore.ServiceClient().Login(context.Background(), req)

	if err != nil {
		return GetHttpErrorFromGRPCError(err, c)
	}

	if !res.Authorized {
		return echo.ErrUnauthorized
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	// TODO -- add user data to claims
	var userJSON map[string]interface{}
	if err := BindJSON(res.User, &userJSON); err != nil {
		return MakeErrorResponse(c, http.StatusInternalServerError, err)
	}

	claims["user"] = userJSON
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token": t,
	})
}
