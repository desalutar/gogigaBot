package controller

import (
	"R2D2/apps/gptAnswer/internal/gptApi"
	"context"
	pb "gptBot/pkg/gen/gpt-service"
	"gptBot/pkg/logger"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
	client *gptApi.OpenAIClient
	logger logger.Logger
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController(ctx context.Context, gptAPi *gptApi.OpenAIClient, logger logger.Logger) *Handler {
	return &Handler{client: gptAPi, logger: logger}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	h.logger.Info("Получен новый запрос", logger.Field{Key: "request", Value: req.Question})

	answer, err := h.client.GetCompletion(ctx, req.Question)
	if err != nil {
		h.logger.Error("Ошибка при вызове OpenAI", logger.Field{Key: "error", Value: err})
		return nil, err
	}

	return &pb.AskResponse{Answer: answer}, nil	
}