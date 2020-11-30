package eval

import (
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
)

// CommandCallback is the callback of the "eval" command.
func CommandCallback(bot *tbot.BotAPI, msg tbot.Update) error {
	if rand.Intn(100) == 0 {
		text := "вурвур гей"

		res := tbot.NewMessage(msg.Message.Chat.ID, text)

		_, err := bot.Send(res)
		return err
	}
	return nil
}
