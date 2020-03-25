// This is a code for a service node, that stores files
// on itself.
package main

import (
	"fmt"
	"flag"
	//"log"
	//"math/rand"
	//"os"
	//"strconv"
	"time"

	"storagePeer/src/peer"

	//"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func сloseConnections(cons []*grpc.ClientConn) {
	for _, c := range cons {
		c.Close()
	}
}

func main() {

	ipPtr := flag.String("ip", "", "External IP of created node")
	numPtr := flag.Uint64("num", 0, "Number of nodes in the network")
	entry := flag.String("entry", "", "Ip of some existing node")
	flag.Parse()

	if *ipPtr == "" {
		panic("ip flag not set")
	}
	if *numPtr == 0 {
		panic("num flag not set")
	}
	if *entry == "" {
		panic("entry flag not set")
	}


	p := peer.NewPeer(*ipPtr, *numPtr, *entry)

	for {

		time.Sleep(1 * time.Second)

		mes, err := p.SendFile(0)

		if err != nil {
			fmt.Println("Couldn't send a file")
		}

		fmt.Printf("Value %t\n", mes)
	}

}

/*
// Ping testing
func testPing() {

	///////////
	// Initialization
	///////////

	// This is a setup for a testing version

	arguments := os.Args

	if len(arguments) == 1 {
		fmt.Println("Please provide the id of a node")
		return
	}

	ports := make([]string, 3)
	var ownPort string

	switch arguments[1] {
	case "0":
		ownPort = "127.0.0.1:9000"
		ports[0] = "127.0.0.1:9001"
		ports[1] = "127.0.0.1:9002"
		ports[2] = "127.0.0.1:9003"
	case "1":
		ownPort = "127.0.0.1:9001"
		ports[0] = "127.0.0.1:9000"
		ports[1] = "127.0.0.1:9002"
		ports[2] = "127.0.0.1:9003"
	case "2":
		ownPort = "127.0.0.1:9002"
		ports[0] = "127.0.0.1:9001"
		ports[1] = "127.0.0.1:9000"
		ports[2] = "127.0.0.1:9003"
	case "3":
		ownPort = "127.0.0.1:9003"
		ports[0] = "127.0.0.1:9001"
		ports[1] = "127.0.0.1:9002"
		ports[2] = "127.0.0.1:9000"
	}

	// Create a peer instance

	intID, err := strconv.Atoi(arguments[1])
	if err != nil {
		fmt.Println("Provided arguments are wrong")
		return
	}

	s := peer.Peer{
		ID:         intID,
		OwnIP:      ownPort,
		Neighbours: ports,
	}

	fmt.Println("Sending to:", ports[0])
	///////////
	// Routing
	///////////

	// Start gRPC server
	errs := s.Start()
	time.Sleep(1 * time.Second)

	/////// Configure queries

	// Get connections and clients

	connections := make([]*grpc.ClientConn, len(s.Neighbours))
	clients := make([]peer.PeerServiceClient, len(s.Neighbours))

	for i, nei := range s.Neighbours {
		connections[i], err = grpc.Dial(nei, grpc.WithInsecure())
		if err != nil {
			fmt.Println("Connection refused")
			return
		}
		clients[i] = peer.NewPeerServiceClient(connections[i])
	}
	defer сloseConnections(connections)

	/////////////////
	// Main loop
	/////////////////

	fmt.Println("Work begins")

	for {
		select {

		case err := <-errs:
			fmt.Println("Cannot serve a request", err)
			return

		default:

			_, err := clients[0].Ping(
				context.Background(),
				&peer.PingMessage{Ok: true},
			)

			if err != nil {
				fmt.Printf("%s is not ok!. error: %s\n", s.Neighbours[0], err)
			} else {
				fmt.Printf("%s is ok\n", s.Neighbours[0])
			}

			time.Sleep(1 * time.Second)
		}
	}

}

// Generate random string with length between lo and hi
func randString(lo, hi int) []byte {
	length := rand.Int()%(hi-lo) + lo
	randString := make([]byte, length)
	for i := range randString {
		randString[i] = byte(rand.Int() % 256)
	}

	return randString
}

// testFiles tests read/write capabilities of a peer
func testFiles() {
	p := peer.Peer{
		ID:         0,
		OwnIP:      "127.0.0.1:9000",
		Neighbours: make([]string, 0), //Don't need them in this test
	}

	p.Start()

	connection, err := grpc.Dial(p.OwnIP, grpc.WithInsecure())
	if err != nil {
		log.Println("Cannot open connection", err)
		return
	}

	client := peer.NewPeerServiceClient(connection)
	defer connection.Close()

	fName := "test_file"
	fContent := randString(10, 256)

	_, err = client.Write(context.Background(), &peer.WriteRequest{Name: fName, Data: fContent})
	if err != nil {
		log.Println("Write request failed:", err)
		return
	}

	read, err := client.Read(context.Background(), &peer.ReadRequest{Name: fName})
	if err != nil {
		log.Println("Read request failed", err)
		return
	}

	if r, w := len(read.Data), len(fContent); r != w {
		fmt.Println("Read", r, "Wrote", w, " - doesn't match!")
		return
	}

	for i, b := range read.Data {
		if b != fContent[i] {
			fmt.Println("Read data different from written data")
			return
		}
	}

	os.Remove(fName)
	fmt.Println("R/W test successful!")
}
*/
