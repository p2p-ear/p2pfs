package dht

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"

  "crypto/sha256"
  "math"
	"math/big"

	"encoding/hex"
  "encoding/json"
)

////////
// Data structures
////////

// Saved neighbour
type finger struct {
	start      uint64
	end        uint64
	ip         string
}

// This is a Chord node
type RingNode struct {
  id          uint64
	ownIP       string
	maxNodes    uint64
	fingerTable []finger
}

////////
// Local functions
////////

func Hash(data []byte, maxNum uint64) uint64 {

	hash := sha256.Sum256(data)

	gen := new(big.Int)
	gen.SetString(hex.EncodeToString(hash[:]), 16)

	base := new(big.Int)
	base.SetUint64(maxNum)

	base.Mod(gen, base)

	return base.Uint64()
}

// RingNode constructor
func NewRingNode(ownIP string, maxNodes uint64) *RingNode {

  id := Hash([]byte(ownIP), maxNodes)

  fingSize := uint64(math.RoundToEven(math.Log2(float64(maxNodes))))

	n := RingNode{
    ownIP: ownIP,
    id: id,
    maxNodes: maxNodes,
    fingerTable: make([]finger, fingSize),
  }

	return &n
}


// Serialization
func (n *RingNode) MarshalJSON() ([]byte, error) {

  type PublicFinger struct {
    Start uint64
    End   uint64
    IP    string
  }

  type PublicRingNode struct {
    ID    uint64
    OwnIP string
    MaxNodes uint64
    FingerTable []PublicFinger
  }

  p := PublicRingNode {
    ID: n.id,
    OwnIP: n.ownIP,
    MaxNodes: n.maxNodes,
    FingerTable: make([]PublicFinger, len(n.fingerTable)),
  }

  for i, f := range n.fingerTable {
    p.FingerTable[i].Start = f.start
    p.FingerTable[i].End = f.end
    p.FingerTable[i].IP = f.ip
  }

  return json.Marshal(p)
}

////////
// Chord operations (names are the same as in the original paper)
////////

// Initialize finger table with peers' values
func (n* RingNode) Join(existingIP string) {

  n.fingerTable[0].start = Hash([]byte(existingIP), n.maxNodes)
  n.fingerTable[0].end   = n.id
  n.fingerTable[0].ip    = existingIP
}

// Find successor in a ring
func (n* RingNode) FindSuccessor(id uint64) (string, error) {

  conn, err := grpc.Dial(n.fingerTable[0].ip, grpc.WithInsecure())
  if err != nil {
    panic(err)
  }
  cl := NewRingServiceClient(conn)

  mes, err := cl.GetSucc(
    context.Background(),
    &GetSuccRequest{ID: id,},
  )

  conn.Close()

  return mes.GetIP(), err
}

////////
// Remote calls
////////

func (n *RingNode) GetSucc(ctx context.Context, in *GetSuccRequest) (*GetSuccReply, error) {

  return &GetSuccReply{IP: n.ownIP, ID: n.id}, nil
}
