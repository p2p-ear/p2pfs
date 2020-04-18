package testfiles

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"storagePeer/src/peer"

	"google.golang.org/grpc"
)

// Generate random string with length between lo and hi
func randString(length int) []byte {

	randString := make([]byte, length)
	for i := range randString {
		randString[i] = byte(rand.Intn(256))
	}

	return randString
}

// TestFiles tests read/write capabilities of a peer
func TestFiles() {
	ownIP := "127.0.0.1:9000"
	peer.NewPeer(ownIP, 2, "")

	connection, err := grpc.Dial(ownIP, grpc.WithInsecure())
	if err != nil {
		log.Println("Cannot open connection", err)
		return
	}

	client := peer.NewPeerServiceClient(connection)
	defer connection.Close()

	fName := "test_file"

	wstream, err := client.Write(context.Background())
	if err != nil {
		log.Println("Creating write stream failed:", err)
		return
	}

	if err := wstream.Send(&peer.WriteRequest{Name: fName}); err != nil {
		log.Println("Initializing write stream failed:", err)
	}

	chunkAmnt := rand.Intn(16) + 16
	fLength := chunkAmnt * 8
	fContent := make([]byte, 0)

	for i := 0; i < chunkAmnt; i++ {
		nextChunk := randString(8)
		fContent = append(fContent, nextChunk...)
		if err = wstream.Send(&peer.WriteRequest{Data: nextChunk}); err != nil {
			log.Println("Error writing to stream:", err)
			return
		}
	}

	lastChunkLen := rand.Intn(4) + 3
	lastChunk := randString(lastChunkLen)

	fLength += lastChunkLen
	fContent = append(fContent, lastChunk...)

	if err := wstream.Send(&peer.WriteRequest{Data: lastChunk}); err != nil {
		fmt.Println("Error writing final bytes to stream:", err)
		return
	}

	writeReply, err := wstream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error closing write stream:", err)
		return
	}

	written := int(writeReply.Written)

	rstream, err := client.Read(context.Background(), &peer.ReadRequest{Name: fName, ChunkSize: 8})
	if err != nil {
		log.Println("Creating read stream failed:", err)
		return
	}

	readContent := make([]byte, 0)
	for {
		readReply, err := rstream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("Error reading from stream:", err)
		}

		nextChunk := readReply.Data[:readReply.Size]

		readContent = append(readContent, nextChunk...)
	}

	if r, w := len(readContent), written; r != w {
		fmt.Println("Read", r, "Wrote", w, " - doesn't match!")
		return
	}

	for i, b := range readContent {
		if b != fContent[i] {
			fmt.Println("Read data different from written data")
			return
		}
	}

	os.Remove(fName)
	fmt.Println("R/W test successful!")
}
