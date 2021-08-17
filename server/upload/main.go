package main

import (
	"fmt"
	pb "grpc-sample/pb/upload"
	"io"
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

const port = ":50051"

type ServerClientSide struct {
	pb.UnimplementedUploadServer
}

func (s *ServerClientSide) Upload(stream pb.Upload_UploadServer) error {
	var sum int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			message := fmt.Sprintf("DONE: sum = %d", sum)
			return stream.SendAndClose(&pb.UploadReply{
				Message: message,
			})
		}
		if err != nil {
			return err
		}
		fmt.Println(req.GetValue())
		sum += req.GetValue()
	}
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}
	s := grpc.NewServer()
	var server ServerClientSide
	pb.RegisterUploadServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "connection failed")
	}
	return nil
}

func main() {
	fmt.Println("起動")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
