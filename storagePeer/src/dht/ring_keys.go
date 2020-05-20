package dht

import (
  "golang.org/x/net/context"
  "google.golang.org/grpc"
)

////////
// Local
///////

func (n *RingNode) SaveKey(key string) {

  // Add it to yourself
  n.keys = append(n.keys, key)

  // Propogate the info back
  keys := make([]string, 1)
  keys[0] = key

  ok, err := n.invokeUpdateKeysInfo(n.predecessor.IP, keys)
  if !ok || err != nil {
    panic(err)
  }
}

////////
// RPC calls
///////

func (n *RingNode) UpdateKeys(ctx context.Context, in *UpdateKeysRequest) (*UpdateReply, error) {

  n.keys = append(n.keys, in.GetKeys()...)
  // Don't forget to add them to the NewKeys channel
  
  return &UpdateReply{OK: true}, nil
}

func (n *RingNode) invokeUpdateKeys(invokeIP string, keys []string) (bool, error) {

	conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.UpdateKeys(
		context.Background(),
		&UpdateKeysRequest{Keys: keys},
	)
	if err != nil {
		return false, err
	}
	conn.Close()

	return mes.GetOK(), nil
}

func (n *RingNode) UpdateKeysInfo(ctx context.Context, in *UpdateKeysInfoRequest) (*UpdateReply, error) {

  id := in.GetID()
  ok := false

  // Decide whether theese keys are relevant to you
  if n.fingerTable[0].ID == id {

    n.succKeys = append(n.succKeys, in.GetKeys()...)
    ok = true
  } else {

    for el := n.succList.Front(); el != nil; el = el.Next() {
      val := el.Value.(neighbour)
      if val.node.ID == id {
        val.keys = append(val.keys, in.GetKeys()...)
        ok = true
      }
    }
  }

  // Decide whether we should propogate theese keys forward
  if ok {
    ok, err := n.invokeUpdateKeysInfo(n.predecessor.IP, in.GetKeys())
    if !ok || err != nil {
      panic(err)
    }
  }

  return &UpdateReply{OK: true}, nil
}

func (n *RingNode) invokeUpdateKeysInfo(invokeIP string, keys []string) (bool, error) {

	conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.UpdateKeysInfo(
		context.Background(),
		&UpdateKeysInfoRequest{Keys: keys, ID: n.self.ID},
	)
	if err != nil {
		return false, err
	}
	conn.Close()

	return mes.GetOK(), nil
}
