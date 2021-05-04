package user

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"
)

type service struct {
	repostory Repository
	logger    log.Logger
}

func NewService(rep Repository) Service {
	return &service{
		repostory: rep,
	}
}

type Service interface {
	CreateHelloWorld(context.Context, HelloWorld) (HelloWorld, error)
}

func (rs *service) CreateHelloWorld(ctx context.Context, request HelloWorld) (HelloWorld, error) {
	if request.Hello != "hello" {
		return HelloWorld{}, errors.New("Request not hello, it must be hello, you say " + request.Hello)
	}
	return HelloWorld{
		Hello: "Hello From Server!",
	}, nil
}
