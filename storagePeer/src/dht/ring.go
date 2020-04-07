// dht is a package which implments chord distributed hash table.
package dht

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
  "math"
  "errors"
  //"fmt"
)

////////
// Chord operations (names are the same as in the original paper)
// title: Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications
// authors: Ion Stoica and others
// link: https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf
////////

//////////////// New node joining ////////////////

// Calculate i'th finger index from the current node
func (n* RingNode) fingerIndex(i int64, clockWise bool) uint64 {

  var ans uint64

  if clockWise {

    ans = (n.self.ID + uint64(math.Pow(2, float64(i)))) % n.maxNodes
  } else {

    temp := (int64(n.self.ID) - int64(math.Pow(2, float64(i)))) % int64(n.maxNodes)

    if temp < 0 {
      temp += int64(n.maxNodes)
    }

    ans = uint64(temp)
  }

  return ans
}

// Initialize finger table with peers' values
func (n* RingNode) Join(existingIP string) {

  if len(existingIP) == 0 {
    //First join
    for i := int64(0); i < int64(len(n.fingerTable)); i++ {
      n.fingerTable[i].start = n.fingerIndex(i, true) // TODO: Might want to refactor it to create a new Node structure
      n.fingerTable[i].IP    = n.self.IP
      n.fingerTable[i].ID    = n.self.ID
    }
    n.predecessor = n.fingerTable[len(n.fingerTable)-1]
  } else {

    n.initFingerTable(existingIP)
    n.updateOthers()
  }
}

// Get information about your neighbours
func (n* RingNode) initFingerTable(existingIP string) {

  //fmt.Printf("===Initiating for node %s with remote %s===\n", n.self.IP, existingIP)

  existingNode := finger{IP: existingIP, ID: Hash([]byte(existingIP), n.maxNodes)}

  // First get successor and predecessor
  n.predecessor = n.recursivePredFindingStep(n.self.ID, existingNode, n.self)

  succ, err := n.invokeGetSucc(n.predecessor.IP)
  if err != nil {
    panic(err)
  }

  n.fingerTable[0] = succ
  n.fingerTable[0].start = n.self.ID + 1

  // Insert new node as this node's predecessor

  ok, err := n.invokeUpdatePredecessor(succ.IP)

  if err != nil || !ok {
    panic(err)
  }

  // Now update other entries in finger table
  for i := int64(1); i < int64(len(n.fingerTable)); i++ {

    start := n.fingerIndex(i, true)

    //fmt.Printf("finger %d with start %d\n", i, start)

    if n.inInterval(n.predecessor.ID, n.self.ID, start) {
    // It means that new node is responsible for theese keys
      n.fingerTable[i] = n.self

    } else {
      if n.inInterval(n.self.ID, n.fingerTable[i-1].ID, start) {

        n.fingerTable[i] = n.fingerTable[i-1]
        //fmt.Printf("using old one\n")

      } else {

        pred := n.recursivePredFindingStep(start, existingNode, n.self)
        succ, err := n.invokeGetSucc(pred.IP)
        if err != nil {
          panic(err)
        }

        n.fingerTable[i] = succ
        //fmt.Printf("got pred %s\n", pred.IP)
        //fmt.Printf("asked %s got %s\n", existingIP, succ.IP)
      }
    }
    n.fingerTable[i].start = start
  }
}

// Update finger tables of other nodes
func (n* RingNode) updateOthers() {

  for i := int64(0); i < int64(len(n.fingerTable)); i++ {

    p := n.recursivePredFindingStep( n.fingerIndex(i, false), n.fingerTable[0], n.self ) // Don't use your own table

    if p.IP == n.self.IP {
      break
    }

    n.invokeUpdateSpecificFinger(p.IP, i, n.self)

  }
}

//////////////// Finding responsible nodes ////////////////

// Check whether id is located inside (start, end) interval.
func (n* RingNode) inInterval(start uint64, end uint64, id uint64) bool{

  start %= n.maxNodes; end %= n.maxNodes; id %= n.maxNodes

  if start == end {
    return true
  } else if start < end {
    return id >= start && id <= end
  }
  return (id >= start && id < n.maxNodes) || (id >= 0 && id <= end)
}

// Find closest predecessing finger from the personal table
func (n* RingNode) getClosestPreceding(id uint64) (finger, error) {

  if n.inInterval(n.self.ID, n.fingerTable[0].ID, id) {
    //fmt.Printf("I am %s and the answer is me\n", n.self.IP)
    return n.self, nil
  }

  for i := len(n.fingerTable) - 1; i >= 0; i-- {
    if n.fingerTable[i].IP != n.self.IP && n.inInterval(n.self.ID, id, n.fingerTable[i].ID) { // TODO: this is a hack! just refactor inInterval not to include start
      //fmt.Printf("I am %s and the answer is %s\n", n.self.IP, n.fingerTable[i].IP)
      return n.fingerTable[i], nil
    }
  }

  return finger{}, errors.New("Couldn't find a place on a circle.")
}

// Going to reuse this function in findPredecessor and during
// construction of finger table.
func (n* RingNode) recursivePredFindingStep(id uint64, remoteNode finger, currNode finger) finger {

  if remoteNode.ID == currNode.ID {
    return currNode
  }

  next, err := n.invokeFindPred(remoteNode.IP, id)
  if err != nil {
    panic(err)
  }

  //fmt.Printf("recursive step: id %d remote %s ans %s \n", id, remoteNode.IP, next.IP)

  return n.recursivePredFindingStep(id, next, remoteNode)
}

// Find successor in a ring
func (n* RingNode) findPredecessor(id uint64) (finger, error) {

  nextTarget, err := n.getClosestPreceding(id)
  if err != nil {
    return finger{}, err
  }

  // Ask the node for the closest ones in its table recursivly
  return n.recursivePredFindingStep(id, nextTarget, n.self), nil
}

// Find successor node for a given id
func (n* RingNode) FindSuccessor(id uint64) (string, error) {

  pred, err := n.findPredecessor(id)
  if err != nil {
    panic(err)
  }

  if pred.IP == n.self.IP {
    return n.self.IP, nil
  }

  ans, err := n.invokeGetSucc(pred.IP)
  return ans.IP, err
}

//////////////// Remote calls ////////////////

func getConn(ip string) (*grpc.ClientConn, RingServiceClient) {

  conn, err := grpc.Dial(ip, grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  cl := NewRingServiceClient(conn)

  return conn, cl

}

//// Simple getters
// Get successor of a node
func (n *RingNode) GetNodeSucc(ctx context.Context, in *GetNodeSuccRequest) (*NodeReply, error) {
  return &NodeReply{IP: n.fingerTable[0].IP, ID: n.fingerTable[0].ID}, nil
}

func (n *RingNode) invokeGetSucc(IP string) (finger, error) {

  conn, cl := getConn(IP)

  mes, err := cl.GetNodeSucc(
    context.Background(),
    &GetNodeSuccRequest{},
  )
  if err != nil {
    return finger{}, err
  }
  conn.Close()

  return finger{ID: mes.GetID(), IP: mes.GetIP()}, nil
}

// Get predecessor of a node
func (n* RingNode) GetNodePred(ctx context.Context, in *GetNodePredRequest) (*NodeReply, error) {
  return &NodeReply{IP: n.predecessor.IP, ID: n.predecessor.ID}, nil
}
func (n* RingNode) invokeGetPred(IP string) (finger, error) {

  conn, cl := getConn(IP)

  mes, err := cl.GetNodePred(
    context.Background(),
    &GetNodePredRequest{},
  )
  if err != nil {
    return finger{}, err
  }
  conn.Close()

  return finger{ID: mes.GetID(), IP: mes.GetIP()}, nil
}

//// Recursive finding
// Find predecessor of certain id
func (n *RingNode) FindPred(ctx context.Context, in *FindPredRequest) (*NodeReply, error) {
  res, err := n.getClosestPreceding(in.ID)
  return &NodeReply{IP: res.IP, ID: res.ID}, err
}

func (n *RingNode) invokeFindPred(invokeIP string, id uint64) (finger, error){

  conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
  if err != nil {
    return finger{}, err
  }
  cl := NewRingServiceClient(conn)

  mes, err := cl.FindPred(
    context.Background(),
    &FindPredRequest{ID: id,},
  )
  if err != nil {
    return finger{}, err
  }
  conn.Close()

  return finger{ID: mes.GetID(), IP: mes.GetIP()}, nil
}

//// Update requests
// Update predecessor of a node with a requesting one
func (n *RingNode) UpdatePredecessor(ctx context.Context, in *UpdatePredRequest) (*UpdateReply, error) {

  ip := in.IP
  id := Hash([]byte(ip), n.maxNodes)

  ok := n.inInterval(n.predecessor.ID, n.self.ID, id)

  if ok {
    n.predecessor = finger{ID: id, IP: ip}
  }

  return &UpdateReply{OK: ok}, nil
}

func (n *RingNode) invokeUpdatePredecessor(invokeIP string) (bool, error) {

  conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
  if err != nil {
    return false, err
  }
  cl := NewRingServiceClient(conn)

  mes, err := cl.UpdatePredecessor(
    context.Background(),
    &UpdatePredRequest{IP: n.self.IP},
  )
  if err != nil {
    return false, err
  }
  conn.Close()

  return mes.GetOK(), nil
}

// Update i'th finger of a node
func (n *RingNode) UpdateSpecificFinger(ctx context.Context, in *UpdateSpecificFingerRequest) (*UpdateReply, error) {

  go func() {
    s := finger{ID: in.GetID(), IP: in.GetIP()}
    i := in.GetFingID()

    if n.inInterval(n.self.ID, n.fingerTable[i].ID, s.ID) {

      n.fingerTable[i].ID = s.ID; n.fingerTable[i].IP = s.IP

      // n.predecessor has already included itself in fingertable buring construction (if necessary)
      if n.predecessor.IP != s.IP {
        _, err := n.invokeUpdateSpecificFinger(n.predecessor.IP, i, s)
        if err != nil {
          panic(err)
        }
      }
    }
  }()

  return &UpdateReply{OK: true}, nil
}

func (n *RingNode) invokeUpdateSpecificFinger(invokeIP string, fingIndex int64, node finger) (bool, error) {

  conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
  if err != nil {
    return false, err
  }
  cl := NewRingServiceClient(conn)

  mes, err := cl.UpdateSpecificFinger(
    context.Background(),
    &UpdateSpecificFingerRequest{FingID: fingIndex, ID: node.ID, IP: node.IP},
  )
  if err != nil {
    return false, err
  }
  conn.Close()

  return mes.GetOK(), nil
}
