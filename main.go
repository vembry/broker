package main

import (
	"log"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"

	"broker/broker"
	grpcserver "broker/server/grpc"
	httpserver "broker/server/http"
)

func main() {
	log.Printf("hello broker!")

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := broker.New(nil) // initiate core queue
	queue.Start()            // restore backed-up queues

	httpServer := httpserver.New(queue) // initiate http server
	grpcServer := grpcserver.New(queue) // iniitate grpc server

	// http server
	go func() {
		if err := httpServer.Start(); err != nethttp.ErrServerClosed {
			log.Fatalf("found error on starting http server. err=%v", err)
		}
	}()

	// grpc server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("found error on starting grpc server. err=%v", err)
		}
	}()

	WatchForExitSignal()
	log.Println("shutting down...")

	httpServer.Stop()
	grpcServer.Stop()
	queue.Stop() // shutdown queue
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
