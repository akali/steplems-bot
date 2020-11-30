package bot

import (
	"fmt"

	"github.com/akali/steplems-bot/app/commands"
	"github.com/akali/steplems-bot/app/database"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Run starts listening to bot api and waits for new messages.
func (b *Bot) Run(timeout int) error {
	defer log.Note.Println("bot stopped running")

	log.Info.Println("update timout:", timeout)

	u := tbot.NewUpdate(0)
	u.Timeout = timeout

	updates, err := b.api.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	log.Succ.Println("successfully started update chan listener")

	for update := range updates {
		b.update(update)
	}

	return nil
}

func (b *Bot) Record(message *tbot.Message) error {
	log.Info.PrintTKV("Recording message from {{from}}", "from", func() string {
		if len(message.From.UserName) > 0 {
			return message.From.UserName
		} else {
			return message.From.FirstName + " " + message.From.LastName
		}
	}())

	return b.Database.SaveMessage(&database.Message{
		Message: *message,
	})
}

func (b *Bot) update(update tbot.Update) {
	// Ignore any non-Message Updates.
	if update.Message == nil {
		return
	}

	log.Info.PrintTKV("[{{update_id}}] {{username}}: {{text}}",
		"update_id", update.UpdateID,
		"username", update.Message.From.UserName,
		"text", update.Message.Text,
	)

	// Record the message
	if err := b.Record(update.Message); err != nil {
		log.Error.Println("unexpected error!", err.Error())
	}

	if name := update.Message.CommandWithAt(); len(name) > 0 {
		commandName := commands.Name(name)

		// Sending help message if the command by the given name wasn't found.
		if callback, ok := b.commands.Get(commandName); !ok {
			log.Warn.PrintT("command '{}' not found", commandName)
			return
		} else {
			b.executeCommand(update, callback, commandName)
		}
	}
}

func (b *Bot) executeCommand(update tbot.Update, callback commands.Callback, name commands.Name) {
	defer func() {
		if err := recover(); err != nil {
			log.Error.Println("unexpected error!", err)
			b.sendErrorMessage(update.Message.Chat.ID, "Sorry, we have some problems here")
		}
	}()

	err := callback(b.api, update)
	if err != nil {
		log.Error.PrintT("error while executing a command '{}'", name)
		b.sendErrorMessage(update.Message.Chat.ID, b.formError(err))
		return
	}
}

func (b *Bot) formError(err error) string {
	return fmt.Sprintf("Error trying to execute your command: %s", err.Error())
}

func (b *Bot) sendErrorMessage(chatID int64, err string) {
	errMsg := tbot.NewMessage(chatID, err)
	_, sendError := b.api.Send(errMsg)
	if sendError != nil {
		log.Error.Println("error trying to send an error message:", err)
	}
}
