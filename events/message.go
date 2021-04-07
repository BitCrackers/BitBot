package events

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct{}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{}
}

func (h *MessageHandler) Handler(s *discordgo.Session, m *discordgo.MessageCreate) {

	// We don't want bot logs.
	if m.Author.Bot || s.State.User.ID == m.Author.ID {
		return
	}

	_, err := s.Channel(m.ChannelID)

	if err != nil {
		fmt.Println("Failed getting channel from MessageCreate event: ", err)
		return
	}
	fmt.Printf("%s#%s: %s\n", m.Author.Username, m.Author.Discriminator, m.Content)
}
