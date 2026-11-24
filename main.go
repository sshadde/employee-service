package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"example.com/employee-service/grpctransport"
	"example.com/employee-service/httptransport"
	"example.com/employee-service/implementation"
	"example.com/employee-service/pb"
	"example.com/employee-service/repository"
	"example.com/employee-service/transport"
)

func main() {
	httpAddr := flag.String("http.addr", ":8080", "HTTP listen address")
	grpcAddr := flag.String("grpc.addr", ":9090", "gRPC listen address")
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.NewSyncLogger(logger)
		logger = level.NewFilter(logger, level.AllowDebug())
		logger = log.With(logger,
			"svc", "employee",
			"ts", log.DefaultTimestampUTC,
			"caller", log.DefaultCaller,
		)
	}

	level.Info(logger).Log("msg", "service started")
	defer level.Info(logger).Log("msg", "service ended")

	repo := repository.NewMemoryRepository()
	svc := implementation.NewService(repo, logger)
	endpoints := transport.MakeEndpoints(svc)

	errs := make(chan error)

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	go func() {
		level.Info(logger).Log("transport", "HTTP", "addr", *httpAddr)
		handler := httptransport.NewService(endpoints, logger)
		server := &http.Server{
			Addr:    *httpAddr,
			Handler: handler,
		}
		errs <- server.ListenAndServe()
	}()

	go func() {
		level.Info(logger).Log("transport", "gRPC", "addr", *grpcAddr)
		listener, err := net.Listen("tcp", *grpcAddr)
		if err != nil {
			errs <- err
			return
		}
		grpcSrv := grpc.NewServer()
		pb.RegisterEmployeeServiceServer(grpcSrv, grpctransport.NewServer(endpoints))
		reflection.Register(grpcSrv)
		errs <- grpcSrv.Serve(listener)
	}()

	level.Error(logger).Log("exit", <-errs)
}
