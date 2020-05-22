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

  ok, err := n.invokeUpdateKeysInfo(n.predecessor.IP, n.self.ID, keys)
  if !ok || err != nil {
    panic(err)
  }
}

////////
// RPC calls
///////

/////////// Update personal keys

func (n *RingNode) UpdateKeys(ctx context.Context, in *UpdateKeysRequest) (*UpdateReply, error) {

  // Add them to the key list
  n.keys = append(n.keys, in.GetKeys()...)

  // Send them to the NewFilesChannel for higher level software to take care of it
  for _, k := range in.GetKeys() {
    n.NewFilesChannel <- k
  }

  // Backpropogate info about new keys
  ok, err := n.invokeUpdateKeysInfo(n.predecessor.IP, n.self.ID, in.GetKeys())
  if !ok || err != nil {
    panic(err)
  }

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

/////////// Update information about keys of others

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
        el.Value = val
        ok = true
      }
    }
  }

  // Decide whether we should propogate theese keys forward
  if ok {
    ok, err := n.invokeUpdateKeysInfo(n.predecessor.IP, id, in.GetKeys())
    if !ok || err != nil {
      panic(err)
    }
  }

  return &UpdateReply{OK: true}, nil
}

func (n *RingNode) invokeUpdateKeysInfo(invokeIP string,  updateForID uint64, keys []string) (bool, error) {

	conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return false, err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.UpdateKeysInfo(
		context.Background(),
		&UpdateKeysInfoRequest{Keys: keys, ID: updateForID},
	)
	if err != nil {
		return false, err
	}
	conn.Close()

	return mes.GetOK(), nil
}

/////////// Get someones keys

func (n *RingNode) GetKeys(ctx context.Context, in *GetKeysRequest) (*KeyReply, error) {

  return &KeyReply{Keys: n.keys}, nil
}

func (n *RingNode) invokeGetKeys(invokeIP string) ([]string, error) {

  conn, err := grpc.Dial(invokeIP, grpc.WithInsecure())
	if err != nil {
		return make([]string,0), err
	}
	cl := NewRingServiceClient(conn)

	mes, err := cl.GetKeys(
		context.Background(),
		&GetKeysRequest{},
	)
	if err != nil {
		return make([]string,0), err
	}
	conn.Close()

	return mes.GetKeys(), nil
}
