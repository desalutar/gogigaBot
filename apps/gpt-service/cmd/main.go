package main

import (
	// "R2D2/apps/gptAnswer/controller"
	"R2D2/apps/gptAnswer/internal/server"
	"context"
	"gptBot/pkg/config"
	"gptBot/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := SetupShutdownContext()
	defer cancel()

	if err := godotenv.Load(); err != nil {
		log.Fatal("Файл .env не найден, продолжаем без него")
	}

	conf := config.LoadConfig()

	logger, err := logger.NewLogger(conf.Logger)
	if err != nil {
		log.Fatalf("error create logger")
	}

	server.StartServer(ctx, logger)

	<-ctx.Done()
	logger.Info("Shutdown signal received, exiting")
}


func SetupShutdownContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-signals
		cancel()
	}()

	return ctx, cancel
}