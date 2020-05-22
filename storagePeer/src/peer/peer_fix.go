package peer

import (
	"bytes"
	"fmt"
	"net/http"
)

func (p *Peer) fixRoutine() error {
	const MAX_FILESIZE = 4096 * 1600 / 8
	filecont := make([]byte, MAX_FILESIZE)

	for {
		newFile := <-p.ring.NewFilesChannel
		readCertificate := "miha.day.sertificat.pochitat"
		empty, err := DownloadFileRSC(p.ring.self.IP, newFile, p.ring.maxnodes, filecont, certificate)
		if err != nil {
			return err
		}

		writeCertificate := "miha.day.sertificat.popisat"
		err = UploadFileRSC(p.ring.self.IP, newFile, p.ring.maxnodes, filecont[:MAX_FILESIZE-empty], certificate)
		if err != nil {
			return err
		}
	}
}

const SERVER_IP = "http://172.104.136.183"

func (p *Peer) notifyAboutArrival() {

	notifyUrl := fmt.Sprintf("%s/auth/node/add", SERVER_IP)

	var jsonStr = []byte(fmt.Sprintf(`{"ip_address":"%s"}`, p.ownIP))
	req, err := http.NewRequest("POST", notifyUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()
}
