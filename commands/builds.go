package commands

import (
	"fmt"
	"github.com/BitCrackers/BitBot/github"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func (ch *CommandHandler) BuildsCommand() *Command {
	return &Command{
		Name:        "builds",
		Description: "Gets the latest AmongUsMenu builds",
		Options:     []*discordgo.ApplicationCommandOption{},
		HandlerFunc: ch.handleBuilds,
	}
}

func (ch *CommandHandler) handleBuilds(s *discordgo.Session, i *discordgo.InteractionCreate) {
	artifacts, err := github.Artifacts("BitCrackers", "AmongUsMenu")
	if err != nil {
		logrus.Errorf("error while getting artifacts: %v", err)
		return
	}

	run, err := github.GetLatestMasterWorkflowRun("BitCrackers", "AmongUsMenu")
	if err != nil {
		logrus.Errorf("error while getting workflow runs: %v", err)
		return
	}

	if len(artifacts) < 4 {
		logrus.Errorf("unexpected amount of artifacts: %v", len(artifacts))
		return
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
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
				},
			},
		},
	})
	if err != nil {
		logrus.Errorf("error while sending reponse embed for build command: %v", err)
	}
}
