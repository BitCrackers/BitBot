package commands

import (
	"github.com/BitCrackers/BitBot/database"
)

// CommandHandler provides a shared state for all command handler functions.
type CommandHandler struct {
	DB         *database.Database
	Moderators []string
	Debug      bool
}

// Commands returns all commands of the CommandHandler. This should be updated
// to contain all `XxxxCommand()` methods of the handler.
func (ch *CommandHandler) Commands() []*Command {
	cmds := []*Command{
		ch.KickCommand(),
		ch.BanCommand(),
		ch.WarnCommand(),
		ch.BuildsCommand(),
	}
	if ch.Debug {
		return append(
			cmds,
			ch.PingCommand(),
			ch.ParseCommand(),
		)
	}
	return cmds
}
