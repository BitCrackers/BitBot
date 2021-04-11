package main

import (
	"fmt"
	"github.com/BitCrackers/BitBot/config"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"regexp"
)

func newFilterHandler(filter config.Filter) (interface{}, error) {
	var exps []*regexp.Regexp
	for _, expString := range filter.RegExp {
		exp, err := regexp.Compile(expString)
		if err != nil {
			return nil, fmt.Errorf("error while parsing regular expression \"%v\": %v", expString, err)
		}
		exps = append(exps, exp)
	}
	if len(exps) == 0 {
		return nil, fmt.Errorf("filter has no (valid) regular expressions")
	}

	return func(s *discordgo.Session, m *discordgo.MessageCreate) {
		var match bool
		for _, exp := range exps {
			if exp.MatchString(m.Content) {
				match = true
				break
			}
		}
		if !match {
			return
		}
		if filter.Delete {
			if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
				logrus.Errorf("error while deleting filtered message: %v", err)
			}
		}

		if filter.Response == "" {
			return
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, filter.Response); err != nil {
			logrus.Errorf("error while responding to filtered message: %v", err)
		}
	}, nil
}
