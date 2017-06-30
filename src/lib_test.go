package microservice

import (
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/server"
	"github.com/quaponatech/golang-extensions/test"
)

var (
	startPort  = 10000
	portCount  = 0
	serverName = "MicroService test"
)

// Returns the next available port for testing
func getNextPort() int {
	port := startPort + portCount
	portCount++
	return port
}

// Returns a default gRPC server instance
func getGRPCServer() *grpcservice.GRPCServer {
	return grpcservice.NewGRPCServer(false, "", "", getNextPort())
}

// Returns a default server logger instance
func getLogger() *server.Logger {
	return server.NewLogger(
		serverName,
		"",
		"",
		make(chan server.Status),
		make(chan error),
		make(chan string),
		make(chan string),
		make(chan string),
		0)
}

// Returns a default microservice instance
func getDefaultService(t *testing.T) MicroService {
	service := MicroService{}
	client := &GrpcClientMicroService{}
	test.AssertThat(t, service.Setup(serverName, getGRPCServer(), getLogger(), client), nil)
	return service
}

func TestSuccessSetupServeStop(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		// Prepare the environment
		service := getDefaultService(t)

		// Test Serve
		go func() {
			err := service.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(time.Second)

		// Test Stop
		err := service.Stop()
		test.AssertThat(t, err, nil)
	})
}

func TestFailSetupGRPCServerNil(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		service := MicroService{}
		test.AssertThat(t, service.Setup(serverName, nil, nil, nil), nil, "not")
	})
}
