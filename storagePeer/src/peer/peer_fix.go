package peer

import (
  "fmt"
  "bytes"
  "net/http"
)

func (p *Peer) fixRoutine() {

  for {
    newFile := <- p.ring.NewFilesChannel
    fmt.Printf("Have to download %s", newFile)
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
