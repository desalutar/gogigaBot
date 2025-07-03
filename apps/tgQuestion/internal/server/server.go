package server

import (
	"context"
	"gptBot/apps/tgQuestion/internal/controller"
	"gptBot/pkg/gen/gpt"
	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func StartServer(ctx context.Context, log logger.Logger) {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Error("error started: %v", logger.Field{Key: "error", Value: err})
	}

	conn, err := grpc.NewClient("gpt:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("не удалось подключиться к gpt-сервису: %v", logger.Field{Key: "Error", Value: err})
	}
	defer conn.Close()

	gptClient :=gpt.NewQAServiceClient(conn)

	serverController := controller.NewController(ctx, gptClient, log)
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