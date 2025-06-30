package controller

import (
	"context"
	"os"
	"strings"

	pb "gptBot/pkg/gen/tgHandlers"
	"gptBot/pkg/logger"

	tgbot "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	if update.Message.IsCommand() && update.Message.Command() == "ask" {
		t.handleAskCommand(update)
	} else {
		t.echoMessage(update)
	}
}

func (t *TelegramBot) handleAskCommand(update tgbot.Update) {
	text := strings.TrimSpace(update.Message.CommandArguments())
	resp, err := t.client.Ask(context.Background(), &pb.AskRequest{Question: text})
	if err != nil {
		t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Ошибка gRPC"))
		return
	}
	t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, resp.Answer))
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