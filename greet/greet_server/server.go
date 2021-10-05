package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/felixwqp/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
)

type server struct {
	greetpb.UnimplementedGreetServiceServer
}

func (*server) Greet(ctx context.Context, req *greetpb.GreetRequest) (*greetpb.GreetResponse, error) {
	firstName := req.GetGreeting().GetFirstName()
	result := "Hello" + firstName
	res := &greetpb.GreetResponse{
		Result: result,
	}
	return res, nil
}

func (*server) GreetManyTimes(req *greetpb.GreetManyTimesRequest, stream greetpb.GreetService_GreetManyTimesServer) error {
	firstName := req.Greeting.GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + firstName + ", Number " + strconv.Itoa(i)
		res := &greetpb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

func (*server) LongGreet(stream greetpb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet Function invoked.")
	result := "Hello "
	for {
		msg, eil := stream.Recv()
		if eil == io.EOF {
			// we have finished reading client streaming
			break
		}
		if eil != nil {
			log.Fatalf("Error while reading client stream: %v", eil)
		}
		firstName := msg.GetGreeting().GetFirstName()
		result += "Hello " + firstName + "! "

	}

	return stream.SendAndClose(&greetpb.LongGreetResponse{
		Result: result,
	})
}

func main() {
	fmt.Println("Hello World")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	greetpb.RegisterGreetServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v")
	}

}
