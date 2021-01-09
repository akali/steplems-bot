package help

import "github.com/akali/steplems-bot/app/commands"

const (
	// CommandName is the name for command "help".
	CommandName = "backup"
)

var (
	// Command is the composed Command for "help".
	// "help" sends a help message that has usage instructions
	// of the curl bot.
	Command = commands.NewCommand(CommandName, CommandCallback)
)
