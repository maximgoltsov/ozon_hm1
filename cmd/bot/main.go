package main

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/maximgoltsov/botproject/config"
	"github.com/maximgoltsov/botproject/internal/commander"
	"github.com/maximgoltsov/botproject/internal/handlers"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		log.Panic(err)
	}

	commander, err := commander.Init(bot)

	if err != nil {
		log.Panic(err)
	}

	handlers.AddHandlers(commander)

	if err := commander.Run(); err != nil {
		log.Panic(err)
	}
}
