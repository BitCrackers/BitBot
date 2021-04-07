package commands

import "github.com/BitCrackers/BitBot/internal/commands"

type CommandPing struct{}

func (c *CommandPing) Invokes() []string {
	return []string{"ping"}
}

func (c *CommandPing) Description() string {
	return "Pong!"
}

func (c *CommandPing) AdminRequired() bool {
	return false
}

func (c *CommandPing) Exec(ctx *commands.Context) error {
	_, err := ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Pong!")
	if err != nil {
		return err
	}
	return nil
}
