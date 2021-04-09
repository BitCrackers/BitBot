package filters

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BitCrackers/BitBot/internal/config"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/BitCrackers/BitBot/responses"
	"github.com/bwmarrin/discordgo"
)

var AutoMod = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) bool {
		perfectMatch := false
		response := ""
		deleteMessage := false

		for _, f := range config.C.Filters {

			r, err := regexp.Compile(f.RegExp)

			if err != nil {
				fmt.Printf("error trying to compile regex %v\n", err)
				return true
			}

			if r.MatchString(m.Content) {
				perfectMatch = true
				response = f.Response
				deleteMessage = f.Delete
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
