package main

import (
	"context"
	pb "gogigaBot/pkg/gen/handlers"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct {
	pb.UnimplementedQAServiceServer
}

func (s *server) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	log.Printf("Получен новый запрос: %v", req)
	return &pb.AskResponse{Answer: "Привет от gRPC сервера!"}, nil	
}

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error started: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterQAServiceServer(s, &server{})
	
	reflection.Register(s)

	log.Println("Server started port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error server %v", err)
	}
}