package filters

import (
	"fmt"
	"github.com/BitCrackers/BitBot/internal/config"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
	"strings"
)

var AutoMod = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) bool {
		containsIllegal := false

		for _, b := range config.C.IllegalWords {
			if !strings.Contains(strings.ToLower(m.Content), strings.ToLower(b)) {
				continue
			}

			containsIllegal = true
			break
		}

		if !containsIllegal {
			return true
		}

		err := s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			fmt.Printf("error trying to delete message %v", err)
		}

		_, err = s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%v> you message was deleted because it contained an illegal word", m.Author.ID))
		if err != nil {
			fmt.Printf("error trying to send message %v", err)
			return true
		}

		return false
	},
}