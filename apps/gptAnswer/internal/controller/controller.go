package controller

import (
	"R2D2/apps/gptAnswer/internal/gptApi"
	"context"
	pb "gptBot/pkg/gen/gpt"
	"log"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
	client *gptApi.OpenAIClient
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController(ctx context.Context, gptAPi *gptApi.OpenAIClient) *Handler {
	return &Handler{client: gptAPi}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	log.Printf("Получен новый запрос: %v", req)

	answer, err := h.client.GetCompletion(ctx, req.Question)
	if err != nil {
		log.Printf("Ошибка при вызове OpenAI: %v", err)
		return nil, err
	}

	return &pb.AskResponse{Answer: answer}, nil	
}