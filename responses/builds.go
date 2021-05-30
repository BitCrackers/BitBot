package responses

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var BuildResponse = Response{
	Name: "builds",
	Send: func(s *discordgo.Session, m *discordgo.Message, reply bool) error {
		embed := discordgo.MessageEmbed{
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
		}

		var err error

		if reply {
			_, err = s.ChannelMessageSendComplex(m.ChannelID, &discordgo.MessageSend{
				Embed:     &embed,
				Reference: m.Reference(),
			})
		} else {
			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		}

		if err != nil {
			return fmt.Errorf("error trying to send embed %v", err)
		}
		return nil
	},
}
