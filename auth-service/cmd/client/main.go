package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/kandithws/sharespace-api/auth-service/src/genproto"
	"google.golang.org/grpc"
)

var (
	host        = flag.String("host", ":8081", "Server address, e.g. :8081")
	requestName = flag.String("request", "Login", "gRPC method to call")
	body        = flag.String("body", "", "Message Body (use if given)")
)

func main() {
	flag.Parse()

	conn, err := grpc.Dial(*host, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := genproto.NewAuthServiceClient(conn)

	// TODO dynamic
	req := &genproto.LoginRequest{
		Username: "test",
		Password: "test",
	}
	res, err := client.Login(context.Background(), req)

	if err != nil {
		log.Fatal(err)
	}

	if res != nil {
		//fmt.Printf("%d task(s)\n", res.Total)
		fmt.Printf("%v\n", res)
	} else {
		fmt.Print("No tasks\n")
	}
}
