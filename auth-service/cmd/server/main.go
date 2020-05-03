package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/kandithws/sharespace-api/auth-service/src/common/db"
	authHandler "github.com/kandithws/sharespace-api/auth-service/src/handler"

	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initConfig() {

	appdir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(appdir)
	viper.SetDefault("PORT", "8081")
	viper.SetDefault("GO_ENV", "development")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config file: %s", err))
	}
	viper.Set("app.root", appdir)
}

func makeDBConfig() *db.DBConfig {
	dbCfg := db.NewDBConfig()
	dbCfg.Username = viper.GetString("db.username")
	dbCfg.Password = viper.GetString("db.password")
	dbCfg.DatabaseName = viper.GetString("db.db_name")
	return &dbCfg
}

func main() {
	// Return Service Interface
	initConfig()
	server := grpc.NewServer()

	db.InitDB(makeDBConfig())
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
