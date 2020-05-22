package peer

import (
<<<<<<< HEAD
  "fmt"
  "bytes"
  "net/http"
  "io/ioutil"
  "encoding/json"
=======
	"bytes"
	"fmt"
	"net/http"
>>>>>>> 69f9b07e550e31f3f2b7538f0a1ee007fea208e7
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

<<<<<<< HEAD
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
=======
	var jsonStr = []byte(fmt.Sprintf(`{"ip_address":"%s"}`, p.ownIP))
	req, err := http.NewRequest("POST", notifyUrl, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	resp.Body.Close()
>>>>>>> 69f9b07e550e31f3f2b7538f0a1ee007fea208e7
}
