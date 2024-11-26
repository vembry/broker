package grpc

import (
	"broker/pb"
	"broker/server"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type grpcserver struct {
	address string
	server  *grpc.Server
}

func New(
	broker server.IBroker,
	address string,
	serverOptions ...grpc.ServerOption,
) *grpcserver {
	// create grpc server
	grpcServer := grpc.NewServer(serverOptions...) // register interceptors

	// construct grpc handler
	handler := NewHandler(broker)

	// register handler onto the grpc server
	pb.RegisterBrokerServer(grpcServer, handler)

	return &grpcserver{
		address: address,
		server:  grpcServer,
	}
}

// Start starts grpc server
func (g *grpcserver) Start() error {
	log.Printf("running grpc server at %s", g.address)

	reflection.Register(g.server)

	lis, err := net.Listen("tcp", g.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return g.server.Serve(lis)
}

// Stop stops grpc server
func (g *grpcserver) Stop() {
	g.server.GracefulStop()
}
