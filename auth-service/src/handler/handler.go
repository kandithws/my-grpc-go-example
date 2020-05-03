package handler

import (
	"context"

	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"github.com/kandithws/sharespace-api/auth-service/src/store"
	"google.golang.org/grpc"
)

func NewAuthGrpcHandler(gserver *grpc.Server) {
	s := &handler{userStore: store.NewUserStore()}

	genproto.RegisterAuthServiceServer(gserver, s)
}

// AuthServiceHandler implements genproto.AuthServiceServer
type handler struct {
	userStore *store.UserStore
}

func (h *handler) Register(context.Context, *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {

	return &genproto.RegisterResponse{}, nil
}

func (h *handler) Login(context.Context, *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	return &genproto.LoginResponse{}, nil
}

func (h *handler) GetUser(context.Context, *genproto.GetUserRequest) (*genproto.User, error) {
	return &genproto.User{}, nil
}
