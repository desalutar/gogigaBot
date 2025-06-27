package controller

import (
	"context"
	pb "gogigaBot/pkg/gen/handlers"
	"log"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController() *Handler {
	return &Handler{}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	log.Printf("Получен новый запрос: %v", req)
	return &pb.AskResponse{Answer: "Привет от gRPC сервера!"}, nil	
}