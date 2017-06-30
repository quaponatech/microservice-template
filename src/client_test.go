package microservice

import (
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/quaponatech/golang-extensions/grpcservice"
	"github.com/quaponatech/golang-extensions/test"
	"github.com/quaponatech/microservice-template/src/mock_client"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

// Test local server and client
func TestSuccessConnect(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		port := 8899

		// Start a server
		service := MicroService{}
		client := &GrpcClientMicroService{}
		test.AssertThat(t, service.Setup(serverName,
			grpcservice.NewGRPCServer(false, "", "", port),
			getLogger(),
			client),
			nil)

		go func() {
			err := service.Serve()
			test.AssertThat(t, err, nil)
		}()
		time.Sleep(time.Second)

		// Create a connection info
		info := &grpcservice.ConnectionInfo{
			UseTLS:              true,
			CertFile:            "",
			ServerHostName:      "",
			IP:                  "127.0.0.1",
			Port:                "8899",
			TimeoutInMilliSecs:  1000,
			RetryTimes:          1,
			RetryAfterMilliSecs: 100,
		}

		// Connect encrypted, which should not work
		c := GrpcClientMicroService{}
		err := c.Connect(info)
		test.AssertThat(t, err, nil, "not")

		// Close the inexisting connection
		err = c.Close()
		test.AssertThat(t, err, nil)

		// Connect plain, this should work
		c = GrpcClientMicroService{}
		info.UseTLS = false
		err = c.Connect(info)
		test.AssertThat(t, err, nil)

		// Close the existing connection
		err = c.Close()
		test.AssertThat(t, err, nil)

		// Stop the server
		err = service.Stop()
		test.AssertThat(t, err, nil)
	})
}

// Test mocked Client
func TestSuccessMock(t *testing.T) {
	// Setup the mocking controller
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	// Create a mocked Client which always returns the testString
	testString := "Hello"
	mockee := mock_protobuf.NewMockMicroServiceClient(mockCtrl)
	mockee.EXPECT().Hello(gomock.Any(), gomock.Any()).Return(&protobuf.Response{Message: testString}, nil)

	// Test the RPC with the mocked client
	response, err := mockee.Hello(nil, nil)
	test.AssertThat(t, err, nil)
	test.AssertThat(t, response.GetMessage(), testString)
}
