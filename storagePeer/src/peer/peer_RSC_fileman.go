// Ring-targeted filemanagement w\ ReedSolomonCodes
package peer

import (
	"fmt"
	"math"
	"os"

	"github.com/klauspost/reedsolomon"
)

const dataRSC = 8
const parityRSC = 2

func getShardName(fname string, number int) string {
	return fmt.Sprintf("%s_rep%d", fname, number)
}

// UploadFileRSC - like UploadFile but with Reed-Solomon erasure coding
func UploadFileRSC(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) error {
	enc, err := reedsolomon.New(dataRSC, parityRSC)
	if err != nil {
		return err
	}

	N := len(fcontent)
	data := make([][]byte, dataRSC+parityRSC)
	shardLen := int(math.Ceil(float64(N) / float64(dataRSC)))
	fcontentSlice := fcontent

	for i := 0; i < dataRSC-1; i++ {
		data[i] = fcontentSlice[:shardLen]
		fcontentSlice = fcontentSlice[shardLen:]
	}

	lastSlice := make([]byte, shardLen)
	copy(lastSlice, fcontentSlice)
	data[dataRSC-1] = lastSlice
	for i := 0; i < parityRSC; i++ {
		data[dataRSC+i] = make([]byte, shardLen)
	}

	err = enc.Encode(data)
	if err != nil {
		return err
	}

	for i, s := range data {
		err := uploadFile(ringIP, getShardName(fname, i), ringsz, s, certificate)
		if err != nil {
			return err
		}
	}

	return nil
}

// DownloadFileRSC downloads file using Reed Solomon Codes
func DownloadFileRSC(ringIP string, fname string, ringsz uint64, fcontent []byte, certificate string) (int, error) {
	enc, _ := reedsolomon.New(dataRSC, parityRSC)

	shards := make([][]byte, dataRSC+parityRSC)
	var shardlen int
	maxshardlen := int(math.Floor(float64(len(fcontent)) / float64(dataRSC)))
	firstShard := make([]byte, maxshardlen)
	shardnum := 0
	nilshards := make(map[int]bool)

	for {
		if shardnum > parityRSC {
			return 0, fmt.Errorf("Too many corrupt files, can't recover")
		}

		empty, err := downloadFile(ringIP, getShardName(fname, shardnum), ringsz, firstShard, certificate)
		if !os.IsNotExist(err) {
			shardlen = maxshardlen - empty
			if (shardnum+1)*shardlen > len(fcontent) {
				return 0, fmt.Errorf("Not enough space in buffer")
			}

			copy(fcontent[shardnum*shardlen:], firstShard)
			//shardnum++
			break
		}

		shards[shardnum] = nil
		nilshards[shardnum] = true
		shardnum++
	}

	totalEmpty := 0

	for ; shardnum < dataRSC+parityRSC; shardnum++ {
		if shardnum < dataRSC {
			shards[shardnum] = fcontent[shardlen*shardnum : shardlen*(shardnum+1)]
		} else {
			shards[shardnum] = make([]byte, shardlen)
		}

		empty, err := downloadFile(ringIP, fmt.Sprintf("%s_rep%d", fname, shardnum), ringsz, shards[shardnum], certificate)
		if os.IsNotExist(err) {
			shards[shardnum] = nil
			nilshards[shardnum] = true
			continue
		}

		if err != nil {
			return 0, err
		}

		totalEmpty += empty
	}

	err := enc.ReconstructData(shards)
	if err != nil {
		return 0, err
	}

	for n := range nilshards {
		copy(fcontent[shardlen*n:shardlen*(n+1)], shards[n])
	}

	return totalEmpty, nil
}

func DeleteFileRSC(ringIP string, fname string, ringsz uint64, certificate string) error {
	for i := 0; i < dataRSC+parityRSC; i++ {
		err := deleteFile(ringIP, getShardName(fname, i), ringsz, certificate)
		if err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}
