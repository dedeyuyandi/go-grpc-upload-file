package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	v1_user_grpc "github.com/dedeyuyandi/go-grpc-upload-file/proto"
	"github.com/go-kit/kit/auth/jwt"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"

	_ "github.com/lib/pq"
)

//GrpcServer ....
type GrpcServer struct {
	createHelloWorld grpctransport.Handler
}

func GRPCServerRun(addr string, userSvr v1_user_grpc.UserServer) {
	// error info
	errs := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGALRM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	fmt.Println(addr, "addr")
	grpcListener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var opts []grpc.ServerOption
	opts = GetElasticAPMServerOptions()

	go func() {
		baseServer := grpc.NewServer(opts...)
		v1_user_grpc.RegisterUserServer(baseServer, userSvr)
		fmt.Println("🚀 Server recruitment started successfully 🚀 - Running on", addr)
		baseServer.Serve(grpcListener)
	}()
	fmt.Println("exit", <-errs)
}

func NewGRPCServer() *GrpcServer {
	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = log.With(logger,
			"service", "user",
			"time:", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}
	fmt.Println(logger, "logger")

	// start database conn
	var db *sql.DB
	{
		var err error
		db, err = sql.Open("postgres", "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable")
		if err != nil {
			level.Error(logger).Log("exit", err)
			os.Exit(-1)
		}
	}

	repository := NewRepo(db)
	srv := NewService(repository)
	endpoints := MakeEndpoints(srv)

	handlerOpt := []grpctransport.ServerOption{
		grpctransport.ServerBefore(jwt.GRPCToContext()),
	}

	return &GrpcServer{
		createHelloWorld: grpctransport.NewServer(
			endpoints.CreateHelloWorld,
			decodeCreateHelloWorld,
			encodeCreateHelloWorld,
			handlerOpt...,
		),
	}
}

func (g *GrpcServer) CreateHelloWorld(ctx context.Context, request *v1_user_grpc.HelloWorld) (*v1_user_grpc.HelloWorld, error) {
	_, resp, err := g.createHelloWorld.ServeGRPC(ctx, request)
	if err != nil {
		return nil, err
	}
	return resp.(*v1_user_grpc.HelloWorld), nil
}

func decodeCreateHelloWorld(ctx context.Context, request interface{}) (interface{}, error) {
	req, ok := request.(*v1_user_grpc.HelloWorld)
	if !ok {
		return nil, errors.New("Invalid type")
	}
	return &HelloWorld{
		Hello: req.Hello,
	}, nil
}

func encodeCreateHelloWorld(ctx context.Context, response interface{}) (interface{}, error) {
	resp, ok := response.(HelloWorld)
	if !ok {
		return nil, errors.New("Invalid type")
	}
	return &v1_user_grpc.HelloWorld{
		Hello: resp.Hello,
	}, nil
}
