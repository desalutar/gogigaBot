package main

import (
	// "fmt"
	"gptBot/apps/tgHandlers/internal/server"
	"gptBot/apps/tgHandlers/internal/controller"
)

func main() {
	go controller.StartTelegramBot()
	server.StartServer()

}