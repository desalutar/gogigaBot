package controller

import (
	"context"
	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"
	"log"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
	log logger.Logger
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController(ctx context.Context, log logger.Logger) *Handler {
	return &Handler{log: log}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	log.Printf("Получен новый запрос: %v", req)
	return &pb.AskResponse{Answer: "Привет от gRPC сервера!"}, nil	
}