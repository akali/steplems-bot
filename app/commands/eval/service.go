package eval

import (
	bot "github.com/akali/steplems-bot/app/bot"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
)

// CommandCallback is the callback of the "eval" command.
func CommandCallback(botAPI *tbot.BotAPI, msg tbot.Update) error {
	if rand.Intn(100) == 0 {
		text := "вурвур гей"

		res := bot.NewMessageReply(msg.Message.Chat.ID, text, msg.Message.MessageID)

		_, err := botAPI.Send(res)
		return err
	}
	return nil
}
