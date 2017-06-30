package microservice

import (
	"golang.org/x/net/context"

	"github.com/quaponatech/microservice-template/src/protobuf"
)

// Hello returns a Response to the given Response
func (service *MicroService) Hello(ctx context.Context, request *protobuf.Request) (*protobuf.Response, error) {
	// TODO: Add your own RPC logic and rename the messages
	return &protobuf.Response{Message: request.GetMessage()}, nil
}
