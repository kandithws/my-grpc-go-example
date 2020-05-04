package main

import (
	"fmt"

	"net"
	"os"
	"os/signal"
	"strings"

	"github.com/kandithws/sharespace-api/auth-service/src/common/db"
	authHandler "github.com/kandithws/sharespace-api/auth-service/src/handler"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
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
	viper.SetDefault("PORT", "8082")
	viper.SetDefault("GO_ENV", "development")

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "__"))
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error loading config file: %s", err))
	}
	viper.Set("app.root", appdir)
}

func initDefaultLogger() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
}

func NewGRPCServerLogger() *logrus.Entry {
	serverLogger := log.New()
	serverLogger.SetFormatter(&log.JSONFormatter{})
	serverLogger.SetReportCaller(true)
	return log.NewEntry(serverLogger)
}

func makeDBConfig() *db.DBConfig {
	dbCfg := db.NewDBConfig()
	viper.SetDefault("db.host", "127.0.0.1")
	viper.SetDefault("db.port", "5432")
	dbCfg.Host = viper.GetString("db.host")
	dbCfg.Port = viper.GetString("db.port")
	dbCfg.Username = viper.GetString("db.username")
	dbCfg.Password = viper.GetString("db.password")
	dbCfg.DatabaseName = viper.GetString("db.db_name")
	return &dbCfg
}

func main() {
	// Return Service Interface
	initDefaultLogger()
	initConfig()
	logrusEntry := NewGRPCServerLogger()
	server := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logrusEntry),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logrusEntry),
		),
	)

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
			log.Infof("shutting down gRPC server...")

			server.GracefulStop()

			<-ctx.Done()
		}
	}()

	log.Infof("Serving on grpc server on %s\n", uri)
	server.Serve(list)
}
