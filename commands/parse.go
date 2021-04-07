package commands

import (
	"github.com/BitCrackers/BitBot/internal/commands"
)

type CommandParse struct{}

func (c *CommandParse) Invokes() []string {
	return []string{"parse"}
}

func (c *CommandParse) Description() string {
	return "Parses a mesage. Debug command."
}

func (c *CommandParse) AdminRequired() bool {
	return false
}

func (c *CommandParse) Exec(ctx *commands.Context) error {
	s := ""
	for _, a := range ctx.Args {
		s += "["
		s += a
		s += "] "
	}

	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, s)
	if err != nil {
		return err
	}
	return nil
}
