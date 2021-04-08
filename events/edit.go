package events

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

type EditHandler struct{
	Filters []router.Filter
}

func NewEditHandler() *EditHandler {
	return &EditHandler{
		Filters: []router.Filter{},
	}
}

func (h *EditHandler) AddFilter(f router.Filter)  {
	h.Filters = append(h.Filters, f)
}

func (h *EditHandler) Handler(s *discordgo.Session, e *discordgo.MessageEdit) {
	m, err := s.ChannelMessage(e.Channel, e.ID)
	if err != nil {
		fmt.Println("Failed getting message from MessageEdit event: ", err)
		return
	}
	// We don't want bot logs.

	if m.Author.Bot || s.State.User.ID == m.Author.ID {
		return
	}
	for _, f := range h.Filters {
		f.Exec(s, m)
	}
}
