package controller

import (
	"context"
	pb "gptBot/pkg/gen/tgHandlers"
	"log"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController(ctx context.Context) *Handler {
	return &Handler{}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	log.Printf("Получен новый запрос: %v", req)
	return &pb.AskResponse{Answer: "Привет от gRPC сервера!"}, nil	
}