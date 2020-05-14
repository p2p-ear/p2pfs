package dht

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"fmt"
)

////////
// Update state of other nodes
////////


/////////////////// Predecessor


// UpdatePredecessor updates predecessor of a node with a requesting one
func (n *RingNode) UpdatePredecessor(ctx context.Context, in *UpdatePredRequest) (*UpdateReply, error) {

	ip := in.IP
	id := Hash([]byte(ip), n.maxNodes)

	//fmt.Printf("Curr pred: %d, self: %d, id: %d, IP: %s", n.predecessor.ID, n.self.ID, id, ip)
	ok := n.inInterval(n.predecessor.ID, n.self.ID, id, true, false)

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


/////////////////// Fingers


// UpdateSpecificFinger updates i'th finger of a node
func (n *RingNode) UpdateSpecificFinger(ctx context.Context, in *UpdateSpecificFingerRequest) (*UpdateReply, error) {

	s := finger{ID: in.GetID(), IP: in.GetIP()}
	i := in.GetFingID()

	if i == 0 {
		// Use update succ for this
		return &UpdateReply{OK: false}, nil
	}

	//fmt.Printf("Asking %d to change %dth finger from %d to %d\n", n.self.ID, i, n.fingerTable[i].ID, s.ID)
	if n.inInterval(n.self.ID, n.fingerTable[i].ID, s.ID, true, false) {

		n.fingerTable[i].ID = s.ID
		n.fingerTable[i].IP = s.IP

		// n.predecessor has already included itself in fingertable buring construction (if necessary)
		if n.predecessor.IP != s.IP {
			// Propogate change
			// TODO: check results and don't just forget about this func
			go func() {
				_, err := n.invokeUpdateSpecificFinger(n.predecessor.IP, i, s)
				if err != nil {
					panic(err)
				}
			}()
		}
	}

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

/////////////////// Successor

// Insert node into succlist and make sure that it's <= succListSize
// Returns true if element was inserted
func (n *RingNode) insertToSuccList(node finger) bool {

	// Finding our place
	for neighb := n.succList.Front(); neighb != nil; neighb = neighb.Next() {

		if n.inInterval(n.fingerTable[0].ID, neighb.Value.(neighbour).node.ID, node.ID, false, false) {

			el := n.succList.InsertBefore(neighbour{node: node}, neighb)
			if el == nil {
				panic("Couldn't insert a node into succ list!")
			}

			if uint64(n.succList.Len()) > n.succListSize {
				n.succList.Remove(n.succList.Back())
			}
			return true
		}
	}

	// Insert first element or in the end
	if uint64(n.succList.Len()) < n.succListSize {
		n.succList.PushBack(neighbour{node: node})
		return true
	}

	return false
}

// Insert new succ and update successor lists on previous nodes
func (n *RingNode) UpdateSucc(ctx context.Context, in *UpdateSuccRequest) (*UpdateReply, error) {

	ip := in.IP
	id := Hash([]byte(ip), n.maxNodes)

	//fmt.Printf("update succ: %d is updated with %d\n", n.self.ID, id)

	oldSuc := n.fingerTable[0]

	n.fingerTable[0].ID = id; n.fingerTable[0].IP = ip

	// Check if it's second node joining
	if oldSuc.ID != n.self.ID {

		if !n.insertToSuccList(oldSuc) {
			panic("Couldn't insert old suc")
		}

		// Propogate change
		go func() {
			_, err := n.invokeUpdateSuccList(n.predecessor.IP, finger{IP: ip, ID: id})
			if err != nil {
				panic(err)
			}
		}()
	}

	return &UpdateReply{OK:true}, nil
}

func (n* RingNode) invokeUpdateSucc(invokeIP string, node finger) (bool, error) {

	conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.UpdateSucc(
		context.Background(),
		&UpdateSuccRequest{IP: node.IP},
	)
	if err != nil {
		return false, err
	}
	conn.Close()

	return mes.GetOK(), nil
}

/////////////////// Successor list

func (n *RingNode) UpdateSuccList(ctx context.Context, in *UpdateSuccListRequest) (*UpdateReply, error) {

	ip := in.IP
	id := Hash([]byte(ip), n.maxNodes)

	//fmt.Printf("update succ list: %d is updated with %d, size %d\n", n.self.ID, id, n.succList.Len())

	// Don't add yourself to a succ list!
	if id != n.self.ID {

		if n.insertToSuccList(finger{IP: ip, ID: id}) {
			// Propogate only changes you made yourself
			go func() {
				_, err := n.invokeUpdateSuccList(n.predecessor.IP, finger{IP: ip, ID: id})
				if err != nil {
					panic(err)
				}
			}()
		}
	}

	return &UpdateReply{OK: true}, nil
}

func (n *RingNode) invokeUpdateSuccList(invokeIP string, node finger) (bool, error) {

	conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.UpdateSuccList(
		context.Background(),
		&UpdateSuccListRequest{IP: node.IP},
	)
	if err != nil {
		return false, err
	}
	conn.Close()

	return mes.GetOK(), nil
}
