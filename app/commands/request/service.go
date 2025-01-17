package request

import (
	"fmt"
	"github.com/akali/steplems-bot/app/bot"
	"github.com/akali/steplems-bot/app/commands/request/arguments"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
	"net/http"
	"strings"
)

// PerformRequest creates and sends a new http request configured by the set
// of arguments in the message.
func PerformRequest(method string, botAPI *tbot.BotAPI, msg tbot.Update) error {
	commandArgs := msg.Message.CommandArguments()
	splitArgs := strings.Split(commandArgs, " ")

	log.Debug.PrintTKV("request command - {{method}}, args: {{args}}",
		"method", method,
		"args", splitArgs,
	)

	args, err := arguments.ParseArguments(splitArgs...)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(method, args.URL.String(), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	log.Debug.Println(*resp)

	text := fmt.Sprintf(
		"Status: %v",
		resp.Status,
	)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		text += "\n" + string(bodyBytes)
	}

	message := bot.NewMessageReply(msg.Message.Chat.ID, text, msg.Message.MessageID)

	_, err = botAPI.Send(message)

	return err
}
