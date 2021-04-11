package main

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

// TODO: Reduce the complexity of this function (i.e. split it up).
func aumLog(s *discordgo.Session, m *discordgo.MessageCreate) {
	if len(m.Attachments) == 0 || len(m.Attachments) > 1 {
		return
	}
	if m.Attachments[0].Filename != "aum-log.txt" {
		return
	}
	resp, err := http.Get(m.Attachments[0].URL)
	if err != nil {
		logrus.Errorf("Error while fetching AmongUsMenu log: %v", err)
		return
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Errorf("Error while reading HTTP response body: %v", err)
		return
	}

	logs := strings.Split(string(b), "\n")
	r := struct{
		commitHash string
		branch     string
		buildType  string
		auVersion  string
	}{}
	infoIndex := -1

	for i, log := range logs {
		if !strings.Contains(log, "[INFO - AUM - Run]") {
			continue
		}
		infoIndex = i
		break
	}

	if infoIndex == -1 && len(logs) < 5 {
		message := "```"
		for _, log := range logs {
			message += log
		}
		message += "```"

		_, err = s.ChannelMessageSendReply(m.ChannelID, message, m.MessageReference)
		if err != nil {
			logrus.Errorf("Error trying to send reply %v", err)
			return
		}
		return
	}

	if infoIndex == -1 {
		// TODO: Make this into a Discord message response.
		logrus.Errorf("Invalid log format: unable to find run info")
		return
	}

	re := *regexp.MustCompile(`^\tBuild:\s(.*)$`)
	if re.MatchString(strings.TrimSpace(logs[infoIndex+2])) {
		// TODO: Make this into a Discord message response.
		logrus.Errorf("Invalid log format: unable to find build type")
		return
	}
	r.buildType = strings.Split(logs[infoIndex+2], ":")[1]

	re = *regexp.MustCompile(`^\tCommit: (.*) - (.*)$`)
	if re.MatchString(strings.TrimSpace(logs[infoIndex+3])) {
		// TODO: Make this into a Discord message response.
		logrus.Errorf("Invalid log format: unable to find git info")
		return
	}
	git := strings.Split(logs[infoIndex+3], ":")[1]
	r.commitHash = strings.TrimSpace(strings.Split(git, " - ")[0])
	r.branch = strings.TrimSpace(strings.Split(git, " - ")[1])

	re = *regexp.MustCompile(`^\tAmong Us Version: (.*)$`)
	if re.MatchString(strings.TrimSpace(logs[infoIndex+4])) {
		// TODO: Make this into a Discord message response.
		logrus.Errorf("Invalid log format: unable to find among us version")
		return
	}
	r.auVersion = strings.Split(logs[infoIndex+4], ":")[1]
	embed := discordgo.MessageEmbed{
		Title: "AmongUsMenu log",
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Build type",
				Value:  r.buildType,
				Inline: false,
			},
			{
				Name:   "AU Version",
				Value:  r.auVersion,
				Inline: false,
			},
			{
				Name:   "Git branch",
				Value:  r.branch,
				Inline: false,
			},
			{
				Name:   "Git commit",
				Value:  fmt.Sprintf("[%s](https://github.com/BitCrackers/AmongUsMenu/commit/%s)", r.commitHash[:7], r.commitHash),
				Inline: false,
			},
		},
	}

	_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		logrus.Errorf("Error trying to send embed %v", err)
		return
	}

	return
}
