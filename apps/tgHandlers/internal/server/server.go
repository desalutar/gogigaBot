package server

import (
	"context"
	"gptBot/apps/tgHandlers/internal/controller"
	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func StartServer(ctx context.Context, log logger.Logger) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("error started: %v", logger.Field{Key: "error", Value: err})
	}

	serverController := controller.NewController(ctx)
	s := grpc.NewServer()
	pb.RegisterQAServiceServer(s, serverController)
	reflection.Register(s)

	go func() {
		<-ctx.Done()
		log.Info("Shutdown signal received, stopping gRPC server")
		s.GracefulStop()
	}()

	log.Info("Server started", logger.Field{Key: "port", Value: 50051})
	if err := s.Serve(lis); err != nil {
		log.Error("Error server %v", logger.Field{Key: "error", Value: err})
	}
}