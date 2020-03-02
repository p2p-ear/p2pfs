// This is a code for a service node, that stores files
// on itself.
package main

import (
  "fmt"
  "log"
  "net"
  "os"
  "strconv"
  "time"

  "p2pClient/src/peer"
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

func CloseConnections(cons []*grpc.ClientConn) {
  for _, c := range cons {
    c.Close()
  }
}
// main start a gRPC server and waits for connection
func main() {

  ///////////
  // Initialization
  ///////////

  // This is a setup for a testing version

  arguments := os.Args

  if len(arguments) == 1 {
    fmt.Println("Please provide the id of a node")
    return
  }

  ports := make([]string, 3)
  var ownPort string

  switch arguments[1] {
  case "0":
    ownPort =  "127.0.0.1:9000"
    ports[0] = "127.0.0.1:9001"
    ports[1] = "127.0.0.1:9002"
    ports[2] = "127.0.0.1:9003"
  case "1":
    ownPort =  "127.0.0.1:9001"
    ports[0] = "127.0.0.1:9000"
    ports[1] = "127.0.0.1:9002"
    ports[2] = "127.0.0.1:9003"
  case "2":
    ownPort =  "127.0.0.1:9002"
    ports[0] = "127.0.0.1:9001"
    ports[1] = "127.0.0.1:9000"
    ports[2] = "127.0.0.1:9003"
  case "3":
    ownPort =  "127.0.0.1:9003"
    ports[0] = "127.0.0.1:9001"
    ports[1] = "127.0.0.1:9002"
    ports[2] = "127.0.0.1:9000"
  }

  // Create a peer instance

  intId, err := strconv.Atoi(arguments[1])
  if err != nil {
    fmt.Println("Provided arguments are wrong")
    return
  }

  s := peer.Peer{
    Id: intId,
    OwnIp: ownPort,
    Neighbours: make([]string, len(ports)),
  }
  copy(s.Neighbours, ports)

  fmt.Println(ports[0])
  ///////////
  // Routing
  ///////////

  ////// Configure listening

  lis, err := net.Listen("tcp", s.OwnIp)
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  // create a gRPC server object
  grpcServer := grpc.NewServer()

  // attach services to handler object
  peer.RegisterPeerServiceServer(grpcServer, &s)

  // Start listening in a separate go routine

  errs := make(chan error, 1)

  go func() {
    errs <- grpcServer.Serve(lis)
    close(errs)
  }()

  time.Sleep(1 * time.Second)

  /////// Configure queries

  // Get connections and clients

  connections := make([]*grpc.ClientConn, len(s.Neighbours))
  clients := make([]peer.PeerServiceClient, len(s.Neighbours))


  for i, nei := range s.Neighbours {
    connections[i], err = grpc.Dial(nei, grpc.WithInsecure())
    if err != nil {
      fmt.Println("Connection refused")
      return
    }
    clients[i] = peer.NewPeerServiceClient(connections[i])
  }
  defer CloseConnections(connections)

  /////////////////
  // Main loop
  /////////////////

  fmt.Println("Work begins")

  for {
    select {

    case err:= <-errs :
      fmt.Println("Cannot serve a request", err)
      return

    default:

      _, err := clients[0].Ping(
        context.Background(),
        &peer.PingMessage{Ok: true},
      )

      if err != nil {
        fmt.Printf("%s is not ok!. error: %s\n", s.Neighbours[0], err)
      } else {
        fmt.Printf("%s is ok\n", s.Neighbours[0])
      }

      time.Sleep(1 * time.Second)
    }
  }

}
