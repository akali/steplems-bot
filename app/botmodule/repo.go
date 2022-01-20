package botmodule

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Module interface {
	MessageUpdate(message *tbot.Message) error
}
