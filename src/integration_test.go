package microservice

import (
	"context"
	"fmt"
	"testing"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

// Returns the current cluster connection information
func getClusterConnectionInfo() *grpcservice.ConnectionInfo {
	return &grpcservice.ConnectionInfo{
		UseTLS:              false,
		CertFile:            "",
		ServerHostName:      "",
		IP:                  clusterIP,
		Port:                fmt.Sprint(clusterPort),
		TimeoutInMilliSecs:  10000,
		RetryTimes:          1,
		RetryAfterMilliSecs: 100,
	}
}

// Simple integration test which tries to connect to a real client
func TestSuccessIntegration(t *testing.T) {
	// Do not exeute this test in short mode
	if testing.Short() {
		return
	}
	t.Run("Succeeds", func(t *testing.T) {
		// Create a connection info
		info := getClusterConnectionInfo()
		test.AssertThat(t, info, nil, "not")
		t.Logf("Connecting to cluster '%s' via port '%s'", info.IP, info.Port)

		// Connect to the cluster
		c := GrpcClientMicroService{}
		err := c.Connect(info)
		test.AssertThat(t, err, nil)

		// Test the RPC
		testString := "Hello World!"
		msg := &protobuf.Request{Message: testString}
		response, err := c.MicroServiceClient.Hello(context.Background(), msg)
		test.AssertThat(t, err, nil)
		test.AssertThat(t, response.GetMessage(), testString)
	})
}
