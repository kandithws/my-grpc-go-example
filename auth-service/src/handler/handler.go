package handler

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"

	"github.com/kandithws/sharespace-api/auth-service/src/common/validator"
	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"github.com/kandithws/sharespace-api/auth-service/src/model"
	"github.com/kandithws/sharespace-api/auth-service/src/store"
	"google.golang.org/grpc"
)

func NewAuthGrpcHandler(gserver *grpc.Server) {
	s := &handler{userStore: store.NewUserStore(), val: validator.NewValidator()}

	genproto.RegisterAuthServiceServer(gserver, s)
}

// AuthServiceHandler implements genproto.AuthServiceServer
type handler struct {
	userStore *store.UserStore
	val       validator.Validator
}

var (
	ErrRequestBindingError = errors.New("Fail to bind the request")
)

func bindRequest(req interface{}, m interface{}) error {
	if reflect.TypeOf(m).Kind() != reflect.Ptr {
		return ErrRequestBindingError
	}

	byt, _ := json.Marshal(req)
	if err := json.Unmarshal(byt, m); err != nil {
		return err
	}
	return nil
}

func (h *handler) Register(c context.Context, req *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {
	var m *model.User
	// Bind
	if err := bindRequest(req, m); err != nil {
		return nil, err
	}

	if err := h.val.Validate(m); err != nil {
		return nil, err
	}
	// h.userStore.CreateUser()
	return &genproto.RegisterResponse{}, nil
}

func (h *handler) Login(context.Context, *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	return &genproto.LoginResponse{}, nil
}

func (h *handler) GetUser(context.Context, *genproto.GetUserRequest) (*genproto.User, error) {
	return &genproto.User{}, nil
}
