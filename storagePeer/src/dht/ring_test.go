package dht

import (
	"fmt"
	"net"
	"testing"

	"google.golang.org/grpc"
)

func startTestServ(node *RingNode) chan error {

	lis, err := net.Listen("tcp", node.self.IP)
	if err != nil {
		panic(err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// attach services to handler object
	RegisterRingServiceServer(grpcServer, node)

	// Start listening in a separate go routine

	errs := make(chan error, 1)
	go func() {
		errs <- grpcServer.Serve(lis)
		close(errs)
	}()

	return errs
}

func TestOneNode(t *testing.T) {

	var maxNum uint64 = 1000

	loner := NewRingNode("localhost:9000", maxNum)
	loner.Join("")

	contentID := Hash([]byte("This is a sample string"), maxNum)

	pred, err := loner.findPredecessor(contentID)
	if err != nil {
		panic(err)
	}
	succ, err := loner.FindSuccessor(contentID)
	if err != nil {
		panic(err)
	}

	if succ != pred.IP {
		t.Errorf("IPs are not equal: %s and %s", succ, pred.IP)
	}
	if succ != loner.self.IP {
		t.Errorf("Got the wrong value: %s", succ)
	}
}

func TestInInterval(t *testing.T) {

	var maxNum uint64 = 1000
	loner := NewRingNode("localhost:9000", maxNum)

	var tests = []struct {
		start, end, id uint64
		want           bool
	}{
		{200, 300, 250, true},
		{200, 300, 400, false},
		{100, 100, 900, true},
		{900, 100, 950, true},
		{900, 200, 50, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d,%d,%d", tt.start, tt.end, tt.id)
		t.Run(testname, func(t *testing.T) {
			ans := loner.inInterval(tt.start, tt.end, tt.id)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestFingerIndex(t *testing.T) {

	var maxNum uint64 = 1000
	loner := NewRingNode("localhost:4000", maxNum) // ID: 415

	var tests = []struct {
		i         int64
		clockWise bool
		want      uint64
	}{
		{0, true, 416},
		{10, true, 439},
		{0, false, 414},
		{10, false, 391},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d, %t", tt.i, tt.clockWise)
		t.Run(testname, func(t *testing.T) {
			ans := loner.fingerIndex(tt.i, tt.clockWise)
			if ans != tt.want {
				t.Errorf("got %d, want %d", ans, tt.want)
			}
		})
	}
}

func findListSuccessor(nodes []*RingNode, id uint64) uint64 {

	succ := nodes[0].self.ID

	for i := 1; i < len(nodes); i++ {
		if nodes[0].inInterval(id, succ, nodes[i].self.ID) {
			succ = nodes[i].self.ID
		}
	}

	return succ
}

func printNode(node *RingNode, i int) {
	js, err := node.MarshalJSON()
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d: %s\n", i, string(js))
}
func printNodes(nodes []*RingNode) {

	for i, el := range nodes {
		printNode(el, i)
	}
}

func TestJoin(t *testing.T) {

	var maxNum uint64 = 1000

	var nodes = []*RingNode{NewRingNode("localhost:8003", maxNum), NewRingNode("localhost:8004", maxNum), NewRingNode("localhost:8005", maxNum), NewRingNode("localhost:8008", maxNum)}
	a := make([](chan error), len(nodes))

	for i, el := range nodes {
		a[i] = startTestServ(el)
	}

	nodes[0].Join("")
	nodes[1].Join(nodes[0].self.IP)
	nodes[2].Join(nodes[1].self.IP)
	nodes[3].Join(nodes[0].self.IP)

	// Test predecessors
	for _, node := range nodes {
		pred := node.predecessor.ID
		succ := findListSuccessor(nodes, pred+1)
		if succ != node.self.ID {
			t.Errorf("On node %s: got predeccessor %d, pred has actual succ %d", node.self.IP, pred, succ)
		}
	}

	firstTime := true
	// Test fingertables
	for _, node := range nodes {
		for fingIdx, fing := range node.fingerTable {
			if fing.ID != findListSuccessor(nodes, fing.start) {
				if firstTime {
					printNodes(nodes)
					firstTime = false
				}
				t.Errorf("On node %s, finger %d: got suc %d must be %d", node.self.IP, fingIdx, fing.ID, findListSuccessor(nodes, fing.start))
			}
		}
	}

}
