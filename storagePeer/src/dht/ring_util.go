package dht

import (
  "crypto/sha256"
	"encoding/hex"
  "math/big"
  "math"
  "google.golang.org/grpc"
)

////////
// Utility functions for all parts of the module
////////

// Hash is a hash function used for node and file ids.
func Hash(data []byte, maxNum uint64) uint64 {

	hash := sha256.Sum256(data)

	gen := new(big.Int)
	gen.SetString(hex.EncodeToString(hash[:]), 16)

	base := new(big.Int)
	base.SetUint64(maxNum)

	base.Mod(gen, base)

	return base.Uint64()
}

// Calculate i'th finger index from the current node
func (n *RingNode) fingerIndex(i int64, clockWise bool) uint64 {

	var ans uint64

	if clockWise {

		ans = (n.self.ID + uint64(math.Pow(2, float64(i)))) % n.maxNodes
	} else {

		temp := (int64(n.self.ID) - int64(math.Pow(2, float64(i)))) % int64(n.maxNodes)

		if temp < 0 {
			temp += int64(n.maxNodes)
		}

		ans = uint64(temp)
	}

	return ans
}

// Check whether id is located inside (start, end) interval.
func (n *RingNode) inInterval(start uint64, end uint64, id uint64) bool {

	start %= n.maxNodes
	end %= n.maxNodes
	id %= n.maxNodes

	if start == end {
		return true
	} else if start < end {
		return id >= start && id <= end
	}
	return (id >= start && id < n.maxNodes) || (id >= 0 && id <= end)
}

// Make a connection with other node
func getConn(ip string) (*grpc.ClientConn, RingServiceClient) {

	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	cl := NewRingServiceClient(conn)

	return conn, cl

}
