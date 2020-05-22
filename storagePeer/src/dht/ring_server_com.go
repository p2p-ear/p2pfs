package dht

import (
  "net/http"
  "fmt"
  "bytes"
)

const SERVER_IP = "http://172.104.136.183"

func (n *RingNode) notifyAboutDeath(deadIP string) {

  deleteUrl := fmt.Sprintf("%s/auth/node/delete", SERVER_IP)

  var jsonStr = []byte(fmt.Sprintf(`{"ip_address":"%s"}`, deadIP))
  req, err := http.NewRequest("DELETE", deleteUrl, bytes.NewBuffer(jsonStr))
  req.Header.Set("Content-Type", "application/json")

  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      panic(err)
  }

  resp.Body.Close()
}
