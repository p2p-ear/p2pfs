package dht

import (
	"encoding/json"
	"math"
	"container/list"
)




////////
// This is a realisation of a Chord ring - names of the methods are consistent with the paper.
// title: Chord: A Scalable Peer-to-peer Lookup Service for Internet Applications
// authors: Ion Stoica and others
// link: https://pdos.csail.mit.edu/papers/chord:sigcomm01/chord_sigcomm.pdf
////////

// Theese constants specify system requirements
const FAIL_PROB = 0.1 // probability that one node will fail in time delta T
const TOLERABLE_FAIL_PROB = 0.001 // tolerable probability of failure (epsilon)

////////
// Data structures
////////

// Saved neighbour
type finger struct {
	start uint64
	IP    string
	ID    uint64
}

// Node and it's contents
type neighbour struct {
	node finger
	contents list.List
}

// RingNode is a Chord node
type RingNode struct {
	maxNodes    uint64
	self        finger
	predecessor finger

	fingerTable []finger
	succList    *list.List
	succListSize uint64
}

// NewRingNode is a RingNode constructor. After constructing an object make sure to enable a gRPC server.
func NewRingNode(ownIP string, maxNodes uint64) *RingNode {

	id := Hash([]byte(ownIP), maxNodes)

	fingSize := uint64(math.RoundToEven(math.Log2(float64(maxNodes))))
	succListSize := uint64(math.Log(TOLERABLE_FAIL_PROB) / math.Log(FAIL_PROB)) - 1 // this -1 apears since first successor is in the fingertable, it's convinient

	n := RingNode{
		self:        finger{IP: ownIP, ID: id, start: id},
		predecessor: finger{},
		maxNodes:    maxNodes,
		fingerTable: make([]finger, fingSize),
		succList:    list.New(), // at first it's empty
		succListSize: succListSize,
	}

	return &n
}

// MarshalJSON serializes node for printing
func (n *RingNode) MarshalJSON() ([]byte, error) {

	type PublicFinger struct {
		Start uint64
		IP    string
		ID    uint64
	}

	type PublicNeighbour struct {
		Node     PublicFinger
		Contents list.List
	}

	type PublicRingNode struct {
		MaxNodes    uint64
		Self        PublicFinger
		Predecessor PublicFinger
		FingerTable []PublicFinger
		SuccList    []PublicNeighbour
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
		SuccList:    make([]PublicNeighbour, n.succListSize),
	}

	for i, f := range n.fingerTable {
		p.FingerTable[i].Start = f.start
		p.FingerTable[i].IP = f.IP
		p.FingerTable[i].ID = f.ID
	}

	i:=0
	for el := n.succList.Front(); el !=nil; el = el.Next() {
		p.SuccList[i].Node = PublicFinger{ID: el.Value.(neighbour).node.ID, IP: el.Value.(neighbour).node.IP}
		p.SuccList[i].Contents = el.Value.(neighbour).contents
		i++
	}

	return json.Marshal(p)
}
