package microservice

import (
	"fmt"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

// GrpcClientMicroService extends a generalized gRPC Client by functions to
// specialize it as MicroService client
type GrpcClientMicroService struct {
	GrpcClient         *grpcservice.GRPCClient
	AccessInfo         *grpcservice.ConnectionInfo
	MicroServiceClient protobuf.MicroServiceClient
}

// Connect sets up a connection to a gRPC Server by given parameters and
// returns an instance of a registered gRPC Service Client
func (client *GrpcClientMicroService) Connect(info *grpcservice.ConnectionInfo) error {

	// Close already open connections
	if err := client.Close(); err != nil {
		return err
	}

	// Create a new gRPC client
	grpcClient := &grpcservice.GRPCClient{}

	// Connect that to the service
	if err := grpcClient.Connect(info); err != nil {
		return err
	}

	// Create a new client
	microServiceClient := protobuf.NewMicroServiceClient(grpcClient.GetConnection())
	if microServiceClient == nil {
		if err := grpcClient.Close(); err != nil {
			return err
		}
		return fmt.Errorf("Setup new MicroService Client")
	}

	// Set the internal variables
	client.GrpcClient = grpcClient
	client.MicroServiceClient = microServiceClient

	return nil
}

// Close shuts down the MicroService Client by closing its gRPC Server connection
func (client *GrpcClientMicroService) Close() error {
	client.MicroServiceClient = nil
	if client.GrpcClient != nil {
		return client.GrpcClient.Close()
	}
	return nil
}
