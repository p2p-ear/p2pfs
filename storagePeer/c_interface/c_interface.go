package main

import (
	"C"
	"fmt"
	"storagePeer/src/peer"
)

// UploadFile writes to the closest node to hash of fname on the same ring as ringIP
//export UploadFile
func UploadFile(ringIP string, fname string, ringsz uint64, fcontent []byte) {
	if err := peer.UploadFile(ringIP, fname, ringsz, fcontent); err != nil {
		fmt.Println("Error uploading file!", err)
	}
}

// DownloadFile downloads the file from the closest node to hash of fname on the same ring as ringIP
//export DownloadFile
func DownloadFile(ringIP string, fname string, ringsz uint64) []byte {
	fcontent, err := peer.DownloadFile(ringIP, fname, ringsz)
	if err != nil {
		fmt.Println("Error downloading file!", err)
	}

	return fcontent
}

func main() {
}
