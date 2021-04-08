package events

import (
	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

type DeleteHandler struct{
	Filters []router.Filter
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		Filters: []router.Filter{},
	}
}

func (h *DeleteHandler) AddFilter(f router.Filter)  {
	h.Filters = append(h.Filters, f)
}

func (h *DeleteHandler) Handler(s *discordgo.Session, d *discordgo.MessageDelete) {
	// We don't want bot logs.

	if d.Message.Author.Bot || s.State.User.ID == d.Message.Author.ID {
		return
	}
	for _, f := range h.Filters {
		f.Exec(s, d.Message)
	}
}
