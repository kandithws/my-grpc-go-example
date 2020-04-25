package main

import (
	"fmt"
	"net"

	authHandler "github.com/kandithws/sharespace-api/auth-service/src/handler"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Return Service Interface
	server := grpc.NewServer()

	authHandler.NewAuthGrpcHandler(server)
	reflection.Register(server)
	list, err := net.Listen("tcp", viper.GetString("server.address"))
	if err != nil {
		fmt.Println("SOMETHING HAPPEN")
	}
	err = server.Serve(list)
	if err != nil {
		fmt.Println("Unexpected Error", err)
	}
}
