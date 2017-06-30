package main

import (
	"context"
	"fmt"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/microservice-template/src"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

func main() {
	// Setup a new client
	client := &microservice.GrpcClientMicroService{}

	// Create a connection info
	info := &grpcservice.ConnectionInfo{IP: "localhost", Port: "42302"}

	// Connect to the instance
	if err := client.Connect(info); err != nil {
		fmt.Printf("Could not connect: %v\n", err)
		return
	}

	// Send a message via the RPC
	message := &protobuf.Request{Message: "Hello, world!"}
	fmt.Printf("Sending %v\n", message)
	response, err := client.MicroServiceClient.Hello(context.Background(), message)

	// Check if the RPC was successful
	if err != nil {
		fmt.Printf("Could not say hello: %v\n", err)
		return
	}

	fmt.Printf("Got response %v\n", response)

	// Close the connection
	if err := client.Close(); err != nil {
		fmt.Printf("Could not close connection: %v\n", err)
		return
	}
	fmt.Println("Closed connection.")
}
