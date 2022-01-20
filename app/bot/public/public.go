package public

import tbot "github.com/go-telegram-bot-api/telegram-bot-api"

// SendMessageRepo handles message to send. Used for replies etc.
type SendMessageRepo interface {
	// SendMessage sends message
	SendMessage(message tbot.Chattable) (tbot.Message, error)
}
