package peer

import(
  "log"

  "golang.org/x/net/context"
)

///////// Object's state
type Peer struct {
  Id int
  OwnIp string
  Neighbours []string
}

//////// gRPC handlers

// Ping generates response to a Ping request
func (p *Peer) Ping(ctx context.Context, in *PingMessage) (*PingMessage, error) {
  log.Printf("Receive message %s", in.Ok)
  return &PingMessage{Ok: true}, nil
}

//////// Request functions
