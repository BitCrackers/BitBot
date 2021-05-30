package commands

import (
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) BuildsCommand() *Command {
	return &Command{
		Name:        "builds",
		Description: "Gets the latest AmongUsMenu release.",
		Options:     []*discordgo.ApplicationCommandOption{},
		HandlerFunc: ch.handleBuilds,
	}
}

func (ch *CommandHandler) handleBuilds(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       "Builds",
					Description: "Here, I fetched the latest releases for you. If you're looking for a pre-release, use the Absolute Latest link.",
					Fields: []*discordgo.MessageEmbedField{
						{
							Name:   "Official Latest Release",
							Value:  "[[download]](https://github.com/BitCrackers/AmongUsMenu/releases/latest)",
							Inline: true,
						},
						{
							Name:   "Absolute Latest Release",
							Value:  "[[download]](https://github.com/BitCrackers/AmongUsMenu/releases)",
							Inline: true,
						},
					},
				},
			},
		},
	})
	if err != nil {
		logrus.Errorf("error while sending reponse embed for build command: %v", err)
	}
}
