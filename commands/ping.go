package commands

import "github.com/bwmarrin/discordgo"

type CommandPing struct{}

func (c *CommandPing) Name() string {
	return "ping"
}

func (c *CommandPing) Description() string {
	return "Pong!"
}

func (c *CommandPing) AdminRequired() bool {
	return false
}

func (c *CommandPing) Options() []*discordgo.ApplicationCommandOption {
	return make([]*discordgo.ApplicationCommandOption, 0)
}

func (c *CommandPing) Exec(s *discordgo.Session, i *discordgo.InteractionCreate) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content: "Pong",
		},
	})
}
