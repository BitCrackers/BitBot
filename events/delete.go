package events

import (
	"fmt"
	"github.com/BitCrackers/BitBot/filters"
	"github.com/BitCrackers/BitBot/internal/router"

	"github.com/bwmarrin/discordgo"
)

type DeleteHandler struct {
	Filters []router.Filter
}

func NewDeleteHandler() *DeleteHandler {
	return &DeleteHandler{
		Filters: []router.Filter{},
	}
}

func (h *DeleteHandler) AddFilter(f router.Filter) {
	h.Filters = append(h.Filters, f)
}

func (h *DeleteHandler) Handler(s *discordgo.Session, d *discordgo.MessageDelete) {
	m := filters.GetMessageFromCache(d.ID)
	if m == nil {
		fmt.Print("Unable to get message from cache")
		return
	}
	filters.DeleteFromCache(d.ID)

	// We don't want bot logs.
	if m.Author.Bot || s.State.User.ID == m.Author.ID {
		return
	}

	for _, f := range h.Filters {
		if !f.Exec(s, m) {
			break
		}
	}
}
