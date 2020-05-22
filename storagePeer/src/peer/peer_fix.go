package peer

import (
  "fmt"
  "bytes"
  "net/http"
  "io/ioutil"
  "encoding/json"
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
