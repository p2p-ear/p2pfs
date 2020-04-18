package main

import "C"
import (
	"fmt"
	"storagePeer/src/peer"
)

// UploadFile writes to the closest node to id on the same ring as ringIP
//export UploadFile
func UploadFile(ringIP string, id uint64, fname string, fcontent []byte) {
	if err := peer.UploadFile(ringIP, id, fname, fcontent); err != nil {
		fmt.Println("Error uploading file!", err)
	}
}

func main() {
}
