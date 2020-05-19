// IP-targeted file management
package peer

import (
	"context"
	"fmt"
	"io"
	"math"
	"os"
	"time"
)

const chunksz = 8

// sendFile sends file to the target IP
func sendFile(targetIP string, fname string, fcontent []byte, certificate string) error {

	conn, cl, err := Connect(targetIP)
	if err != nil {
		return err
	}
	defer conn.Close()

	fmt.Printf("Opening write stream to %s...", targetIP)
	// Stream to write
	wstream, err := cl.Write(context.Background())
	for err != nil {
		fmt.Println(err.Error())
		fmt.Println("Couldn't initialize remote write stream")
		time.Sleep(time.Second * 1)

		wstream, err = cl.Write(context.Background())
	}

	fmt.Printf("Sending filename %s...", fname)
	// Send request
	err = wstream.Send(&WriteRequest{Name: fname, Certificate: certificate})

	for err != nil {
		fmt.Println(err.Error())
		fmt.Println("Couldn't send filename")
		time.Sleep(time.Second * 1)
		err = wstream.Send(&WriteRequest{Name: fname})
	}

	chunkSize := chunksz
	chunkAmnt := int(math.Ceil(float64(len(fcontent)) / float64(chunkSize)))

	fmt.Println("Writing to file, total chunks:", chunkAmnt)
	for i := 0; i < chunkAmnt; i++ {

		curChunk := fcontent[i*chunkSize:]
		if len(curChunk) > chunkSize {
			fmt.Println()
			curChunk = curChunk[:chunkSize]
		}

		err := wstream.Send(&WriteRequest{Data: curChunk})

		for err != nil {
			fmt.Println(err.Error())
			fmt.Printf("Unable to send chunk #%d", i)
			time.Sleep(time.Second * 1)
			err = wstream.Send(&WriteRequest{Name: fname})
		}
	}

	reply, err := wstream.CloseAndRecv()
	fmt.Printf("Finished writing %d bytes", reply.Written)
	return err
}

// recvFile recieves file w/ filename=fname, from node targetIP - returns how much empty space is at the end (negative, if buffer is too small)
func recvFile(targetIP string, fname string, fcontent []byte, certificate string) (int, error) {

	conn, peer, err := Connect(targetIP)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	rstream, err := peer.Read(context.Background(), &ReadRequest{Name: fname, ChunkSize: chunksz, Certificate: certificate})
	if err != nil {
		return 0, err
	}

	readReply, err := rstream.Recv()
	if !readReply.Exists {
		return 0, os.ErrNotExist
	}

	contentSlice := fcontent[:]
	bufferSmall := false
	emptySpace := len(fcontent)

	for {
		readReply, err := rstream.Recv()

		if err == io.EOF {
			break
		}

		if err != nil {
			return 0, err
		}

		if !bufferSmall && int(readReply.Size) > emptySpace {
			for i := range contentSlice {
				contentSlice[i] = readReply.Data[i]
			}

			emptySpace = 0
			bufferSmall = true
		}

		if !bufferSmall {
			for i := 0; i < int(readReply.Size); i++ {
				contentSlice[i] = readReply.Data[i]
			}
			contentSlice = contentSlice[readReply.Size:]
		}

		emptySpace -= int(readReply.Size)
	}

	return emptySpace, nil
}

func remvFile(targetIP string, fname string, certificate string) error {
	conn, peer, err := Connect(targetIP)
	if err != nil {
		return err
	}
	defer conn.Close()

	r, err := peer.Delete(context.Background(), &DeleteRequest{Fname: fname, Certificate: certificate})
	if err == nil && !r.Exists {
		err = os.ErrNotExist
	}

	return err
}
