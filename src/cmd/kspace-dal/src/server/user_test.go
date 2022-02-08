package server

import (
	"context"
	rpc2 "github.com/kercylan98/kspace/src/cmd/kspace-dal/src/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
	"time"
)

func TestUser_Create(t *testing.T) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("127.0.0.1:9500", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := rpc2.NewDalUserClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.Create(ctx, &rpc2.User{
		Account:  "test",
		Password: "123456",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetUser())
}
