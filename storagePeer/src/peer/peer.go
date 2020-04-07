package peer

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"encoding/json"
	"google.golang.org/grpc"

	"storagePeer/src/dht"
	"fmt"
	"time"
)

////////
// Data structures
////////

type Peer struct {

	ownIP string
	ring  *dht.RingNode
	Errs  chan error
}

////////
// Local functions
////////

func NewPeer(ownIP string, maxNodes uint64, existingIP string) *Peer {

	p := Peer{ownIP: ownIP, ring: dht.NewRingNode(ownIP, maxNodes), Errs: make(chan error, 1)}

	p.start()
	p.ring.Join(existingIP)

	return &p
}

func (p *Peer) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct{
			OwnIP string
			Ring  *dht.RingNode
		}{
			OwnIP: p.ownIP,
			Ring:  p.ring,
		})
}

// Start starts gRPC server for peer in a seperate go routine
func (p *Peer) start() {
	// Configure listening

	lis, err := net.Listen("tcp", p.ownIP)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach services to handler object
	RegisterPeerServiceServer(grpcServer, p)
	dht.RegisterRingServiceServer(grpcServer, p.ring)

	// Start listening in a separate go routine
	go func() {
		p.Errs <- grpcServer.Serve(lis)
		close(p.Errs)
	}()

	// Join the network. Build finger table and adapt the other ones.

	// Download from predecessor files that are now yours.

}

////////
// Send and recieve files
////////

// Send file to a specific node
func (p *Peer) SendFile(id uint64) (bool, error) {

	ip, err := p.ring.FindSuccessor(id)

	for {
		ip, err = p.ring.FindSuccessor(id)

		if err == nil {
			break
		} else {
			fmt.Println(err.Error())
			fmt.Println("Couldn't fetch ip")
			time.Sleep(time.Second * 1)
		}
	}

	fmt.Printf("Ring has answered with ip %s\n", ip)

	conn, err := grpc.Dial(ip, grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  cl := NewPeerServiceClient(conn)

	mes, err := cl.Ping(
		context.Background(),
		&PingMessage{Ok: true},
	)

	return mes.GetOk(), err
}

////////
// Remote calls
////////


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
