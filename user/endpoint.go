package user

import (
	"context"

	"github.com/go-kit/kit/endpoint"
)

type Endpoints struct {
	CreateHelloWorld endpoint.Endpoint
}

func MakeEndpoints(svc Service) Endpoints {
	return Endpoints{
		CreateHelloWorld: makeHelloWorld(svc),
	}
}
func makeHelloWorld(svc Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(*HelloWorld)
		return svc.CreateHelloWorld(ctx, HelloWorld{
			Hello: req.Hello,
		})
	}
}
