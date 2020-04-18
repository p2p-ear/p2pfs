package dht

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math"
	"math/big"
)

////////
// Data structures
////////

// Saved neighbour
type finger struct {
	start uint64
	IP    string
	ID    uint64
}

// RingNode is a Chord node
type RingNode struct {
	maxNodes    uint64
	self        finger
	predecessor finger
	fingerTable []finger
}

// Hash is a hash function used for node and file ids.
func Hash(data []byte, maxNum uint64) uint64 {

	hash := sha256.Sum256(data)

	gen := new(big.Int)
	gen.SetString(hex.EncodeToString(hash[:]), 16)

	base := new(big.Int)
	base.SetUint64(maxNum)

	base.Mod(gen, base)

	return base.Uint64()
}

// NewRingNode is a RingNode constructor. After constructing an object make sure to enable a gRPC server.
func NewRingNode(ownIP string, maxNodes uint64) *RingNode {

	id := Hash([]byte(ownIP), maxNodes)

	fingSize := uint64(math.RoundToEven(math.Log2(float64(maxNodes))))

	n := RingNode{
		self:        finger{IP: ownIP, ID: id, start: id},
		predecessor: finger{},
		maxNodes:    maxNodes,
		fingerTable: make([]finger, fingSize),
	}

	return &n
}

// GetPersonalSuccessor gets your node's successor
func (n *RingNode) GetPersonalSuccessor() finger {
	return finger{IP: n.fingerTable[0].IP, ID: n.fingerTable[0].ID}
}

// GetPersonalPredecessor gets your node's predecessor
func (n *RingNode) GetPersonalPredecessor() finger {
	return finger{IP: n.predecessor.IP, ID: n.predecessor.ID}
}

// MarshalJSON serializes node for printing
func (n *RingNode) MarshalJSON() ([]byte, error) {

	type PublicFinger struct {
		Start uint64
		IP    string
		ID    uint64
	}

	type PublicRingNode struct {
		MaxNodes    uint64
		Self        PublicFinger
		Predecessor PublicFinger
		FingerTable []PublicFinger
	}

	p := PublicRingNode{
		Self: PublicFinger{
			Start: n.self.start,
			IP:    n.self.IP,
			ID:    n.self.ID,
		},
		Predecessor: PublicFinger{
			Start: n.predecessor.start,
			IP:    n.predecessor.IP,
			ID:    n.predecessor.ID,
		},
		MaxNodes:    n.maxNodes,
		FingerTable: make([]PublicFinger, len(n.fingerTable)),
	}

	for i, f := range n.fingerTable {
		p.FingerTable[i].Start = f.start
		p.FingerTable[i].IP = f.IP
		p.FingerTable[i].ID = f.ID
	}

	return json.Marshal(p)
}
