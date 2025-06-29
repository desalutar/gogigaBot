package controller

import (
	"context"
	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"
	"os"
	"strings"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func StartTelegramBot(ctx context.Context, log logger.Logger) {
	bot, err := tgbot.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot.Debug = true

	updateConfig := tgbot.NewUpdate(0)
	updateConfig.Timeout = 30
	updates := bot.GetUpdatesChan(updateConfig)

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Не удалось подключиться к gRPC", logger.Field{Key: "error", Value: err})
		return
	}
	defer conn.Close()

	client := pb.NewQAServiceClient(conn)

	for {
		select {
		case <-ctx.Done():
			log.Info("Shutdown signal received, stopping Telegram bot")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}

			if update.Message.IsCommand() && update.Message.Command() == "ask" {
				text := strings.TrimSpace(update.Message.CommandArguments())
				resp, err := client.Ask(context.Background(), &pb.AskRequest{Question: text})
				if err != nil {
					bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Ошибка gRPC"))
					continue
				}
				bot.Send(tgbot.NewMessage(update.Message.Chat.ID, resp.Answer))
			} else {
				msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
				msg.ReplyToMessageID = update.Message.MessageID
				bot.Send(msg)
			}
		}
	}
}
