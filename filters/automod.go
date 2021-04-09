package filters

import (
	"fmt"
	"github.com/BitCrackers/BitBot/responses"
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
		deleteMessage := false

		for _, f := range config.C.Filters {
			for _, w := range f.Words {
				if !strings.Contains(strings.ToLower(m.Content), strings.ToLower(w)) {
					continue
				}
				numMatch += 1

				// We have successfully hit a Filter from the config file.
				if numMatch == len(f.Words) {
					perfectMatch = true
					response = f.Response
					deleteMessage = f.Delete
				}
			}
		}

		if !perfectMatch {
			return false
		}

		if deleteMessage {
			err := s.ChannelMessageDelete(m.ChannelID, m.ID)
			if err != nil {
				fmt.Printf("error trying to delete message %v", err)
			}
		}

		if strings.Contains(response, "custom#") {
			resp, err := responses.GetCustomResponse(strings.Split(response, "#")[1])
			if err != nil {
				fmt.Printf("error trying to get custom response: %v\n", err)
				return false
			}
			err = resp.Send(s, m)
			if err != nil {
				fmt.Printf("error trying to get send response: %v\n", err)
			}
			return false
		}

		_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("<@%v>: %s", m.Author.ID, response))
		if err != nil {
			fmt.Printf("error trying to send message %v\n", err)
			return false
		}

		return deleteMessage
	},
}
