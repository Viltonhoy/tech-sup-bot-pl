package main

import (

	//"github.com/joho/godotenv"

	"log"
	clientchanel "tech-sup-bot-pl/internal/client_chanel"
	"tech-sup-bot-pl/internal/telegram_bot"

	_ "github.com/lib/pq"
	"go.uber.org/zap"
	//"tg-bot-for-ts/repository"
	//"tg-bot-for-ts/service"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("zap.NewDevelopment: %v", err)
	}
	defer logger.Sync()

	bot, err := telegram_bot.NewBot(logger)
	if err != nil {
		logger.Fatal("failed to create a new BotAPI instance", zap.Error(err))
	}

	cc := clientchanel.New()

	err = bot.BotWorker(cc)
	if err != nil {
		// logger.Fatal("")
	}
}
