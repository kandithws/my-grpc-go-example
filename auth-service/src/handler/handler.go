package handler

import (
	"context"

	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"google.golang.org/grpc"
)

func NewAuthGrpcHandler(gserver *grpc.Server) {
	s := &server{}

	genproto.RegisterAuthServiceServer(gserver, s)
}

// AuthServiceHandler implements genproto.AuthServiceServer
type server struct {
}

func (h *server) Register(context.Context, *genproto.RegisterRequest) (*genproto.RegisterResponse, error) {
	return nil, nil
}

func (h *server) Login(context.Context, *genproto.LoginRequest) (*genproto.LoginResponse, error) {
	return nil, nil
}

func (h *server) GetUser(context.Context, *genproto.GetUserRequest) (*genproto.User, error) {
	return nil, nil
}
