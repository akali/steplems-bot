package bot

import (
	"github.com/akali/steplems-bot/app/bot/public"
	"github.com/akali/steplems-bot/app/botmodule"
	"github.com/akali/steplems-bot/app/botmodule/youtube"
	"github.com/akali/steplems-bot/app/commands"
	"github.com/akali/steplems-bot/app/database"
	"github.com/akali/steplems-bot/app/logger"
	tbot "github.com/go-telegram-bot-api/telegram-bot-api"
)

type (
	// Bot is a wrapper for tbot.BotAPI that stricts and simplifies
	// its functionality.
	Bot struct {
		RunBotRepo
		RecordMessageRepo
		public.SendMessageRepo
		api      *tbot.BotAPI
		commands commands.CallbackMap
		Database database.Database
		Modules  []botmodule.Module
	}
)

var (
	log = logger.Factory.NewLogger("bot")
)

// NewBot initializes bot api and returns a new *Bot.
func NewBot(token string, commands commands.CallbackMap, enableMongo bool, url string, databaseName string) (*Bot, error) {
	api, err := tbot.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	var databaseConfig database.Database

	if enableMongo {
		databaseConfig = &database.DatabaseImpl{
			Url:      url,
			Database: databaseName,
		}
	} else {
		databaseConfig = &database.DatabaseNoOp{}
	}

	b := &Bot{
		api:      api,
		commands: commands,
		Database: databaseConfig,
	}
	b.RegisterModule(youtube.NewModule(b))

	return b, nil
}
