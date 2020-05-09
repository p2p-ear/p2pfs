package peer

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"google.golang.org/grpc"
)

func genIP() func() string {
	port := 9000
	return func() string {
		ip := fmt.Sprintf("127.0.0.1:%d", port)
		port++
		return ip
	}
}

var IP = genIP()

////////
// Service funcs
////////

// Generate random string with specified length
func randString(length int) []byte {

	randString := make([]byte, length)
	for i := range randString {
		randString[i] = byte(rand.Intn(256))
	}

	return randString
}

// Make one peer
func makePeer() (string, uint64, PeerServiceClient, *grpc.ClientConn, error) {
	ownIP := IP()

	ringsz := uint64(1000)
	NewPeer(ownIP, ringsz, "")

	connection, err := grpc.Dial(ownIP, grpc.WithInsecure())
	if err != nil {
		return "", 0, nil, nil, err
	}

	client := NewPeerServiceClient(connection)

	return ownIP, ringsz, client, connection, nil
}

// Make n peers in one ring
func makeRing(n uint) (string, uint64) {

	ringsz := uint64(1000)
	host := IP()

	NewPeer(host, ringsz, "")

	ips := make([]string, n)
	for i := uint(0); i < n; i++ {
		ips[i] = IP()
		NewPeer(ips[i], ringsz, host)
	}

	return host, ringsz
}

// TestRW tests read/write capabilities of a peer
func TestRW(t *testing.T) {

	_, _, client, connection, err := makePeer()
	defer connection.Close()
	if err != nil {
		t.Error(err)
	}

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

func TestUpload(t *testing.T) {

	ownIP, ringsz, _, connection, err := makePeer()
	defer connection.Close()
	if err != nil {
		t.Error(err)
	}

	fcontent := randString(4096)
	fname := "testfile.txt"

	err = UploadFile(ownIP, fname, ringsz, fcontent)
	if err != nil {
		t.Error("Unable to send file", err)
	}

	fcontentRead, err := ioutil.ReadFile(fname)
	if err != nil {
		t.Error("Unable to read sent file", err)
	}

	if len(fcontentRead) != len(fcontent) {
		t.Error("Lengths don't match: written", len(fcontent), "read", len(fcontentRead))
	}

	for i, b := range fcontentRead {
		if fcontent[i] != b {
			t.Error("Content doesn't match!")
		}
	}

	os.Remove(fname)
}

func TestDownload(t *testing.T) {

	ownIP, ringsz, _, connection, err := makePeer()
	defer connection.Close()
	if err != nil {
		t.Error(err)
	}

	fcontent := randString(4096)
	fname := "testfile.txt"

	ioutil.WriteFile(fname, fcontent, 0644)

	fcontentRead := make([]byte, len(fcontent))
	empty, err := DownloadFile(ownIP, fname, ringsz, fcontentRead)
	if err != nil {
		t.Error("Unable to download file", err)
	}

	if empty != 0 {
		t.Error("Lengths don't match: empty =", empty)
	}

	for i, b := range fcontentRead {
		if fcontent[i] != b {
			t.Error("Content doesn't match!")
		}
	}

	os.Remove(fname)
}

func findBin() (string, error) {
	absPath, err := filepath.Abs("./")
	if err != nil {
		return "", fmt.Errorf("Abs error: %s", err.Error())
	}

	projName := "storagePeer"
	projPathInd := strings.LastIndex(absPath, projName)
	projPath := absPath[:projPathInd+len(projName)]

	return projPath + "/bin", nil
}

func TestC(t *testing.T) {
	ip, ringsz, _, conn, err := makePeer()
	if err != nil {
		t.Error("Unable to create peer", err)
	}
	defer conn.Close()

	binpath, err := findBin()
	if err != nil {
		t.Error("Unable to find file:", err)
	}

	var errStream bytes.Buffer
	cTest := exec.Cmd{Path: "./c_test", Dir: binpath, Args: []string{binpath + "/c_test", ip, strconv.Itoa(int(ringsz))}, Stderr: &errStream}
	err = cTest.Run()
	if err != nil {
		t.Error("Run error:", err, "stderr:", errStream.String())
	}
}

func TestRSC(t *testing.T) {
	host, ringsz := makeRing(10)

	fname := "testfile"
	fcontent := randString(4096)

	err := UploadFileRSC(host, fname, ringsz, fcontent)
	if err != nil {
		t.Error("UploadRSC error:", err)
	}

	fcontentRead := make([]byte, len(fcontent)*2)
	empty, err := DownloadFileRSC(host, fname, ringsz, fcontentRead)
	if err != nil {
		t.Error("DownloadRSC error:", err)
	}

	if empty < 0 {
		t.Error("File read too large, empty =", empty)
	}

	for i, b := range fcontent {
		if b != fcontentRead[i] {
			t.Error("Bytes at place", i, "don't match")
		}
	}
}
