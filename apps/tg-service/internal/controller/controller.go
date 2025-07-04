package controller

import (
	"context"
	gpt "gptBot/pkg/gen/gpt-service"
	pb "gptBot/pkg/gen/tg-service"
	"gptBot/pkg/logger"
)

type Handler struct {
	pb.UnimplementedQAServiceServer
	gptClient gpt.QAServiceClient
	logger logger.Logger
}

type Handlerer interface {
	Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error)
}

func NewController(ctx context.Context, gptClient gpt.QAServiceClient, logger logger.Logger) *Handler {
	return &Handler{gptClient: gptClient, logger: logger}
}

func (h *Handler) Ask(ctx context.Context, req *pb.AskRequest) (*pb.AskResponse, error) {
	h.logger.Info("Получен новый запрос", logger.Field{Key: "request", Value: req.Question})

	resp, err := h.gptClient.Ask(ctx, &gpt.AskRequest{
		Question: req.Question,
	})
	if err != nil {
		h.logger.Error("Ошибка при вызове gpt-сервиса", logger.Field{Key: "error", Value: err})
		return nil, err
	}

	return &pb.AskResponse{Answer: resp.Answer}, nil	
}