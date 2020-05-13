package dht

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	//"fmt"
)

////////
// Update state of other nodes
////////

// UpdatePredecessor updates predecessor of a node with a requesting one
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

// UpdateSpecificFinger updates i'th finger of a node
func (n *RingNode) UpdateSpecificFinger(ctx context.Context, in *UpdateSpecificFingerRequest) (*UpdateReply, error) {

	go func() {
		s := finger{ID: in.GetID(), IP: in.GetIP()}
		i := in.GetFingID()

		if n.inInterval(n.self.ID, n.fingerTable[i].ID, s.ID) {

			n.fingerTable[i].ID = s.ID
			n.fingerTable[i].IP = s.IP

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
