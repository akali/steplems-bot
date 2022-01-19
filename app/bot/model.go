package bot

import (
	"github.com/akali/steplems-bot/app/commands"
	"github.com/akali/steplems-bot/app/database"
	"github.com/akali/steplems-bot/app/logger"
	"github.com/akali/steplems-bot/app/youtube"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	// Bot is a wrapper for tbot.BotAPI that stricts and simplifies
	// its functionality.
	Bot struct {
		RunBotRepo
		RecordMessageRepo
		api      *tbot.BotAPI
		commands commands.CallbackMap
		Database *database.Database
		Youtube  *youtube.Youtube
	}
)

var (
	log = logger.Factory.NewLogger("bot")
)

// NewBot initializes bot api and returns a new *Bot.
func NewBot(token string, commands commands.CallbackMap, url string, databaseName string) (*Bot, error) {
	api, err := tbot.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	databaseConfig := &database.Database{
		Url:      url,
		Database: databaseName,
	}

	return &Bot{api: api, commands: commands, Database: databaseConfig, Youtube: youtube.NewYoutube()}, nil
}
