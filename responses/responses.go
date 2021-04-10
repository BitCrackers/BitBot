package responses

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BitCrackers/BitBot/internal/github"
	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"net/http"
)

var responses = []router.Response{
	{
		Name: "builds",
		Send: func(s *discordgo.Session, m *discordgo.Message) error {
			resp, err := http.Get("https://api.github.com/repos/BitCrackers/AmongUsMenu/actions/artifacts")
			if err != nil {
				return err
			}

			b, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			artifacts := github.ArtifactsResponse{}
			err = json.Unmarshal(b, &artifacts)
			if err != nil {
				return err
			}

			resp, err = http.Get("https://api.github.com/repos/BitCrackers/AmongUsMenu/actions/runs")
			if err != nil {
				return err
			}

			b, err = ioutil.ReadAll(resp.Body)
			if err != nil {
				return err
			}

			runs := github.RunResponse{}
			err = json.Unmarshal(b, &runs)
			if err != nil {
				return err
			}

			if len(artifacts.Artifacts) < 4 {
				return errors.New("incorrect amount of artifacts")
			}

			if len(artifacts.Artifacts) < 1 {
				return errors.New("incorrect amount of workflow runs")
			}

			embed := discordgo.MessageEmbed{
				Title:       "Builds",
				Description: "You have to be logged into github to download the following artifacts",
				Fields: []*discordgo.MessageEmbedField{
					{
						Name:   "Version Proxy",
						Value:  fmt.Sprintf("[[download]](https://github.com/BitCrackers/AmongUsMenu/suites/%v/artifacts/%v)", runs.WorkflowRuns[0].CheckSuiteID, artifacts.Artifacts[0].ID),
						Inline: true,
					},
					{
						Name:   "Injectable",
						Value:  fmt.Sprintf("[[download]](https://github.com/BitCrackers/AmongUsMenu/suites/%v/artifacts/%v)", runs.WorkflowRuns[0].CheckSuiteID, artifacts.Artifacts[1].ID),
						Inline: true,
					},
				},
			}

			_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
			if err != nil {
				return err
			}
			return nil
		},
	},
}

func GetCustomResponse(name string) (*router.Response, error) {
	for _, response := range responses {
		if response.Name != name {
			continue
		}
		return &response, nil
	}
	return nil, errors.New("could not find custom response")
}
