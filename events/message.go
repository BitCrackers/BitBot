package events

import (
	"fmt"
	"github.com/BitCrackers/BitBot/filters"
	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

type MessageHandler struct {
	Filters []router.Filter
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		Filters: []router.Filter{
			filters.Cache,
			filters.AutoMod,
			filters.AutoReply,
			filters.LogParser,
		},
	}
}

func (h *MessageHandler) AddFilter(f router.Filter) {
	h.Filters = append(h.Filters, f)
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

	for _, f := range h.Filters {
		if !f.Exec(s, m.Message) {
			break
		}
	}

	fmt.Printf("%s#%s: %s\n", m.Author.Username, m.Author.Discriminator, m.Content)
}
