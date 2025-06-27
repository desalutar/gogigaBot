package server

import (
	"gptBot/apps/handlers/internal/controller"
	pb "gptBot/pkg/gen/handlers"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("error started: %v", err)
	}
	serverController := controller.NewController()
	s := grpc.NewServer()
	pb.RegisterQAServiceServer(s, serverController)
	
	reflection.Register(s)

	log.Println("Server started port :50051")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Error server %v", err)
	}
}