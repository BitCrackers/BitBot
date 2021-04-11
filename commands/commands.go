package commands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// Command represents a Discord slash command.
type Command struct {
	Name        string
	Description string
	Options     []*discordgo.ApplicationCommandOption
	HandlerFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)

	session                  *discordgo.Session
	appCmd                   *discordgo.ApplicationCommand
	guildID                  string
	removeHandlerFromSession func()
}

// Register registers an interaction with Discord and adds a handler to the
// session for it.
func (cmd *Command) Register(session *discordgo.Session, guildID string) error {
	if cmd.session != nil {
		panic("attempted to register an already registered command")
	}

	appCmd, err := session.ApplicationCommandCreate(session.State.User.ID, guildID, &discordgo.ApplicationCommand{
		ID:            "",
		ApplicationID: "",
		Name:          cmd.Name,
		Description:   cmd.Description,
		Options:       cmd.Options,
	})
	if err != nil {
		return fmt.Errorf("error while creating application command: %v", err)
	}

	cmd.session, cmd.guildID, cmd.appCmd = session, guildID, appCmd
	cmd.removeHandlerFromSession = session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if i.Data.Name != cmd.Name {
			return
		}
		cmd.HandlerFunc(s, i)
	})
	return nil
}

// Delete deletes the command from the session where it was previously
// registered. After deleting it, attempting to register it again will result in
// a panic.
func (cmd *Command) Delete() error {
	if cmd.session == nil || cmd.removeHandlerFromSession == nil {
		panic("attempted to delete a command which is not registered")
	}

	err := cmd.session.ApplicationCommandDelete(cmd.session.State.User.ID, cmd.guildID, cmd.appCmd.ID)
	if err != nil {
		return fmt.Errorf("error while deleting application command: %v", err)
	}

	cmd.removeHandlerFromSession()
	return nil
}
