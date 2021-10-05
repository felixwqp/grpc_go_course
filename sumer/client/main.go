package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/felixwqp/grpc_go_course/sumer/sumpb"
	"google.golang.org/grpc"
)

func Sum(c sumpb.SumServiceClient) {
	fmt.Println("Starting to Sum RPC")
	req := &sumpb.Arguments{
		ArgumentA: 10,
		ArgumentB: 12,
	}
	res, err := c.Sum(context.Background(), req)
	if err != nil {
		fmt.Println("Cannot get response from server")
	}
	fmt.Printf("Get Response Result: %d", res.GetSum())

}

func BiDiSum(c sumpb.SumServiceClient) {
	fmt.Println("Starting to BiDi Sum RPC")
	stream, err := c.BiDiSum(context.Background())
	if err != nil {
		log.Fatalf("Cannot create the stream: %v", err)
	}
	block := make(chan struct{})
	// request
	go func() {
		primes := [6]int32{2, 3, 5, 7, 11, 13}
		for idx, v := range primes {
			SendErr := stream.Send(&sumpb.Arguments{
				ArgumentA: int32(idx),
				ArgumentB: v,
			})
			fmt.Printf("Request: %d\n", v)
			if SendErr != nil {
				log.Fatalf("Cannot send: %v", SendErr)
			}
			time.Sleep(1000 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	// receive
	go func() {
		for {
			recv, err := stream.Recv()
			if err == io.EOF {
				break
			}
			fmt.Printf("Recv: %d\n", recv.GetSum())
		}
		close(block)
	}()
	<-block
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
	c := sumpb.NewSumServiceClient(cc)
	fmt.Printf("Created Client: %s", c)

	BiDiSum(c)
}
