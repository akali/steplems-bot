package main

import (
	"github.com/akali/steplems-bot/app/bot"
	"github.com/akali/steplems-bot/app/commands"
	"github.com/akali/steplems-bot/app/commands/eval"
	"github.com/akali/steplems-bot/app/commands/help"
	"github.com/akali/steplems-bot/app/commands/backup"
	"github.com/akali/steplems-bot/app/commands/request"
	"github.com/akali/steplems-bot/app/config"
	"github.com/akali/steplems-bot/app/logger"
	"time"
)

var (
	log = logger.Factory.NewLogger("main")
	// Mapping commands into a map to make command selection easier.
	cmds = commands.NewCallbackMap(
		help.Command,
		eval.Command,
		backup.Command,
		request.CommandGet,
		request.CommandHead,
		request.CommandPost,
		request.CommandPut,
		request.CommandPatch,
		request.CommandDelete,
		request.CommandConnect,
		request.CommandOptions,
		request.CommandTrace,
	)
)

func main() {
	// Creating and setting up a new bot api client.
	b, err := bot.NewBot(config.BotAPIToken, cmds, config.MongoConnectionString, config.MongoDatabaseName)
	if err != nil {
		log.Panic.Println("error trying to initialize a new bot:", err)
	}

	err, _ = b.Database.Init(time.Duration(config.UpdateTimeout))

	if err != nil {
		log.Panic.Println("error trying to connect to mongo:", err.Error())
	}

	// Run is going to loop a continues chan that will block
	// the further execution of main func.
	err = b.Run(config.UpdateTimeout)
	if err != nil {
		log.Panic.Println("error trying to run bot:", err)
	}
}
