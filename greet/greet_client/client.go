package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/felixwqp/grpc_go_course/greet/greetpb"
	"google.golang.org/grpc"
)

func DoUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Startto do undary RPC")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Felix",
			LastName:  "Wang",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		fmt.Println("Cannot get response from server")
	}
	fmt.Printf("Get Response Result: %s", res.Result)

}
func doServerStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a streaming RPC")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Felix",
			LastName:  "Wang",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error While calling GreetManyTimes RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response: %v", msg.GetResult())

	}

}

func doClientStreaming(c greetpb.GreetServiceClient) {
	fmt.Println("Starting Streaming client RPC")
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error while calling LongGreet: %v", err)
	}
	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Felix",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Come On",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Coding",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Design",
			},
		},
	}
	for _, req := range requests {
		fmt.Printf("Sending Request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}
	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving: %v\n", err)
	}
	fmt.Println(res)

}

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	// production use SSL certificate,
	// now without one, use in-secure
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	fmt.Printf("Created Client: %s", c)
	// DoUnary(c)
	// doServerStreaming(c)
	doClientStreaming(c)

}
