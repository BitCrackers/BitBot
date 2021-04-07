package router

import "github.com/bwmarrin/discordgo"

type Command struct {
	Name          string
	Description   string
	Options       []*discordgo.ApplicationCommandOption
	AdminRequired bool
	Exec          func(s *discordgo.Session, i *discordgo.InteractionCreate)
}
