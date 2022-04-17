package main

import (
	"context"
	"grpc_example/core"
	"io"
	"log"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial(":5001", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Could not connect to server. %v", err)
	}

	defer conn.Close()

	client := core.NewGreetingServiceClient(conn)
	greet(client)
	greetManyTimes(client)
	longGreet(client)
	greetEveryone(client)
}

func greet(c core.GreetingServiceClient) {
	request := core.GreetingRequest{Greeting: &core.Greeting{FirstName: "FirstName", LastName: "LastName"}}
	res, err := c.Greet(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error greeting. %v", err)
	}

	log.Printf("Greet response: %v", res.Result)
}

func greetManyTimes(c core.GreetingServiceClient) {
	request := core.GreetManyTimesRequest{Greeting: &core.Greeting{FirstName: "FirstName", LastName: "LastName"}}
	stream, err := c.GreetManyTimes(context.Background(), &request)
	if err != nil {
		log.Fatalf("Error greeting many times. %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("Error getting response from stream of greetmanytimes. %v", err)
		}

		log.Printf("GreetManyTimes response: %v", res.Result)
	}
}

func longGreet(c core.GreetingServiceClient) {
	requests := []*core.LongGreetRequest{}
	for i := 1; i < 5; i++ {
		index := strconv.Itoa(i)
		requests = append(requests, &core.LongGreetRequest{Greeting: &core.Greeting{FirstName: "FirstName" + index, LastName: "LastName" + index}})
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("Error long greet. %v", err)
	}

	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error sending long greet request to stream. %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error getting response from stream of longgreet. %v", err)
	}

	log.Printf("LongGreet response: %v", res.Result)
}

func greetEveryone(c core.GreetingServiceClient) {
	stream, err := c.GreetEveryone(context.Background())
	waitCh := make(chan struct{})
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				close(waitCh)
				break
			}
			if err != nil {
				log.Fatalf("Error getting response from stream of greeteveryone. %v", err)
			}

			log.Printf("GreetEveryone response: %v", res.Result)
		}
	}()

	requests := []*core.GreetEveryoneRequest{}
	for i := 1; i < 5; i++ {
		index := strconv.Itoa(i)
		requests = append(requests, &core.GreetEveryoneRequest{Greeting: &core.Greeting{FirstName: "FirstName" + index, LastName: "LastName" + index}})
	}

	if err != nil {
		log.Fatalf("Error greeting everyone. %v", err)
	}

	for _, req := range requests {
		if err := stream.Send(req); err != nil {
			log.Fatalf("Error sending greet everyone request to stream. %v", err)
		}

		time.Sleep(time.Second)
	}

	stream.CloseSend()
	<-waitCh
}
