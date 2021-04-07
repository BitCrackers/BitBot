package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	CommandInstances    []Command
	CommandMap          map[string]Command
	ApplicationCommands []*discordgo.ApplicationCommand
	CommandFunctions    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	OnError             func(err error, ctx *Context)
}

func NewCommandHandler(prefix string) *CommandHandler {
	return &CommandHandler{
		CommandInstances:    make([]Command, 0),
		CommandMap:          make(map[string]Command),
		ApplicationCommands: make([]*discordgo.ApplicationCommand, 0),
		CommandFunctions:    make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		OnError:             func(error, *Context) {},
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.CommandInstances = append(c.CommandInstances, cmd)
	c.CommandMap[cmd.Name()] = cmd
	c.CommandFunctions[cmd.Name()] = cmd.Exec
	appCommand := discordgo.ApplicationCommand{
		Name:        cmd.Name(),
		Description: cmd.Description(),
		Options:     cmd.Options(),
	}

	c.ApplicationCommands = append(c.ApplicationCommands, &appCommand)
}

func (c *CommandHandler) Handler(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if h, ok := c.CommandFunctions[i.Data.Name]; ok {
		h(s, i)
	}
}

func (c *CommandHandler) CreateCommands(s *discordgo.Session, guildId string) {
	for _, v := range c.ApplicationCommands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}
