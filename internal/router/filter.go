package router

import "github.com/bwmarrin/discordgo"

type Filter struct {
	Exec        func(s *discordgo.Session, m *discordgo.Message)
}
