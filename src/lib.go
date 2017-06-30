// Package microservice provides ... TODO
// This microservice should ... TODO
package microservice

import (
	"fmt"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/server"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

// Version of the Library
const Version = "0.1.0"

// MicroService structure as main entry point
type MicroService struct {
	*grpcservice.GRPCService
	// TODO: A flag to prepare the service for special test environment
	// test bool
}

// Setup the microservice
func (service *MicroService) Setup(
	serverName string,
	serverInstance *grpcservice.GRPCServer,
	serverLogger *server.Logger,
	client *GrpcClientMicroService,
) error {
	// Set the internal variables
	service.GRPCService = &grpcservice.GRPCService{}

	// Setup the internal gRPC service
	if err := service.GRPCService.Setup(serverName, serverInstance, serverLogger, make(chan bool)); err != nil {
		return fmt.Errorf(service.Prefix + err.Error())
	}

	// Set the server state to 'starting'
	service.StatusChan <- server.StateStarting
	service.LogChan <- "Register microservice"

	// Register the protocol buffers server
	protobuf.RegisterMicroServiceServer(service.GetInstance(), service)

	// Connect to other microservices
	if client == nil {
		err := fmt.Errorf("Microservice Client not initialized")
		service.ErrorChan <- err
		service.StatusChan <- server.StateError
		if sErr := service.Stop(); sErr != nil {
			return sErr
		}
		return fmt.Errorf(service.Prefix + err.Error())
	}

	/* DEMO part start
	// TODO: Connect to another microservice
	if client.AccessInfo != nil { // The real connection path
		service.LogChan <- "Connect to microservice"
		if err := client.Connect(client.AccessInfo); err != nil {
			service.ErrorChan <- err
			service.StatusChan <- server.StateError
			service.Stop()
			return nil
		}
		service.LogChan <- "Successfully connected to microservice!"
	} else if client.MicroServiceClient == nil { // The failure path
		service.ErrorChan <- fmt.Errorf("Microservice client equals nil")
		service.StatusChan <- server.StateError
		service.Stop()
		return nil
	} else { // The mocking path
		service.LogChan <- "Use connected microservice"
		// TODO: Use this path to inject mocks or already connected clients
	}
	DEMO part end */

	// Set the server state to 'started'
	service.StatusChan <- server.StateStarted
	return nil
}

// Serve the microservice
func (service *MicroService) Serve() error {
	return service.GRPCService.Serve()
}

// Stop the microservice
func (service *MicroService) Stop() error {
	// TODO: Close connections to other microservices ...

	// Stop the server
	return service.GRPCService.Stop()
}
