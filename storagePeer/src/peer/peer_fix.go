package peer

import (
  "fmt"
  "bytes"
  "net/http"
  "io/ioutil"
  "encoding/json"
)

func (p *Peer) fixRoutine() error {
	const MAX_FILESIZE = 4096 * 1600 / 8
	filecont := make([]byte, MAX_FILESIZE)
	selfIP, ringsz := p.ring.RingInfo()

	for {
		newFile := <-p.ring.NewFilesChannel
		readCertificate := "miha.day.sertificat.pochitat"
		empty, err := DownloadFileRSC(selfIP, newFile, ringsz, filecont, readCertificate)
		if err != nil {
			return err
		}

		writeCertificate := "miha.day.sertificat.popisat"
		err = UploadFileRSC(selfIP, newFile, ringsz, filecont[:MAX_FILESIZE-empty], writeCertificate)
		if err != nil {
			return err
		}
	}
}

const SERVER_IP = "http://172.104.136.183"

func (p *Peer) notifyAboutArrival() {

	notifyUrl := fmt.Sprintf("%s/auth/node/add", SERVER_IP)

  request, err := json.Marshal(map[string]string{
    "ip_address": p.ownIP,
  })

  if err != nil {
    panic(err)
  }

  resp, err := http.Post(notifyUrl, "application/json", bytes.NewBuffer(request))

  body, _ := ioutil.ReadAll(resp.Body)
  fmt.Println("response Body:", string(body))

  resp.Body.Close()
}
