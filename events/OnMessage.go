package events

import (
	"github.com/BitCrackers/BitBot/commands"
	"github.com/bwmarrin/discordgo"
)

func OnMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// Ignore all messages from Bots and Itself.
	if m.Author.Bot || m.Author.ID == s.State.User.ID {
		return
	}

	// Basic Ping Command.
	if m.Message.Content == "!ping" {
		commands.Ping(s, m)
	}
}
