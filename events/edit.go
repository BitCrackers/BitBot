package events

import (
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

func (h *EditHandler) Handler(s *discordgo.Session, e *discordgo.MessageUpdate) {
	// We don't want bot logs.
	if e.Message.Author.Bot || s.State.User.ID == e.Message.Author.ID {
		return
	}
	for _, f := range h.Filters {
		if !f.Exec(s, e.Message) {
			break
		}
	}
}
