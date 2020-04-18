// This is a code for a service node, that stores files
// on itself.
package main

import (
	"flag"
	"fmt"
	"time"

	"storagePeer/src/peer"

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
	entry := flag.String("entry", "", "Ip of some existing node (if not set this node is considered first).")
	flag.Parse()

	if *ipPtr == "" {
		panic("ip flag not set")
	}
	if *numPtr == 0 {
		panic("num flag not set")
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
