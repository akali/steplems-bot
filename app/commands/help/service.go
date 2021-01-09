package help

import (
	"github.com/akali/steplems-bot/app/bot"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

// CommandCallback is the callback of the "help" command.
func CommandCallback(botAPI *tbot.BotAPI, msg tbot.Update) error {
	text := "This botAPI is records all messages in group for analysis and cool stuff: \n" +
		"1. /help - show this message\n" +
		"2. /{http_method} url - start to build new request\n"

	res := bot.NewMessageReply(msg.Message.Chat.ID, text, msg.Message.MessageID)

	_, err := botAPI.Send(res)
	return err
}
