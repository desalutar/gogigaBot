package main

import (
	"gptBot/apps/tgHandlers/internal/controller"
	"gptBot/apps/tgHandlers/internal/server"
	"gptBot/pkg/config"
	"gptBot/pkg/logger"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Файл .env не найден, продолжаем без него")
	}

	conf := config.LoadConfig()

	logger, err := logger.NewLogger(conf.Logger)
	if err != nil {
		log.Fatalf("error create logger")
	}

	go controller.StartTelegramBot(logger)
	server.StartServer(logger)
}