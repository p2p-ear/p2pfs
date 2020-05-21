package dht

import (
	"fmt"
	"net"
	"testing"
	"time"
	"google.golang.org/grpc"
	"math/rand"
	"strconv"
)


func startTestServ(node *RingNode) (chan error, net.Listener) {

	lis, err := net.Listen("tcp", node.self.IP)
	if err != nil {
		panic(err)
	}

	// create a gRPC server object
	grpcServer := grpc.NewServer()

	// Paerform start routine
	node.Start(grpcServer)

	// Start listening in a separate go routine

	errs := make(chan error, 1)
	go func() {
		errs <- grpcServer.Serve(lis)
		close(errs)
	}()

	return errs, lis
}

////////
// Test base situations
///////

func TestOneNode(t *testing.T) {

	var maxNum uint64 = 1000

	loner := NewRingNode("localhost:9000", maxNum, time.Second)
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
	loner.Stop()
}

func TestInInterval(t *testing.T) {

	var maxNum uint64 = 1000
	loner := NewRingNode("localhost:9000", maxNum, time.Second)

	var tests = []struct {
		start, end, id uint64
		inc_start, inc_end, want           bool
	}{
		{200, 300, 250, true, true, true},
		{200, 300, 200, false, true, false},
		{200, 300, 300, true, false, false},
		{200, 300, 400, true, true, false},
		{100, 100, 900, true, false, true},
		{900, 100, 950, true, true, true},
		{900, 200, 50, true, false, true},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%d,%d,%d", tt.start, tt.end, tt.id)
		t.Run(testname, func(t *testing.T) {
			ans := loner.inInterval(tt.start, tt.end, tt.id, tt.inc_start, tt.inc_end)
			if ans != tt.want {
				t.Errorf("got %t, want %t", ans, tt.want)
			}
		})
	}
}

func TestFingerIndex(t *testing.T) {

	var maxNum uint64 = 1000
	loner := NewRingNode("localhost:4000", maxNum, time.Second) // ID: 415

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

func findFingerSuccessor(nodes []*RingNode, id uint64) uint64 {

	succ := nodes[0].self.ID

	for i := 1; i < len(nodes); i++ {
		if succ == id {
			break
		}
		if nodes[0].inInterval(id, succ, nodes[i].self.ID, true, false) {
			succ = nodes[i].self.ID
		}
	}

	return succ
}

func findNext(nodes []*RingNode, id uint64) uint64 {
	succ := nodes[0].self.ID

	for i := 1; i < len(nodes); i++ {
		if nodes[0].inInterval(id, succ, nodes[i].self.ID, false, true) {
			succ = nodes[i].self.ID
		}
	}

	return succ
}

func findPrev(nodes []*RingNode, id uint64) uint64 {
	pred := nodes[0].self.ID

	for i := 1; i < len(nodes); i++ {
		if nodes[0].inInterval(pred, id, nodes[i].self.ID, false, false) {
			pred = nodes[i].self.ID
		}
	}

	return pred
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

func generateRing(num uint64, maxNum uint64, deltaT time.Duration, random bool) []*RingNode{

	var nodes = make([]*RingNode, num)
	var begin int

	if random {
		begin = 32768 + rand.Intn(10000)
	} else {
		begin = 30000 + int(num*2)
	}

	for i, _ := range nodes{

		nodes[i] = NewRingNode(fmt.Sprintf("localhost:%s", strconv.Itoa(begin+i)), maxNum, deltaT)
	}

	return nodes
}

////////
// Test ring construction
///////

func validateFingTable(nodes []*RingNode, t *testing.T) bool {

	firstTime := true
	// Test fingertables
	for _, node := range nodes {
		for fingIdx, fing := range node.fingerTable {
			if fing.ID != findFingerSuccessor(nodes, fing.start) {
				if firstTime {
					printNodes(nodes)
					firstTime = false
				}
				t.Errorf("On node %s, finger %d: got suc %d must be %d", node.self.IP, fingIdx, fing.ID, findFingerSuccessor(nodes, fing.start))
			}
		}
	}

	return true
}

func validateMainInfo(nodes []*RingNode, t *testing.T) {

	firstTime := true

	for _, n := range nodes {

		// Test succ
		succ := n.fingerTable[0].ID

		if findNext(nodes, n.self.ID) != succ {
			t.Errorf("Node %d has succ %d but actual is %d", n.self.ID, succ, findNext(nodes, n.self.ID))
			if firstTime {
				printNodes(nodes)
				firstTime = false
			}
		}

		// Test preds
		pred := n.predecessor.ID

		if findPrev(nodes, n.self.ID) != pred {
			t.Errorf("Node %d has pred %d but actual is %d", n.self.ID, pred, findPrev(nodes, n.self.ID))
			if firstTime {
				printNodes(nodes)
				firstTime = false
			}
		}

		// Test succ list
		if uint64(n.succList.Len()) > n.succListSize {
			t.Errorf("SuccList size is %d but must be %d at max", n.succList.Len(), n.succListSize)
		}

		for el := n.succList.Front(); el != nil; el = el.Next() {

			inlist := el.Value.(neighbour).node.ID

			if findNext(nodes, succ) == n.self.ID {
				break
			}

			if findNext(nodes, succ) != inlist {
				t.Errorf("Node %d has %d in succ list but actual is %d", n.self.ID, inlist, findNext(nodes, n.self.ID))
				if firstTime {
					printNodes(nodes)
					firstTime = false
				}
			}
			succ = findNext(nodes, succ)
		}
	}
}

func _TestJoin(t *testing.T) {

	var maxNum uint64 = 123456
	var start uint64 = 10
	var maxRingSize uint64 = 11
	var step uint64 = 10
	var deltaT time.Duration = time.Second

	fmt.Println("Testing construction of different ring topologies...")

	for i := start; i <= maxRingSize; i += step {

		fmt.Println("Ring size: ", i)

		nodes := generateRing(i, maxNum, deltaT, false)
		a := make([](chan error), len(nodes))
		b := make([](net.Listener), len(nodes))

		for j, el := range nodes {
			a[j], b[j] = startTestServ(el)
		}

		nodes[0].Join("")

		for j:= 1; j < len(nodes); j++ {
			nodes[j].Join(nodes[j-1].self.IP)
			time.Sleep(time.Millisecond * 10)
			validateFingTable(nodes[:j+1], t)
			validateMainInfo(nodes[:j+1], t)
		}

		// Close everything
		for i, _ := range nodes {
			killNode(nodes, b, i)
		}
	}
}

////////
// Test node death
///////

func killNode(nodes []*RingNode, listeners []net.Listener, i int) {

	listeners[i].Close()
	nodes[i].Stop()

}

func prepareRing(numNodes uint64, maxNum uint64, deltaT time.Duration, t *testing.T) (nodes []*RingNode, a [](chan error), b [](net.Listener)) {

	// Start everything
	nodes = generateRing(numNodes, maxNum, deltaT, false)
	a = make([](chan error), len(nodes))
	b = make([](net.Listener), len(nodes))

	for j, el := range nodes {
		a[j], b[j] = startTestServ(el)
	}

	// Connect
	nodes[0].Join("")

	for j:= 1; j < len(nodes); j++ {
		nodes[j].Join(nodes[j-1].self.IP)
		time.Sleep(time.Millisecond * 10)
	}

	validateMainInfo(nodes, t)
	validateFingTable(nodes, t)

	return nodes, a, b
}

func _TestFailures(t *testing.T) {

	var maxNum uint64 = 123456
	var deltaT time.Duration = time.Second

	var tests = []struct {
		numNodes      uint64
		waitTime      time.Duration
		deleteNum     int
	}{
		{10, 2*time.Second, 2}, // Just a delete
		{20, 5*time.Second, 3}, // Two in a row
	}

	fmt.Println("Testing fix routine...")

	for i, tt := range tests {
		testname := fmt.Sprintf("%d", i)
		t.Run(testname, func(t *testing.T) {

			nodes, _, b := prepareRing(tt.numNodes, maxNum, deltaT, t)

			// Test

			for i := 1; i <= tt.deleteNum; i++ {
				killNode(nodes, b, len(nodes)-i)
			}

			time.Sleep(tt.waitTime)
			validateMainInfo(nodes[:len(nodes)-tt.deleteNum], t)

			// Close everything
			for i, _ := range nodes[:len(nodes)-tt.deleteNum] {
				killNode(nodes[:len(nodes)-tt.deleteNum], b, i)
			}
		})
	}
}

////////
// Test file lists
///////

func findNodeIdx(nodes []*RingNode, nodeId uint64) (int){
	for i, node := range nodes {
		if node.self.ID == nodeId {
			return i
		}
	}
	panic("Haven't found node for some reason.")
}

func findSuccNode(nodes []*RingNode, id uint64) (int) {
	succID := findNext(nodes, id)
	return findNodeIdx(nodes, succID)
}

func findPredNode(nodes []*RingNode, id uint64) (int) {
	predID := findPrev(nodes, id)
	return findNodeIdx(nodes, predID)
}


func inSlice(slice []string, key string) (bool){

	ans := false
	for _, k := range slice {
		ans = ans || (k == key)
	}
	return ans
}

func checkKeys(keys []string, nodes []*RingNode, t *testing.T) {

	show := false

	for _, key := range keys {
		succIdx := findSuccNode(nodes, Hash([]byte(key), nodes[0].maxNodes))

		// Test personal key list
		if !inSlice(nodes[succIdx].keys, key) {
			t.Errorf("Node %d doesn't have it's key: %s", nodes[succIdx].self.ID, key)
			show = true
		}

		// Test predecessor if it has it in the succKeys
		predIdx := findPredNode(nodes, nodes[succIdx].self.ID)

		if !inSlice(nodes[predIdx].succKeys, key) {
			t.Errorf("Node %d doesn't have it's succesor key: %s", nodes[predIdx].self.ID, key)
			show = true
		}

		// Test multiple step predecessor if they have the key in succ list
		for i := uint64(0); i < nodes[predIdx].succListSize; i++ {
			predIdx = findPredNode(nodes, nodes[predIdx].self.ID)

			var j uint64 = 0

			for el := nodes[predIdx].succList.Front(); el != nil; el = el.Next() {

				if j == i {
					neighb := el.Value.(neighbour)
					if !inSlice(neighb.keys, key) {
						t.Errorf("Node %d doesn't have %d key: %s", nodes[predIdx].self.ID, nodes[succIdx].self.ID, key)
						show = true
					}
					break
				}
				j++
		}}
	}

	// If something went wrong - for debug
	if show {
		printNodes(nodes)
	}
}

func insertKeys(nodes []*RingNode, numKeys int) []string{

	// Generate
	keys := make([]string, numKeys)
	for i := 0; i < numKeys; i++ {
		keys[i] = fmt.Sprintf("key%d",i)
	}

	//Insert
	for _, key := range keys {
		keyHash := Hash([]byte(key), nodes[0].maxNodes)

		placeIP, err := nodes[0].FindSuccessor(keyHash)
		if err != nil {
			fmt.Println("Failed while making a call...")
			panic(err)
		}

		nodeIdx := findNodeIdx(nodes, Hash([]byte(placeIP), nodes[0].maxNodes))
		nodes[nodeIdx].SaveKey(key)

		time.Sleep(time.Millisecond * 1)
	}

	return keys
}

func TestKeyInsert(t *testing.T) {

	fmt.Println("Testing key list construction...")

	// Construct a ring

	var maxNum uint64 = 123456
	var deltaT time.Duration = time.Second * 10
	var numNodes uint64 = 15

	nodes, _, b := prepareRing(numNodes, maxNum, deltaT, t)

	// Insert some keys
	keys := insertKeys(nodes, 50)
	checkKeys(keys, nodes, t)

	// Close everything
	for i, _ := range nodes[:len(nodes)] {
		killNode(nodes[:len(nodes)], b, i)
	}
}

func TestKeyNewNode(t *testing.T) {

	fmt.Println("Check that keys get split with new node...")

	
}
