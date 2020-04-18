package peer

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"fmt"
	"storagePeer/src/dht"
	"time"
)

////////
// Data structures
////////

// Peer is the peer struct
type Peer struct {
	ownIP string
	ring  *dht.RingNode
	Errs  chan error
}

////////
// Local functions
////////

// NewPeer creates new peer
func NewPeer(ownIP string, maxNodes uint64, existingIP string) *Peer {

	p := Peer{ownIP: ownIP, ring: dht.NewRingNode(ownIP, maxNodes), Errs: make(chan error, 1)}

	p.start()
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

// SendFile sends file to a specific node
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

// Read reads the content of a specified file
func (p *Peer) Read(r *ReadRequest, stream PeerService_ReadServer) error {

	f, err := os.Open(r.Name)
	defer f.Close()

	if err != nil {
		log.Fatal(err)

		return err
	}

	reader := bufio.NewReader(f)
	b := make([]byte, r.ChunkSize)

	for {
		n, readErr := reader.Read(b)

		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			return err
		}

		if err := stream.Send(&ReadReply{Data: b, Size: int64(n)}); err != nil {
			return err
		}
	}

	return nil
}

// Write writes the content of the request r onto the disk
func (p *Peer) Write(stream PeerService_WriteServer) error {

	writeInfo, err := stream.Recv()

	if err != nil {
		return err
	}

	f, err := os.Create(writeInfo.Name)
	defer f.Close()

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)

	n, err := writer.Write(writeInfo.Data)

	if err != nil {
		return err
	}

	written := int64(n)

	for {
		toWrite, readErr := stream.Recv()

		if readErr == io.EOF {
			if err = writer.Flush(); err != nil {
				return err
			}

			return stream.SendAndClose(&WriteReply{Written: int64(written)})
		}

		if readErr != nil {
			return readErr
		}

		n, err := writer.Write(toWrite.Data)

		if err != nil {
			return err
		}

		written += int64(n)
	}
}
