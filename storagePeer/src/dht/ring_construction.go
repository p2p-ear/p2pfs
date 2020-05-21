package dht

import (
	"fmt"
)

////////
// Code for initial join of a node to the ring
////////

// Join initializes finger table with peers' values
func (n *RingNode) Join(existingIP string) {

	if len(existingIP) == 0 {

		// First join

		// Init finger table
		for i := int64(0); i < int64(len(n.fingerTable)); i++ {
			n.fingerTable[i].start = n.fingerIndex(i, true) // TODO: Might want to refactor it to create a new Node structure
			n.fingerTable[i].IP = n.self.IP
			n.fingerTable[i].ID = n.self.ID
		}

		// Init a predecessor
		n.predecessor = n.fingerTable[len(n.fingerTable)-1] // TODO: panics when one node in network

	} else {

		// First initialize your succ and pred
		n.initClosest(existingIP)

		// First we deal with succ lists since they are important for correctness
		err := n.initSuccList()
		if err != nil {
			panic(err)
		}

		// Now finger tables for perfomance
		n.initFingerTable()
		n.updateOthersFingerTables()
	}

	// Launch fix routine
	go n.fixRoutine(n.deltaT)
}

///// Say hello to your closest friends

// First init your succ and pred and tell them about yourself
func (n *RingNode) initClosest(existingIP string) {

	existingNode := finger{IP: existingIP, ID: Hash([]byte(existingIP), n.maxNodes)}

	// First get successor and predecessor
	n.predecessor = n.recursivePredFindingStep(n.self.ID, existingNode, n.self)

	succ, err := n.invokeGetSucc(n.predecessor.IP)
	if err != nil {
		panic("Couldn't get a successor!")
	}

	n.fingerTable[0] = succ
	n.fingerTable[0].start = n.self.ID + 1

	// Update others
	n.insertYourself(n.predecessor.IP, succ.IP)

}

// Insert yourself as succ and pred of neighbour node
func (n *RingNode) insertYourself(predIP string, succIP string) {

	// Insert new node as succ's predecessor
	ok, err := n.invokeUpdatePredecessor(succIP)

	if err != nil || !ok {
		fmt.Println(ok)
		panic("Couldn't update pred!")
	}

	// Now as pred's successor
	n.invokeUpdateSucc(predIP, n.self)
}

///// Succ lists

// Get information about closest neighbours
func (n *RingNode) initSuccList() error{

	first, err := n.invokeGetSucc(n.fingerTable[0].IP)
	if err != nil {
		return err
	}

	n.succList.PushBack(neighbour{node:finger{IP: first.IP, ID: first.ID}, keys:make([]string,0)})

	for i := uint64(1); i < n.succListSize; i++ {

		node, err := n.invokeGetSucc(n.succList.Back().Value.(neighbour).node.IP)
		if err != nil {
			return err
		}

		n.succList.PushBack(neighbour{node:finger{IP: node.IP, ID: node.ID}, keys:make([]string,0)})
	}

	return nil
}

///// FT
// Get information about your neighbours
func (n *RingNode) initFingerTable() {

	//fmt.Printf("===Initiating for node %s with remote %s===\n", n.self.IP, existingIP)

	// Now update other entries in finger table
	for i := int64(1); i < int64(len(n.fingerTable)); i++ {

		start := n.fingerIndex(i, true)

		//fmt.Printf("finger %d with start %d\n", i, start)

		if n.inInterval(n.predecessor.ID, n.self.ID, start, true, false) {
			// It means that new node is responsible for theese keys
			n.fingerTable[i] = n.self

		} else {
			if n.inInterval(n.self.ID, n.fingerTable[i-1].ID, start, true, false) {

				n.fingerTable[i] = n.fingerTable[i-1]
				//fmt.Printf("using old one\n")

			} else {

				pred := n.recursivePredFindingStep(start, n.fingerTable[0], n.self)
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
func (n *RingNode) updateOthersFingerTables() {

	//fmt.Printf("Node %d is updating others\n", n.self.ID)

	for i := int64(0); i < int64(len(n.fingerTable)); i++ {

		p := n.recursivePredFindingStep(n.fingerIndex(i, false), n.fingerTable[0], n.self) // Don't use your own table

		if p.IP == n.self.IP {
			break
		}

		//fmt.Printf("For %d: %d id got %d\n", i, n.fingerIndex(i, false), p.ID)

		// If this anticlockwise finger hits some node exactly then we also have to change it's fingers
		succ, err := n.invokeGetSucc(p.IP)
		if err != nil {
			panic(err)
		}

		var target string
		target = p.IP

		if succ.ID == n.fingerIndex(i, false) {
			target = succ.IP
		} else {
			target = p.IP
		}

		n.invokeUpdateSpecificFinger(target, i, n.self)

	}
}
