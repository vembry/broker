package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"broker/message"
	grpcserver "broker/server/grpc"
	"broker/server/grpc/interceptor"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	log.Printf("hello broker!")

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := message.NewBroker(nil) // initiate core queue
	queue.Start()                   // restore backed-up queues

	// initate grpc server
	grpcServer := grpcserver.New(
		queue,   // dependencies
		":4000", // grpc server address

		// register interceptors
		grpc.StatsHandler(otelgrpc.NewServerHandler()),                     // for otel related
		grpc.UnaryInterceptor(interceptor.Authenticate("some-passphrase")), //basic authenticate
	)

	// run grpc server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("found error on starting grpc server. err=%v", err)
		}
	}()

	WatchForExitSignal()

	log.Println("shutting down...")

	grpcServer.Stop() // shutdown grpc server
	queue.Stop()      // shutdown queue
}

// WatchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func WatchForExitSignal() os.Signal {
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
