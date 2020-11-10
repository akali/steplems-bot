package bot

import tbot "github.com/go-telegram-bot-api/telegram-bot-api"

type (
	// RunBotRepo handles bots update chan.
	RunBotRepo interface {
		// Run starts listening to bot api and waits for new messages.
		Run(timeout int) error
	}

	// RecordMessageRepo handles messages for record.
	RecordMessageRepo interface {
		// Records the messages
		Record(message *tbot.Update) error
	}
)
