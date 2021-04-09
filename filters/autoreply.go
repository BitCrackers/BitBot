package filters

import (
	"fmt"
	"strings"

	"github.com/BitCrackers/BitBot/internal/config"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var AutoReply = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) bool {
		shouldReply := false

		for _, b := range config.C.AutoReplyWithBuild {
			if !strings.Contains(strings.ToLower(m.Content), strings.ToLower(b)) {
				continue
			}

			shouldReply = true
			break
		}

		if !shouldReply {
			return true
		}

		embed := discordgo.MessageEmbed{
			Title: "Builds",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Version Proxy",
					Value:  "[download](https://github.com/BitCrackers/AmongUsMenu)",
					Inline: true,
				},
				{
					Name:   "Injectable",
					Value:  "[download](https://github.com/BitCrackers/AmongUsMenu)",
					Inline: true,
				},
			},
		}

		_, err := s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			fmt.Printf("error trying to send message %v", err)
			return true
		}

		return false
	},
}
