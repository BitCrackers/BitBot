package commands

import "github.com/bwmarrin/discordgo"

type Command interface {
	Name() string
	Description() string
	Options() []*discordgo.ApplicationCommandOption
	AdminRequired() bool
	Exec(s *discordgo.Session, i *discordgo.InteractionCreate)
}
