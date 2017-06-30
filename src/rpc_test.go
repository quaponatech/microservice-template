package microservice

import (
	"context"
	"testing"
	"time"

	"github.com/quaponatech/golang-extensions/test"
	"github.com/quaponatech/microservice-template/src/protobuf"
)

func getService(t *testing.T) MicroService {
	service := getDefaultService(t)
	go func() {
		err := service.Serve()
		test.AssertThat(t, err, nil)
	}()
	time.Sleep(time.Second)
	return service
}

func TestSuccessHello(t *testing.T) {
	t.Run("Succeeds", func(t *testing.T) {
		// Start a new service
		service := getService(t)
		testString := "Hello World!"

		// Test the RPC
		msg := &protobuf.Request{Message: testString}
		response, err := service.Hello(context.Background(), msg)
		test.AssertThat(t, err, nil)
		test.AssertThat(t, response.GetMessage(), testString)
	})
}
