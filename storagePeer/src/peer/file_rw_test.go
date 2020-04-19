package peer

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"os"
	"testing"

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

// TestRW tests read/write capabilities of a peer
func TestRW(t *testing.T) {
	ownIP := "127.0.0.1:9000"
	NewPeer(ownIP, 2, "")

	connection, err := grpc.Dial(ownIP, grpc.WithInsecure())
	if err != nil {
		t.Error("Cannot open connection", err)
	}

	client := NewPeerServiceClient(connection)
	defer connection.Close()

	fName := "test_file"

	wstream, err := client.Write(context.Background())
	if err != nil {
		t.Error("Creating write stream failed:", err)
	}

	if err := wstream.Send(&WriteRequest{Name: fName}); err != nil {
		t.Error("Initializing write stream failed:", err)
	}

	chunkAmnt := rand.Intn(16) + 16
	fLength := chunkAmnt * 8
	fContent := make([]byte, 0)

	for i := 0; i < chunkAmnt; i++ {
		nextChunk := randString(8)
		fContent = append(fContent, nextChunk...)
		if err = wstream.Send(&WriteRequest{Data: nextChunk}); err != nil {
			t.Error("Error writing to stream:", err)

		}
	}

	lastChunkLen := rand.Intn(4) + 3
	lastChunk := randString(lastChunkLen)

	fLength += lastChunkLen
	fContent = append(fContent, lastChunk...)

	if err := wstream.Send(&WriteRequest{Data: lastChunk}); err != nil {
		fmt.Println("Error writing final bytes to stream:", err)
	}

	writeReply, err := wstream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error closing write stream:", err)
	}

	written := int(writeReply.Written)

	rstream, err := client.Read(context.Background(), &ReadRequest{Name: fName, ChunkSize: 8})
	if err != nil {
		t.Error("Creating read stream failed:", err)
	}

	readContent := make([]byte, 0)
	for {
		readReply, err := rstream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			t.Error("Error reading from stream:", err)
		}

		nextChunk := readReply.Data[:readReply.Size]

		readContent = append(readContent, nextChunk...)
	}

	if r, w := len(readContent), written; r != w {
		t.Error("Read", r, "Wrote", w, " - doesn't match!")
	}

	for i, b := range readContent {
		if b != fContent[i] {
			t.Error("Read data different from written data")
		}
	}

	os.Remove(fName)
}
