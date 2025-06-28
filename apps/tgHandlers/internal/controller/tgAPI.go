package controller

import (
	"context"
	pb "gptBot/pkg/gen/tgHandlers"
	"log"
	"os"
	"strings"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Файл .env не найден, продолжаем без него")
	}
}

func StartTelegramBot() {
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
		log.Fatalf("Не удалось подключиться к gRPC: %v", err)
	}
	defer conn.Close()

	client := pb.NewQAServiceClient(conn)

	for update := range updates {
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
			// Просто повторяем сообщение
			msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}
}
