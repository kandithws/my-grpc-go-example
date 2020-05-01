package handler

import (
	"context"

	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"github.com/kandithws/sharespace-api/auth-service/src/store"
	"google.golang.org/grpc"
)

func NewAuthGrpcHandler(gserver *grpc.Server) {
	s := &server{userStore: store.NewUserStore()}

	genproto.RegisterAuthServiceServer(gserver, s)
}

// AuthServiceHandler implements genproto.AuthServiceServer
type server struct {
	userStore *store.UserStore
}

func (h *server) Register(context.Context, *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {
	return &genproto.RegisterResponse{}, nil
}

func (h *server) Login(context.Context, *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	return &genproto.LoginResponse{}, nil
}

func (h *server) GetUser(context.Context, *genproto.GetUserRequest) (*genproto.User, error) {
	return &genproto.User{}, nil
}
