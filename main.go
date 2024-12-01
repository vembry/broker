package main

import (
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"broker/message"
	grpcserver "broker/server/grpc"
	"broker/server/grpc/interceptor"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("starting...")

	// compile configs
	appConfig := newConfig()

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := message.NewBroker(nil) // initiate core queue
	queue.Start()                   // restore backed-up queues

	// prep grpc interceptor
	interceptors := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()), // for otel related
	}

	// when defined, include authorization interceptor
	if !isStringEmpty(appConfig.Authorization) {
		interceptors = append(interceptors, grpc.UnaryInterceptor(interceptor.Authenticate(appConfig.Authorization)))
	} else {
		log.Printf("running without authorization")
	}

	// initate grpc server
	grpcServer := grpcserver.New(
		queue,                 // dependencies
		appConfig.GrpcAddress, // grpc server address

		// register interceptors
		interceptors...,
	)

	// run grpc server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("found error on starting grpc server. err=%v", err)
		}
	}()

	watchForExitSignal()

	log.Println("shutting down...")

	grpcServer.Stop() // shutdown grpc server
	queue.Stop()      // shutdown queue
}

// WatchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func watchForExitSignal() os.Signal {
	log.Printf("awaiting sigterm...")
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}

// config contains all the configurations for the app to use
type config struct {
	Authorization string
	GrpcAddress   string
}

// newConfig compiles app's config provided on os env var
func newConfig() *config {
	cfg := &config{
		Authorization: os.Getenv("AUTHORIZATION"),
		GrpcAddress:   os.Getenv("GRPC_ADDRESS"),
	}

	if isStringEmpty(cfg.GrpcAddress) {
		// defaults address
		cfg.GrpcAddress = ":4000"
	}

	return cfg
}

func isStringEmpty(str string) bool {
	str = strings.TrimSpace(str)
	return str == ""
}
