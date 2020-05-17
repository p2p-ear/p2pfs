// This is a code for a service node, that stores files
// on itself, or for a user who sends files
package main

import (
	"flag"
	"fmt"
	"storagePeer/src/peer"
	"time"
)

func main() {

	ipPtr := flag.String("ip", "", "External IP of created node")
	listenPtr := flag.String("list", "", "Local ip that we will listen to")
	numPtr := flag.Uint64("num", 0, "Number of nodes in the network")
	deltaT := flag.Int("refreshTime", 0, "Time in which fix routine is invoked (in seconds)")
	entry := flag.String("entry", "", "Ip of some existing node (if not set this node is considered first)")

	flag.Parse()

	if *ipPtr == "" {
		panic("ip flag not set")
	}
	if *numPtr == 0 {
		panic("num flag not set")
	}

	p := peer.NewPeer(*ipPtr, *listenPtr, *numPtr, *entry, time.Duration(*deltaT) * time.Second)

	err := <-p.Errs
	fmt.Println("Error!:", err)
}
