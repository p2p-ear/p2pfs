package main

import (
	"C"
	"fmt"
	"storagePeer/src/peer"
)

// UploadFile writes to the closest node to id on the same ring as ringIP
//export UploadFile
func UploadFile(ringIP string, fname string, ringsz uint64, fcontent []byte) {
	if err := peer.UploadFile(ringIP, fname, ringsz, fcontent); err != nil {
		fmt.Println("Error uploading file!", err)
	}
}

func main() {
}
