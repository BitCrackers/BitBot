package responses

import (
	"fmt"

	"github.com/BitCrackers/BitBot/github"
	"github.com/bwmarrin/discordgo"
)

var BuildResponse = Response{
	Name: "builds",
	Send: func(s *discordgo.Session, m *discordgo.Message, reply bool) error {
		artifacts, err := github.Artifacts("BitCrackers", "AmongUsMenu")
		if err != nil {
			return fmt.Errorf("error while getting artifacts: %v", err)
		}

		run, err := github.GetLatestMasterWorkflowRun("BitCrackers", "AmongUsMenu")
		if err != nil {
			return fmt.Errorf("error while getting workflow runs: %v", err)

		}

		if len(artifacts) < 4 {
			return fmt.Errorf("unexpected amount of artifacts: %v", len(artifacts))
		}

		embed := discordgo.MessageEmbed{
			Title:       "Builds",
			Description: "You have to be logged into github to download the following artifacts",
			Fields: []*discordgo.MessageEmbedField{
				{
					Name: "Version Proxy",
					Value: fmt.Sprintf(
						"[[download]](https://github.com/BitCrackers/AmongUsMenu/suites/%v/artifacts/%v)",
						run.CheckSuiteID,
						artifacts[0].ID,
					),
					Inline: true,
				},
				{
					Name: "Injectable",
					Value: fmt.Sprintf(
						"[[download]](https://github.com/BitCrackers/AmongUsMenu/suites/%v/artifacts/%v)",
						run.CheckSuiteID,
						artifacts[1].ID,
					),
					Inline: true,
				},
			},
		}

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
