package main

import (
	"context"
	"fmt"
	"grpc_example/core"
	"io"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
)

type GreetingServer struct {
	core.UnimplementedGreetingServiceServer
}

func (s *GreetingServer) Greet(ctx context.Context, request *core.GreetingRequest) (*core.GreetingResponse, error) {
	firstName := request.Greeting.FirstName
	lastName := request.Greeting.LastName

	response := core.GreetingResponse{Result: fmt.Sprintf("Hello %s %s", firstName, lastName)}

	return &response, nil
}

func (s *GreetingServer) GreetManyTimes(request *core.GreetManyTimesRequest, stream core.GreetingService_GreetManyTimesServer) error {
	firstName := request.Greeting.FirstName
	lastName := request.Greeting.LastName

	response := core.GreetManyTimesResponse{Result: fmt.Sprintf("Hello %s %s", firstName, lastName)}

	for i := 0; i < 5; i++ {
		if err := stream.Send(&response); err != nil {
			return err
		}

		time.Sleep(time.Second)
	}

	return nil
}

func (s *GreetingServer) LongGreet(stream core.GreetingService_LongGreetServer) error {
	response := core.LongGreetResponse{}
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&response)
		}

		firstName := request.Greeting.FirstName
		lastName := request.Greeting.LastName

		response.Result += fmt.Sprintf("Hello %s %s\n", firstName, lastName)

		if err != nil {
			return err
		}
	}
}

func (s *GreetingServer) GreetEveryone(stream core.GreetingService_GreetEveryoneServer) error {
	for {
		request, err := stream.Recv()
		if err == io.EOF {
			break
		}

		firstName := request.Greeting.FirstName
		lastName := request.Greeting.LastName
		stream.Send(&core.GreetEveryoneResponse{Result: fmt.Sprintf("Hello %s %s\n", firstName, lastName)})

		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":5001")
	if err != nil {
		log.Fatalf("Failed to listen. %v", err)
	}

	grpcServer := grpc.NewServer()
	core.RegisterGreetingServiceServer(grpcServer, &GreetingServer{})
	log.Printf("Server listening at %v", lis.Addr())
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve. %v", err)
	}
}
