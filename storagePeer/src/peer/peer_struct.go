// Peer structure definition
package peer

import (
	"storagePeer/src/dht"
)

// Peer is the peer struct
type Peer struct {
	ownIP string
	ring  *dht.RingNode
	Errs  chan error
}
