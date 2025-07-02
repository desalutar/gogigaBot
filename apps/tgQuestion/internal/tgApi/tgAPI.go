package controller

import (
	"os"
	"context"
	"strings"
	"gptBot/pkg/logger"
	"google.golang.org/grpc"
	pb "gptBot/pkg/gen/tgHandlers"
	"google.golang.org/grpc/credentials/insecure"
	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	bot    *tgbot.BotAPI
	client pb.QAServiceClient
	log    logger.Logger
	ctx    context.Context
	conn   *grpc.ClientConn
}

func NewTelegramBot(ctx context.Context, log logger.Logger) (*TelegramBot, error) {
	botAPI, err := tgbot.NewBotAPI(os.Getenv("TG_TOKEN"))
	if err != nil {
		return nil, err
	}
	botAPI.Debug = false

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Не удалось подключиться к gRPC", logger.Field{Key: "error", Value: err})
		return &TelegramBot{}, err
	}

	client := pb.NewQAServiceClient(conn)

	return &TelegramBot{
		bot:    botAPI,
		client: client,
		log:    log,
		ctx:    ctx,
		conn: 	conn,
	}, nil
}

func (t *TelegramBot) Run() {
	updateConfig := tgbot.NewUpdate(0)
	updates := t.bot.GetUpdatesChan(updateConfig)

	for {
		select {
		case <-t.ctx.Done():
			t.log.Info("Shutdown signal received, stopping Telegram bot")
			return
		case update := <-updates:
			if update.Message == nil {
				continue
			}
			t.handleUpdate(update)
		}
	}
}

func (t *TelegramBot) handleUpdate(update tgbot.Update) {
	if update.Message == nil || update.Message.Text == "" {
		return
	}

	text := update.Message.Text

	// Проверяем, начинается ли сообщение с @botUsername (например, "@my_bot ")
	if strings.HasPrefix(text, "@"+t.bot.Self.UserName) {
		// Убираем упоминание из текста, чтобы получить вопрос
		question := strings.TrimSpace(strings.TrimPrefix(text, "@"+t.bot.Self.UserName))

		if question == "" {
			// Если вопрос пустой, можно проигнорировать или ответить "Задайте вопрос"
			return
		}

		t.handleAsk(update, question)
		return
	}

	// Если это команда /ask - если хотите оставить поддержку команд
	if update.Message.IsCommand() && update.Message.Command() == "ask" {
		question := strings.TrimSpace(update.Message.CommandArguments())
		t.handleAsk(update, question)
		return
	}

	// В остальных случаях просто эхо
	t.echoMessage(update)
}

// func (t *TelegramBot) handleAskCommand(update tgbot.Update) {
// 	text := strings.TrimSpace(update.Message.CommandArguments())
// 	resp, err := t.client.Ask(context.Background(), &pb.AskRequest{Question: text})
// 	if err != nil {
// 		t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Ошибка gRPC"))
// 		return
// 	}
// 	t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, resp.Answer))
// }

func (t *TelegramBot) handleAsk(update tgbot.Update, question string) {
	resp, err := t.client.Ask(context.Background(), &pb.AskRequest{Question: question})
	if err != nil {
		t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Ошибка gRPC"))
		return
	}

	msg := tgbot.NewMessage(update.Message.Chat.ID, resp.Answer)
	msg.ReplyToMessageID = update.Message.MessageID
	t.bot.Send(msg)
}

func (t *TelegramBot) echoMessage(update tgbot.Update) {
	msg := tgbot.NewMessage(update.Message.Chat.ID, update.Message.Text)
	msg.ReplyToMessageID = update.Message.MessageID
	t.bot.Send(msg)
}

func (t *TelegramBot) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}