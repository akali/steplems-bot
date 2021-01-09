package eval

import "github.com/akali/steplems-bot/app/commands"

const (
	// CommandName is the name for command "eval".
	CommandName = "eval"
)

var (
	// Command is the composed Command for "eval".
	// "eval" sends a help message that has usage instructions
	// of the curl bot.
	Command = commands.NewCommand(CommandName, CommandCallback)
)
