package user

import (
	"errors"
	"log"

	middle "github.com/grpc-ecosystem/go-grpc-middleware"
	recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"go.elastic.co/apm/module/apmgrpc"
	"google.golang.org/grpc"
)

//GetElasticAPMServerOptions returns grpc server option with validator and recovery
func GetElasticAPMServerOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.ServerOption {
	serverOptions := []grpc.ServerOption{
		middle.WithUnaryServerChain(
			getDefaultUnaryOptions(interceptors...)...,
		),
		middle.WithStreamServerChain(
			getDefaultStreamOption()...,
		)}
	return serverOptions
}

//GetElasticAPMClientOption ...
func GetElasticAPMClientOption() []grpc.DialOption {
	clientOptions := []grpc.DialOption{
		grpc.WithChainUnaryInterceptor(
			apmgrpc.NewUnaryClientInterceptor(),
		),
	}
	return clientOptions
}

func defatulRecoveryHandler() func(interface{}) error {
	return func(p interface{}) (err error) {
		log.Printf("Error creating handler: %s", err.Error())
		return errors.New("Internal server Error")
	}
}

func getDefaultRecoveryOptions() []recovery.Option {
	return []recovery.Option{
		recovery.WithRecoveryHandler(defatulRecoveryHandler()),
	}
}

func getDefaultStreamOption(interceptors ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	interceptors = append(interceptors,
		validator.StreamServerInterceptor(),
		recovery.StreamServerInterceptor(getDefaultRecoveryOptions()...),
	)
	return interceptors
}

func getDefaultUnaryOptions(interceptors ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	interceptors = append(interceptors,
		validator.UnaryServerInterceptor(),
		apmgrpc.NewUnaryServerInterceptor(apmgrpc.WithRecovery()),
		recovery.UnaryServerInterceptor(getDefaultRecoveryOptions()...),
	)
	return interceptors
}
