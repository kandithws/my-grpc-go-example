package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	authHandler "github.com/kandithws/sharespace-api/auth-service/src/handler"
	"github.com/kandithws/sharespace-api/common/db"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initConfig() {
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("GO_ENV", "development")

	viper.AutomaticEnv()
}

func main() {
	// Return Service Interface
	initConfig()
	server := grpc.NewServer()
	dbCfg := db.NewDBConfig()
	dbCfg.Username = "kandithws"
	dbCfg.Password = "gunto1166"
	dbCfg.DatabaseName = "auth_service"
	db.InitDB(dbCfg)
	authHandler.NewAuthGrpcHandler(server)
	reflection.Register(server)
	uri := fmt.Sprintf(":%s", viper.GetString("PORT"))
	list, err := net.Listen("tcp", uri)
	if err != nil {
		panic("Error on setting up Port")
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx := context.Background()
	go func() {
		for range c {
			// sig is a ^C, handle it
			fmt.Println("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	fmt.Printf("Serving on grpc server on %s\n", uri)
	server.Serve(list)
}
