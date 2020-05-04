package auth

import (
	"github.com/kandithws/sharespace-api/api-gateway/src/genproto"
	"google.golang.org/grpc"
)

type AuthStore struct {
	client genproto.AuthServiceClient
}

func initAuthServiceClient(url string) (genproto.AuthServiceClient, error) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return genproto.NewAuthServiceClient(conn), nil
}

func NewAuthStore(serviceURL string) *AuthStore {
	cl, err := initAuthServiceClient(serviceURL)
	if err != nil {
		panic(err)
	}
	return &AuthStore{client: cl}
}

func (s *AuthStore) ServiceClient() genproto.AuthServiceClient {
	return s.client
}
