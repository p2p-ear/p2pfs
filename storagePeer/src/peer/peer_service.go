//Various service functions
package peer

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"storagePeer/src/dht"
	"time"

	"google.golang.org/grpc"
)

// NewPeer creates new peer
func NewPeer(ownIP string, maxNodes uint64, existingIP string) *Peer {

	const deltaT = time.Second //TODO: make it configurable

	p := Peer{ownIP: ownIP, ring: dht.NewRingNode(ownIP, maxNodes, deltaT), Errs: make(chan error, 1)}

	p.start()

	// Join the network. Build finger table and adapt the other ones.
	p.ring.Join(existingIP)

	return &p
}

// MarshalJSON converts peer to JSON
func (p *Peer) MarshalJSON() ([]byte, error) {

	return json.Marshal(struct {
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

	p.ring.Start(grpcServer)
	RegisterPeerServiceServer(grpcServer, p)

	// Start listening in a separate go routine
	go func() {
		p.Errs <- grpcServer.Serve(lis)
		close(p.Errs)
	}()
}

// Connect connects to peer with specified IP
func Connect(targetIP string) (*grpc.ClientConn, PeerServiceClient, error) {
	conn, err := grpc.Dial(targetIP, grpc.WithInsecure())

	if err != nil {
		for {
			conn, err = grpc.Dial(targetIP, grpc.WithInsecure())
			if err != nil {
				fmt.Println(err)
				fmt.Println("Couldn't connect to node", targetIP)
			} else {
				break
			}

			time.Sleep(time.Second * 1)
		}
	}

	cl := NewPeerServiceClient(conn)
	return conn, cl, nil
}
