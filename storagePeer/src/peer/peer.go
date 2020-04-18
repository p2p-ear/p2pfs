package peer

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"log"
	"math"
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
	RegisterPeerServiceServer(grpcServer, p)
	dht.RegisterRingServiceServer(grpcServer, p.ring)

	// Start listening in a separate go routine
	go func() {
		p.Errs <- grpcServer.Serve(lis)
		close(p.Errs)
	}()

	// Download from predecessor files that are now yours.

}

////////
// Send and recieve files
////////

// Кто за лопату ответственный
func findSuccessor(ringIP string, id uint64) (string, error) {

	someConn, err := grpc.Dial(ringIP, grpc.WithInsecure())
	if err != nil {
		return "", err
	}
	defer someConn.Close()

	somePeer := NewPeerServiceClient(someConn)
	ip := ""

	for {
		succReply, err := somePeer.FindSuccessor(context.Background(), &FindSuccRequest{Id: id})

		if err == nil {
			ip = succReply.Ip
			break
		} else {
			fmt.Println(err.Error())
			fmt.Println("Couldn't fetch ip")
			time.Sleep(time.Second * 1)
		}
	}

	return ip, nil
}

// UploadFile uploads file to the successor of an id. ringIP - ip of someone on the ring
//export UploadFile
func UploadFile(ringIP string, id uint64, fname string, fcontent []byte) error {

	ip, err := findSuccessor(ringIP, id)
	fmt.Printf("Ring has answered with ip %s\n", ip)

	conn, err := grpc.Dial(ringIP, grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	cl := NewPeerServiceClient(conn)

	fmt.Println("Opening write stream...")
	// Stream to write
	wstream, err := cl.Write(context.Background())
	for err != nil {
		fmt.Println(err.Error())
		fmt.Println("Couldn't initialize remote write stream")
		time.Sleep(time.Second * 1)

		wstream, err = cl.Write(context.Background())
	}

	fmt.Println("Sending filename...")
	// Send filename
	err = wstream.Send(&WriteRequest{Name: fname})

	for err != nil {
		fmt.Println(err.Error())
		fmt.Println("Couldn't send filename")
		time.Sleep(time.Second * 1)
		err = wstream.Send(&WriteRequest{Name: fname})
	}

	chunkSize := 8
	chunkAmnt := int(math.Ceil(float64(len(fcontent)) / float64(chunkSize)))

	fmt.Println("Writing to file, total chunks:", chunkAmnt)
	for i := 0; i < chunkAmnt; i++ {

		curChunk := fcontent[i*chunkSize:]
		fmt.Println("ChunkSize", chunkSize, "Left", len(curChunk))
		if len(curChunk) > chunkSize {
			fmt.Println()
			curChunk = curChunk[:chunkSize]
		}

		err := wstream.Send(&WriteRequest{Data: curChunk})

		for err != nil {
			fmt.Println(err.Error())
			fmt.Println("Unable to send file contents")
			time.Sleep(time.Second * 1)
			err = wstream.Send(&WriteRequest{Name: fname})
		}

		fmt.Println("Written chunk", i)
	}

	_, err = wstream.CloseAndRecv()
	fmt.Println("Finished writing, err:", err)
	return err
}

////////
// Remote calls
////////

// Ping generates response to a Ping request
func (p *Peer) Ping(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %t", in.Ok)
	return &PingMessage{Ok: true}, nil
}

// FindSuccessor finds id's successor
func (p *Peer) FindSuccessor(ctx context.Context, r *FindSuccRequest) (*FindSuccReply, error) {
	ip, err := p.ring.FindSuccessor(r.Id)
	return &FindSuccReply{Ip: ip}, err
}

// Read & Write ----------------------

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
