package dht

import (
  "time"
  "container/list"
  "fmt"
)
////////
// Fix routines for fallen nodes and broken finger tables
////////

func (n *RingNode) fixSuccList() {

  var res bool = false

  // Check how are your successors doing
  for el := n.succList.Front(); el != nil; el = el.Next() {

    val := el.Value.(neighbour)
    _, err := n.invokeGetPred(val.node.IP)

    res = res || (err!=nil)
  }

  // TODO: Change total reconstruction to something more intellengent (I am sorry for this :( )
  if res || (uint64(n.succList.Len()) != n.succListSize) {
    // In case they haven't updated their succ yet. Try again
    for {
      n.succList = list.New()
      err := n.initSuccList()
      if err == nil {
        //fmt.Printf("%d: trying again\n", n.self.ID)
        break
      }
    }

  }
}

func (n *RingNode) fixSuccessor() {

  // Check successor
  oldSucc := n.fingerTable[0].IP
  _, err := n.invokeGetPred(oldSucc)

  if err != nil {

    deadKeys := n.succKeys
    // Find the first alive node and give it the ownership of dead nodes keys
    for el := n.succList.Front(); el != nil; el = n.succList.Front() {

      val := el.Value.(neighbour)
      _, err := n.invokeGetPred(val.node.IP)

      if err != nil {
        deadKeys = append(deadKeys, val.keys...)
        n.succList.Remove(n.succList.Front())
        // And remember his keys
      } else {
        // Update ourselfs
        newSucc := n.succList.Front().Value.(neighbour).node
        n.fingerTable[0].IP = newSucc.IP
        n.fingerTable[0].ID = newSucc.ID
        n.succList.Remove(n.succList.Front())

        // Update this dude
        n.invokeUpdatePredecessor(n.fingerTable[0].IP)

        // Send him new keys
        ok, err := n.invokeUpdateKeys(newSucc.IP, deadKeys)
        if !ok || err != nil {
          panic(err)
        }

        // Get all of his keys as succKeys
        succKeys, err := n.invokeGetKeys(newSucc.IP)
        if err != nil {
          panic(err)
        }
        n.succKeys = succKeys

        break
      }
    }

    if (n.fingerTable[0].IP == oldSucc) && (n.succList.Len() == 0) {
      panic(fmt.Sprintf("%d lost everyone during fix", n.self.ID))
    }
  } else {
    // Check if pred points to us incase something went wrong (concurrent join)
  }

}

func (n *RingNode) fixRoutine(deltaT time.Duration) {

  for {

    time.Sleep(deltaT)

    select {
    case <- n.stopSignal:
      return
    default:
      n.fixSuccessor()
      n.fixSuccList()
    }
  }
}
