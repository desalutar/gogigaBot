package server

import (
	"R2D2/apps/gptAnswer/internal/controller"
	"R2D2/apps/gptAnswer/internal/gptApi"
	"context"
	pb "gptBot/pkg/gen/gpt"
	"gptBot/pkg/logger"
	"net"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// TODO: ДОПИСАТЬ GATEWAY
func StartServer(ctx context.Context, log logger.Logger) {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Error("error started: %v", logger.Field{Key: "error", Value: err})
	}

	client := gptApi.NewClient()

	serverController := controller.NewController(ctx, client, log)
	
	s := grpc.NewServer()
	pb.RegisterQAServiceServer(s, serverController)
	reflection.Register(s)

	go func() {
		<-ctx.Done()
		log.Info("Shutdown signal received, stopping gRPC server")
		s.GracefulStop()
	}()

	log.Info("Server started", logger.Field{Key: "port", Value: 50052})
	if err := s.Serve(lis); err != nil {
		log.Error("Error server %v", logger.Field{Key: "error", Value: err})
	}
}