package commands

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	Prefix           string
	CommandInstances []Command
	CommandMap       map[string]Command

	OnError func(err error, ctx *Context)
}

func NewCommandHandler(prefix string) *CommandHandler {
	return &CommandHandler{
		Prefix:           prefix,
		CommandInstances: make([]Command, 0),
		CommandMap:       make(map[string]Command),
		OnError:          func(error, *Context) {},
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.CommandInstances = append(c.CommandInstances, cmd)
	for _, invoke := range cmd.Invokes() {
		c.CommandMap[invoke] = cmd
	}
}

func (c *CommandHandler) HandleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot || s.State.User.ID == m.Author.ID || !strings.HasPrefix(m.Content, c.Prefix) {
		return
	}

	split := strings.Split(m.Content[len(c.Prefix):], " ")
	if len(split) < 1 {
		return
	}

	// This represents the "command," or the arg just after the prefix. i.e.: !kick -> [kick]
	invoke := split[0]
	// Everything after the first arg.
	args := split[1:]

	cmd, ok := c.CommandMap[invoke]
	if !ok || cmd == nil {
		return
	}

	ctx := &Context{
		Session: s,
		Args:    args,
		Handler: c,
		Message: m.Message,
	}

	if err := cmd.Exec(ctx); err != nil {
		c.OnError(err, ctx)
	}
}
