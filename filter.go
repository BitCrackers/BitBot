package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/BitCrackers/BitBot/config"
	"github.com/BitCrackers/BitBot/responses"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func newFilterHandler(c *config.Config, f config.Filter, rh *responses.CustomResponseHandler) (interface{}, error) {
	var exps []*regexp.Regexp
	for _, expString := range f.RegExp {
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
		for _, a := range c.Moderators {
			if m.Author.ID == a {
				return
			}
		}

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
		if f.Delete {
			if err := s.ChannelMessageDelete(m.ChannelID, m.ID); err != nil {
				logrus.Errorf("error while deleting filtered message: %v", err)
			}
		}

		if strings.Contains(f.Response, "custom#") {
			if len(strings.Split(f.Response, "#")) < 2 {
				logrus.Error("error while responding to with custom filter: invalid custom response formatting")
				return
			}
			res, err := rh.GetCustomResponse(strings.Split(f.Response, "#")[1])
			if err != nil {
				logrus.Errorf("error while responding to filtered message with custom response: %v", err)
			}
			err = res.Send(s, m.Message, !f.Delete)
			if err != nil {
				logrus.Errorf("error while responding to filtered message with custom response: %v", err)
			}
			return
		}

		if f.Response == "" {
			return
		}

		if _, err := s.ChannelMessageSend(m.ChannelID, f.Response); err != nil {
			logrus.Errorf("error while responding to filtered message: %v", err)
		}
	}, nil
}
