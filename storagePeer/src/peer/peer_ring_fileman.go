// Ring-targeted filemanagemenet
package peer

import (
	"context"
	"fmt"
	"storagePeer/src/dht"
	"time"
)

// Service func to connect to the ringIP, and find successor on that ring
func findSuccessorWithRingIP(ringIP string, id uint64) (string, error) {

	someConn, somePeer, err := Connect(ringIP)
	if err != nil {
		return "", err
	}
	defer someConn.Close()

	ip := ""
	for {
		succReply, err := somePeer.FindSuccessorInRing(context.Background(), &FindSuccRequest{Id: id})

		if err == nil {
			ip = succReply.Ip
			break
		} else {
			fmt.Println(err.Error())
			fmt.Println("Couldn't fetch ip")
			time.Sleep(time.Second * 1)
		}
	}

	fmt.Printf("Ring has answered with ip %s\n", ip)
	return ip, nil
}

// uploadFile uploads file to the successor of an id. ringIP - ip of someone on the ring
func uploadFile(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) error {

	id := dht.Hash([]byte(fname), ringsz)
	targetIP, err := findSuccessorWithRingIP(ringIP, id)
	if err != nil {
		return err
	}

	return sendFile(targetIP, fname, fcontent, certificate)
}

// downloadFile downloads file from the corresponding node
func downloadFile(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) (int, error) {
	id := dht.Hash([]byte(fname), ringsz)
	targetIP, err := findSuccessorWithRingIP(ringIP, id)

	if err != nil {
		return 0, err
	}

	return recvFile(targetIP, fname, fcontent, certificate)
}

func deleteFile(ringIP string, fname string, ringsz uint64, certificate string) error {
	id := dht.Hash([]byte(fname), ringsz)
	targetIP, err := findSuccessorWithRingIP(ringIP, id)

	if err != nil {
		return err
	}

	return remvFile(targetIP, fname, certificate)
}
