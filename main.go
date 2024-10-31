package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"broker/broker"
	grpcserver "broker/server/grpc"
)

func main() {
	log.Printf("hello broker!")

	// setup app's tracer
	shutdownHandler := NewTracer()
	defer shutdownHandler()

	queue := broker.New(nil) // initiate core queue
	queue.Start()            // restore backed-up queues

	grpcServer := grpcserver.New(queue) // iniitate grpc server

	// grpc server
	go func() {
		if err := grpcServer.Start(); err != nil {
			log.Fatalf("found error on starting grpc server. err=%v", err)
		}
	}()

	WatchForExitSignal()
	log.Println("shutting down...")

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
