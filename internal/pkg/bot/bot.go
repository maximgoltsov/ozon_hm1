package commander

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/pkg/errors"

	"github.com/maximgoltsov/botproject/config"
	commandPkg "github.com/maximgoltsov/botproject/internal/pkg/bot/command"
)

var UnknownCommand = errors.New("unknown command")

type Interface interface {
	Run() error
	RegisterHandler(cmd commandPkg.Interface)
}

func MustNew() Interface {
	bot, err := tgbotapi.NewBotAPI(config.ApiKey)
	if err != nil {
		log.Panic(errors.Wrap(err, "init tgbot"))
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	return &commander{
		bot:   bot,
		route: make(map[string]commandPkg.Interface),
	}
}

type commander struct {
	bot   *tgbotapi.BotAPI
	route map[string]commandPkg.Interface
}

// RegisterHandler - not thread-safe
func (c *commander) RegisterHandler(cmd commandPkg.Interface) {
	c.route[cmd.Name()] = cmd
}

func (c *commander) Run() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		if cmdName := update.Message.Command(); cmdName != "" {
			if cmd, ok := c.route[cmdName]; ok {
				msg.Text = cmd.Process(update.Message.CommandArguments())
			} else {
				msg.Text = "Unknown command"
			}
		} else {
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			msg.Text = fmt.Sprintf("you send <%v>", update.Message.Text)

		}
		_, err := c.bot.Send(msg)

		if err != nil {
			return errors.Wrap(err, "send tg message")
		}

	}
	return nil
}
