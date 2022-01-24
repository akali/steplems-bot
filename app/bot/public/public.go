package public

import tbot "github.com/go-telegram-bot-api/telegram-bot-api"

// BotApiRepo handles message to send. Used for replies etc.
type BotApiRepo interface {
	// SendMessage sends message
	SendMessage(message tbot.Chattable) (tbot.Message, error)
	// DeleteMessage deletes message by chatID and messageID
	DeleteMessage(chatID int64, messageID int) error
}
