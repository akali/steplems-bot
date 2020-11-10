package help

import tbot "github.com/go-telegram-bot-api/telegram-bot-api"

// CommandCallback is the callback of the "help" command.
func CommandCallback(bot *tbot.BotAPI, msg tbot.Update) error {
	text := "This bot is records all messages in group for analysis and cool stuff: \n" +
		"1. /help - show this message\n" +
		"2. /{http_method} url - start to build new request\n"

	res := tbot.NewMessage(msg.Message.Chat.ID, text)

	_, err := bot.Send(res)
	return err
}
