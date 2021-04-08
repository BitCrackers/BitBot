package filters

import (
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var messageCache = map[string]*discordgo.Message{}

var Cache = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) {
		messageCache[m.ID] = m
	},
}

func GetMessageFromCache(id string) *discordgo.Message {
	return messageCache[id]
}

func DeleteFromCache(id string) {
	messageCache[id] = nil
}