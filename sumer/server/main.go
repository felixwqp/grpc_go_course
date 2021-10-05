package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/felixwqp/grpc_go_course/sumer/sumpb"
	"google.golang.org/grpc"
)

type server struct {
	sumpb.UnimplementedSumServiceServer
}

func (*server) Sum(ctx context.Context, req *sumpb.Arguments) (*sumpb.Result, error) {
	A := req.GetArgumentA()
	B := req.GetArgumentB()
	result := A + B
	res := &sumpb.Result{
		Sum: result,
	}
	return res, nil
}

func (*server) BiDiSum(stream sumpb.SumService_BiDiSumServer) error {
	fmt.Println("Start BiDi Stream Sum")
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Cannot recv stream from client: %v", err)
			return err
		}
		argA := req.GetArgumentA()
		argB := req.GetArgumentB()
		SendErr := stream.Send(&sumpb.Result{
			Sum: argA + argB,
		})
		if SendErr != nil {
			log.Fatalf("Cannot send stream from client: %v", err)
			return SendErr
		}
	}
}

func main() {
	fmt.Println("Hello World")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	sumpb.RegisterSumServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v")
	}
	// s.Sum

}
