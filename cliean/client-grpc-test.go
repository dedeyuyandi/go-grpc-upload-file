package main

import (
	"context"
	"fmt"
	"log"

	v1_user_grpc "github.com/dedeyuyandi/go-grpc-upload-file/proto"
	"google.golang.org/grpc"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":7000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %s", err.Error())
	}
	defer conn.Close()

	c := v1_user_grpc.NewUserClient(conn)
	resp, err := c.CreateHelloWorld(context.Background(), &v1_user_grpc.HelloWorld{
		Hello: "hello",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.String())
}
