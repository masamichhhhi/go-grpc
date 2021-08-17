package main

import (
	"context"
	"fmt"
	pb "grpc-sample/pb/upload"
	"io"
	"log"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func request(client pb.UploadClient) error {
	stream, err := client.Upload(context.Background())
	if err != nil {
		return err
	}

	values := []int32{1, 2, 3, 4, 5}
	for _, value := range values {
		fmt.Println("send:", value)
		if err := stream.Send(&pb.UploadRequest{
			Value: value,
		}); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		time.Sleep(time.Second * 1)
	}

	reply, err := stream.CloseAndRecv()
	if err != nil {
		return err
	}

	log.Printf("result: %v", reply)
	return nil
}

func exec() error {
	address := "localhost:50051"
	conn, err := grpc.Dial(address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return errors.Wrap(err, "connection error")
	}
	defer conn.Close()
	client := pb.NewUploadClient(conn)
	return request(client)
}

func main() {
	if err := exec(); err != nil {
		log.Println(err)
	}
}
