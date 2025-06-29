package server

import (
	"gptBot/apps/tgHandlers/internal/controller"
	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer(log logger.Logger) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("error started: %v", logger.Field{Key: "error", Value: err})
	}
	serverController := controller.NewController(log)
	s := grpc.NewServer()
	pb.RegisterQAServiceServer(s, serverController)
	
	reflection.Register(s)

	log.Info("Server started", logger.Field{Key: "port", Value: 50051})
	if err := s.Serve(lis); err != nil {
		log.Error("Error server %v", logger.Field{Key: "error", Value: err})
	}
}