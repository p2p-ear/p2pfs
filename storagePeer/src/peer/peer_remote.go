//Remote calls to the peer grpc server
package peer

import (
	"bufio"
	"context"
	"io"
	"log"
	"os"
)

// Ping generates response to a Ping request
func (p *Peer) Ping(ctx context.Context, in *PingMessage) (*PingMessage, error) {
	log.Printf("Receive message %t", in.Ok)
	return &PingMessage{Ok: true}, nil
}

// FindSuccessorInRing finds id's successor in p's ring
func (p *Peer) FindSuccessorInRing(ctx context.Context, r *FindSuccRequest) (*FindSuccReply, error) {
	ip, err := p.ring.FindSuccessor(r.Id)
	return &FindSuccReply{Ip: ip}, err
}

// Read reads the content of a specified file
func (p *Peer) Read(r *ReadRequest, stream PeerService_ReadServer) error {

	f, err := os.Open(r.Name)
	if os.IsNotExist(err) {
		return stream.Send(&ReadReply{Exists: false})
	}
	if err != nil {
		return err
	}
	defer f.Close()

	if err = ValidateFile(r.Name, r.Certificate, READACT); err != nil {
		return err
	}

	if err = stream.Send(&ReadReply{Exists: true}); err != nil {
		return err
	}

	reader := bufio.NewReader(f)
	b := make([]byte, r.ChunkSize)

	for {
		n, readErr := reader.Read(b)

		if readErr == io.EOF {
			break
		}

		if readErr != nil {
			return err
		}

		if err := stream.Send(&ReadReply{Data: b, Size: int64(n)}); err != nil {
			return err
		}
	}

	return nil
}

// Write writes the content of the request r onto the disk
func (p *Peer) Write(stream PeerService_WriteServer) error {

	writeInfo, err := stream.Recv()

	if err != nil {
		return err
	}

	f, err := os.Create(writeInfo.Name)
	defer f.Close()

	if err != nil {
		return err
	}

	if err = ValidateFile(writeInfo.Name, writeInfo.Certificate, WRITACT); err != nil {
		return err
	}

	writer := bufio.NewWriter(f)

	n, err := writer.Write(writeInfo.Data)

	if err != nil {
		return err
	}

	written := int64(n)

	for {
		toWrite, readErr := stream.Recv()

		if readErr == io.EOF {
			if err = writer.Flush(); err != nil {
				return err
			}

			return stream.SendAndClose(&WriteReply{Written: int64(written)})
		}

		if readErr != nil {
			return readErr
		}

		n, err := writer.Write(toWrite.Data)

		if err != nil {
			return err
		}

		written += int64(n)
	}
}

func (p *Peer) Delete(ctx context.Context, r *DeleteRequest) (*DeleteReply, error) {

	err := ValidateFile(r.Fname, r.Certificate, DELEACT)
	if err != nil {
		return &DeleteReply{}, err
	}

	err = os.Remove(r.Fname)
	if os.IsNotExist(err) {
		return &DeleteReply{Exists: false}, nil
	}
	if err != nil {
		return &DeleteReply{}, err
	}

	return &DeleteReply{Exists: true}, err
}
