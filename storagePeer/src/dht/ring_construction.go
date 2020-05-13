package dht


////////
// Code for initial join of a node to the ring
////////

// Join initializes finger table with peers' values
func (n *RingNode) Join(existingIP string) {

	if len(existingIP) == 0 {
		//First join
		for i := int64(0); i < int64(len(n.fingerTable)); i++ {
			n.fingerTable[i].start = n.fingerIndex(i, true) // TODO: Might want to refactor it to create a new Node structure
			n.fingerTable[i].IP = n.self.IP
			n.fingerTable[i].ID = n.self.ID
		}
		n.predecessor = n.fingerTable[len(n.fingerTable)-1] // TODO: panics when one node in network
	} else {

		n.initFingerTable(existingIP)
		n.updateOthers()
	}
}

// Get information about your neighbours
func (n *RingNode) initFingerTable(existingIP string) {

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
func (n *RingNode) updateOthers() {

	for i := int64(0); i < int64(len(n.fingerTable)); i++ {

		p := n.recursivePredFindingStep(n.fingerIndex(i, false), n.fingerTable[0], n.self) // Don't use your own table

		if p.IP == n.self.IP {
			break
		}

		n.invokeUpdateSpecificFinger(p.IP, i, n.self)

	}
}
