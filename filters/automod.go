package filters

import (
	"fmt"
	"strings"

	"github.com/BitCrackers/BitBot/internal/config"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

var AutoMod = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) bool {
		numMatch := 0
		perfectMatch := false
		response := ""
		delete := false

		for _, f := range config.C.Filters {
			for _, w := range f.Words {
				if strings.Contains(strings.ToLower(m.Content), strings.ToLower(w)) {
					numMatch += 1

					// We have successfully hit a Filter from the config file.
					if numMatch == len(f.Words) {
						perfectMatch = true
						response = f.Response
						delete = f.Delete
					}

					continue
				}
			}
		}

		if perfectMatch {
			if delete {
				err := s.ChannelMessageDelete(m.ChannelID, m.ID)
				if err != nil {
					fmt.Printf("error trying to delete message %v", err)
				}
			}

			_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%v>: %s", m.Author.ID, response))
			if err != nil {
				fmt.Printf("error trying to send message %v", err)
				return true
			}

			return false
		}

		return true
	},
}
