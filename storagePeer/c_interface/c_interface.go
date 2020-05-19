package main

import (
	"C"
	"log"
	"storagePeer/src/peer"
)

//export UploadFileRSC
func UploadFileRSC(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) {

	if err := peer.UploadFileRSC(ringIP, fname, ringsz, fcontent, certificate); err != nil {
		log.Println("Error uploading file (RSC)!", err)
	}
}

//export DownloadFileRSC
func DownloadFileRSC(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) int {

	emptySpace, err := peer.DownloadFileRSC(ringIP, fname, ringsz, fcontent, certificate)
	if err != nil {
		log.Println("Error downloading file (RSC)!", err)
	}

	return emptySpace
}

//export DeleteFileRSC
func DeleteFileRSC(ringIP string, fname string, ringsz uint64, certificate string) {

	if err := peer.DeleteFileRSC(ringIP, fname, ringsz, certificate); err != nil {
		log.Println("Error deleting file (RSC!", err)
	}
}

func main() {
}
