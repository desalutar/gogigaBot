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

	if !t.isBotMentioned(update.Message) {
		return
	}

	question := strings.Replace(update.Message.Text, "@"+t.bot.Self.UserName, "", 1)
	question = strings.TrimSpace(question)
	if question == "" {
		return
	}

	t.handleAsk(update, question)
}


func (t *TelegramBot) isBotMentioned(message *tgbot.Message) bool {
	for _, entity := range message.Entities {
		if entity.Type == "mention" {
			mention := message.Text[entity.Offset : entity.Offset+entity.Length]
			if mention == "@"+t.bot.Self.UserName {
				return true
			}
		}
	}
	return false
}

func (t *TelegramBot) handleAsk(update tgbot.Update, question string) {
	resp, err := t.client.Ask(t.ctx, &pb.AskRequest{Question: question})
	if err != nil {
        t.log.Error("gRPC Ask failed", logger.Field{Key: "error", Value: err})
        t.bot.Send(tgbot.NewMessage(update.Message.Chat.ID, "Ошибка gRPC: "+err.Error()))
        return
    }

	msg := tgbot.NewMessage(update.Message.Chat.ID, resp.Answer)
	msg.ReplyToMessageID = update.Message.MessageID
	t.bot.Send(msg)
}

func (t *TelegramBot) Close() {
	if t.conn != nil {
		t.conn.Close()
	}
}