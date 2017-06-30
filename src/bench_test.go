package microservice

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

// TODO: Add your own benchmarks
var port = 13333
var connections = 0
var messages = 0
var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GetRandString(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func prepareClient(b *testing.B) GrpcClientMicroService {
	// Start a server
	port++
	service := MicroService{}
	client := &GrpcClientMicroService{}
	test.AssertThat(b, service.Setup(serverName,
		grpcservice.NewGRPCServer(false, "", "", port),
		getLogger(),
		client),
		nil)

	go func() {
		err := service.Serve()
		test.AssertThat(b, err, nil)
	}()
	time.Sleep(time.Second)

	// Create a connection info
	info := &grpcservice.ConnectionInfo{
		UseTLS:              false,
		CertFile:            "",
		ServerHostName:      "",
		IP:                  "127.0.0.1",
		Port:                fmt.Sprint(port),
		TimeoutInMilliSecs:  1000,
		RetryTimes:          1,
		RetryAfterMilliSecs: 100,
	}

	// Connect encrypted, which should not work
	c := GrpcClientMicroService{}
	err := c.Connect(info)
	test.AssertThat(b, err, nil)
	connections++

	return c
}

// A simple performance test for the client
func BenchmarkClient(b *testing.B) {
	// Get a new client instance
	c := prepareClient(b)
	msg := &protobuf.Request{Message: GetRandString(1450)}

	// Run the test
	for n := 0; n < b.N; n++ {
		messages++
		_, err := c.MicroServiceClient.Hello(context.Background(), msg)
		test.AssertThat(b, err, nil)
	}

	fmt.Printf("Messages: %v\n", messages)
	fmt.Printf("Connections: %v\n", connections)
}

// A simple performance test for the client
func BenchmarkClientParallel(b *testing.B) {
	// Get a new client instance
	c := prepareClient(b)
	msg := &protobuf.Request{Message: GetRandString(1450)}

	// Run the test
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			messages++
			_, err := c.MicroServiceClient.Hello(context.Background(), msg)
			test.AssertThat(b, err, nil)
		}
	})

	fmt.Printf("Messages: %v\n", messages)
	fmt.Printf("Connections: %v\n", connections)
}
