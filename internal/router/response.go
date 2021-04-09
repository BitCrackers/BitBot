package router

import "github.com/bwmarrin/discordgo"

type Response struct {
	Name string
	Send func(s *discordgo.Session, m *discordgo.Message) error
}
