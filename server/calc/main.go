package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"google.golang.org/grpc"

	pb "grpc-sample/pb/calc"
)

const port = ":50051"

type ServerUnary struct {
	pb.UnimplementedCalcServer
}

func (s *ServerUnary) Sum(ctx context.Context, in *pb.SumRequest) (*pb.SumReply, error) {
	a := in.GetA()
	b := in.GetB()
	fmt.Println(a, b)
	reply := fmt.Sprintf("%d + %d = %d", a, b, a+b)
	return &pb.SumReply{
		Message: reply,
	}, nil
}

func set() error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return errors.Wrap(err, "ポート失敗")
	}
	s := grpc.NewServer()
	var server ServerUnary
	pb.RegisterCalcServer(s, &server)
	if err := s.Serve(lis); err != nil {
		return errors.Wrap(err, "サーバー起動失敗")
	}

	return nil
}

func main() {
	fmt.Println("start")
	if err := set(); err != nil {
		log.Fatalf("%v", err)
	}
}
