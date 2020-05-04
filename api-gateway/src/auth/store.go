package auth

import (
	authproto "github.com/kandithws/sharespace-api/api-gateway/src/genproto/auth-service"
	"google.golang.org/grpc"
)

type AuthStore struct {
	client authproto.AuthServiceClient
}

func initAuthServiceClient(url string) (authproto.AuthServiceClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return authproto.NewAuthServiceClient(conn), nil
}

func NewAuthStore(serviceURL string) *AuthStore {
	cl, err := initAuthServiceClient(serviceURL)
	if err != nil {
		panic(err)
	}
	return &AuthStore{client: cl}
}

func (s *AuthStore) ServiceClient() authproto.AuthServiceClient {
	return s.client
}
