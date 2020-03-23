package peer

import (
	"context"
	"io/ioutil"
	"log"
	"net"

	"google.golang.org/grpc"
)

// Peer state
type Peer struct {
	ID         int
	OwnIP      string
	Neighbours []string
}

// Start starts gRPC server for peer in a seperate go routine
func (p *Peer) Start() chan error {
	////// Configure listening

	lis, err := net.Listen("tcp", p.OwnIP)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach services to handler object
	RegisterPeerServiceServer(grpcServer, p)

	// Initialize error channel
	errs := make(chan error, 1)

	// Start listening in a separate go routine
	go func() {
		errs <- grpcServer.Serve(lis)
		close(errs)
	}()

	return errs
}

// gRPC handlers

// Ping generates response to a Ping request
func (p *Peer) Ping(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %t", in.Ok)
	return &PingMessage{Ok: true}, nil
}

func (p *Peer) Read(ctx context.Context, r *ReadRequest) (*ReadReply, error) {
	data, err := ioutil.ReadFile(r.Name)
	if err != nil {
		log.Fatal(err)

		return &ReadReply{}, err
	}

	return &ReadReply{Data: data}, nil
}

func (p *Peer) Write(ctx context.Context, r *WriteRequest) (*Empty, error) {
	return &Empty{}, ioutil.WriteFile(r.Name, r.Data, 0644)
}
