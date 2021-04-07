package router

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type CommandHandler struct {
	ApplicationCommands []*discordgo.ApplicationCommand
	CommandFunctions    map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)
	OnError             func(err error, ctx *Context)
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		ApplicationCommands: make([]*discordgo.ApplicationCommand, 0),
		CommandFunctions:    make(map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)),
		OnError:             func(error, *Context) {},
	}
}

func (c *CommandHandler) RegisterCommand(cmd Command) {
	c.CommandFunctions[cmd.Name] = cmd.Exec
	appCommand := discordgo.ApplicationCommand{
		Name:        cmd.Name,
		Description: cmd.Description,
		Options:     cmd.Options,
	}

	fmt.Printf("> Registering command: %s.\n", cmd.Name)

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
			fmt.Printf("Cannot create '%v' command: %v", v.Name, err)
		}
	}
}

func (c *CommandHandler) ClearCommands(s *discordgo.Session, guildId string) {
	commands, err := s.ApplicationCommands(s.State.User.ID, guildId)
	if err != nil {
		fmt.Printf("Cannot fetch existing commands")
	}

	for _, c := range commands {
		err := s.ApplicationCommandDelete(s.State.User.ID, guildId, c.ID)
		if err != nil {
			fmt.Printf("Cannot delete command: %v", err)
		}
	}
}
