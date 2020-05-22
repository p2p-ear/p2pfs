package dht

import (
	"encoding/json"
	"math"
	"container/list"
	"time"
	"google.golang.org/grpc"
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

// Node structure
type finger struct {
	start uint64
	IP    string
	ID    uint64
}

// Node and it's contents
type neighbour struct {
	node finger
	keys []string
}

// RingNode is a Chord node
type RingNode struct {
	maxNodes    uint64

	// Ring information
	self        finger
	predecessor finger

	fingerTable []finger
	succList    *list.List
	succListSize uint64

	// Keys information
	keys            []string
	succKeys        []string
	keysStartSize   int
	NewFilesChannel chan string

	// Fix routine information
	stopSignal  chan struct{}
	deltaT      time.Duration
}

// NewRingNode is a RingNode constructor. After constructing an object make sure to enable a gRPC server.
func NewRingNode(ownIP string, maxNodes uint64, deltaT time.Duration) *RingNode {

	id := Hash([]byte(ownIP), maxNodes)

	fingSize := uint64(math.RoundToEven(math.Log2(float64(maxNodes))))
	succListSize := uint64(math.Log(TOLERABLE_FAIL_PROB) / math.Log(FAIL_PROB)) - 1 // this -1 apears since first successor is in the fingertable, it's convinient

	const keysStartSize = 0

	n := RingNode{
		self:        finger{IP: ownIP, ID: id, start: id},
		predecessor: finger{},
		maxNodes:    maxNodes,
		fingerTable: make([]finger, fingSize),
		succList:    list.New(), // at first it's empty
		succListSize: succListSize,
		stopSignal: make(chan struct{}),
		deltaT: deltaT,
		keys: make([]string, keysStartSize),
		succKeys: make([]string, keysStartSize),
		keysStartSize: keysStartSize,
		NewFilesChannel: make(chan string, 100),
	}

	return &n
}

// Connect service to the gRPC server and start fix routine
func (n *RingNode) Start(grpcServer *grpc.Server) (){

	RegisterRingServiceServer(grpcServer, n)

}

// Gracefull shutdown
func (n *RingNode) Stop() {
	close(n.stopSignal)
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
		Keys     []string
	}

	type PublicRingNode struct {
		MaxNodes    uint64
		Self        PublicFinger
		Predecessor PublicFinger
		FingerTable []PublicFinger
		SuccList    []PublicNeighbour
		Keys        []string
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
		Keys:        make([]string, len(n.keys)),
	}

	for i, f := range n.fingerTable {
		p.FingerTable[i].Start = f.start
		p.FingerTable[i].IP = f.IP
		p.FingerTable[i].ID = f.ID
	}

	for i, k := range n.keys {
		p.Keys[i] = k
	}

	neighbIdx:=0
	for el := n.succList.Front(); el !=nil; el = el.Next() {
		p.SuccList[neighbIdx].Node = PublicFinger{ID: el.Value.(neighbour).node.ID, IP: el.Value.(neighbour).node.IP}
		p.SuccList[neighbIdx].Keys = make([]string, len(el.Value.(neighbour).keys))
		for i, k := range el.Value.(neighbour).keys {
			p.SuccList[neighbIdx].Keys[i] = k
		}
		neighbIdx++
	}

	return json.Marshal(p)
}
