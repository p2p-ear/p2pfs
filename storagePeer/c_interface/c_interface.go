package main

import (
	"C"
	"log"
	"storagePeer/src/peer"
)

// UploadFile writes to the closest node to hash of fname on the same ring as ringIP
//export UploadFile
func UploadFile(ringIP string, fname string, ringsz uint64, fcontent []byte) {
	if err := peer.UploadFile(ringIP, fname, ringsz, fcontent); err != nil {
		log.Println("Error uploading file!", err)
	}
}

// DownloadFile downloads the file from the closest node to hash of fname on the same ring as ringIP
//export DownloadFile
func DownloadFile(ringIP string, fname string, ringsz uint64, fcontent []byte) int {
	emptySpace, err := peer.DownloadFile(ringIP, fname, ringsz, fcontent)
	if err != nil {
		log.Println("Error downloading file!", err)
	}

	return emptySpace
}

//export DeleteFile
func DeleteFile(ringIP string, fname string, ringsz uint64) {
	err := peer.DeleteFile(ringIP, fname, ringsz)
	if err != nil {
		log.Println("Error deleting file!", err)
	}
}

func main() {
}
