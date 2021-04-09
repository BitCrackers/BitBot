package filters

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/BitCrackers/BitBot/internal/router"
	"github.com/bwmarrin/discordgo"
)

type runInfo struct {
	commitHash string
	branch     string
	buildType  string
	auVersion  string
}

var LogParser = router.Filter{
	Exec: func(s *discordgo.Session, m *discordgo.Message) bool {
		if len(m.Attachments) == 0 || len(m.Attachments) > 1 {
			return true
		}

		if m.Attachments[0].Filename != "aum-log.txt" {
			return true
		}

		resp, err := http.Get(m.Attachments[0].URL)
		if err != nil {
			fmt.Printf("error trying to fetch log %v\n", err)
			return true
		}

		b, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Printf("error trying read log body %v\n", err)
			return true
		}

		logs := strings.Split(string(b), "\n")
		r := runInfo{}
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
				fmt.Printf("error trying to send reply %v\n", err)
				return true
			}
			return true
		}

		if infoIndex == -1 {
			fmt.Printf("invalid log format: unable to find run info\n")
			return true
		}

		re := *regexp.MustCompile(`^\tBuild:\s(.*)$`)
		if re.MatchString(strings.TrimSpace(logs[infoIndex+2])) {
			fmt.Printf("invalid log format: unable to find build type\n")
			return true
		}
		r.buildType = strings.Split(logs[infoIndex+2], ":")[1]

		re = *regexp.MustCompile(`^\tCommit: (.*) - (.*)$`)
		if re.MatchString(strings.TrimSpace(logs[infoIndex+3])) {
			fmt.Printf("invalid log format: unable to find git info\n")
			return true
		}
		git := strings.Split(logs[infoIndex+3], ":")[1]
		r.commitHash = strings.TrimSpace(strings.Split(git, " - ")[0])
		r.branch = strings.TrimSpace(strings.Split(git, " - ")[1])

		re = *regexp.MustCompile(`^\tAmong Us Version: (.*)$`)
		if re.MatchString(strings.TrimSpace(logs[infoIndex+4])) {
			fmt.Printf("invalid log format: unable to find among us version\n")
			return true
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
					Value:  fmt.Sprintf("[%s](https://github.com/BitCrackers/AmongUsMenu/commit/%s)", r.commitHash[0:7], r.commitHash),
					Inline: false,
				},
			},
		}

		_, err = s.ChannelMessageSendEmbed(m.ChannelID, &embed)
		if err != nil {
			fmt.Printf("error trying to send embed %v\n", err)
			return true
		}

		return false
	},
}
