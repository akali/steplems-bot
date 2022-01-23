package bot

import (
	"fmt"

	"github.com/akali/steplems-bot/app/botmodule"
	"github.com/akali/steplems-bot/app/commands"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	multierror "github.com/hashicorp/go-multierror"
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
		go b.update(update)
	}

	return nil
}

func (b *Bot) notifyModules(message *tbot.Message) error {
	var moduleErrs error
	for _, module := range b.Modules {
		if err := module.MessageUpdate(message); err != nil {
			moduleErrs = multierror.Append(moduleErrs, err)
		}
	}
	return moduleErrs
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

	// Run through errors
	if err := b.notifyModules(update.Message); err != nil {
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

func (b *Bot) RegisterModule(module botmodule.Module) error {
	if module == nil {
		return fmt.Errorf("module should not be nil")
	}
	b.Modules = append(b.Modules, module)
	return nil
}

// NewMessageReply creates a new Message with reply.
//
// chatID is where to send it, text is the message text, replyMessageID is to whom reply.
func NewMessageReply(chatID int64, text string, replyMessageID int) tbot.MessageConfig {
	message := tbot.NewMessage(chatID, text)
	message.ReplyToMessageID = replyMessageID
	return message
}

func (b *Bot) SendMessage(message tbot.Chattable) (tbot.Message, error) {
	return b.api.Send(message)
}
