package modlog

import (
	"errors"
	"github.com/BitCrackers/BitBot/config"
	"github.com/bwmarrin/discordgo"
)

type ModLogHandler struct {
	Config *config.Config
	Enabled bool
}

func Create(cfg *config.Config, s *discordgo.Session) (ModLogHandler, error) {
	if cfg.ModLogChannelId == "" {
		return ModLogHandler{ Enabled: false }, nil
	}

	_, err := s.Channel(cfg.ModLogChannelId)
	if err != nil {
		return ModLogHandler{}, errors.New("cannot find modlog channel")
	}
	return ModLogHandler{Config: cfg, Enabled: true}, nil
}

func (m *ModLogHandler) SendEmbed(s *discordgo.Session, embed *discordgo.MessageEmbed) error {
	if !m.Enabled {
		return nil
	}
	_, err := s.ChannelMessageSendEmbed(m.Config.ModLogChannelId, embed)
	return err
}

func (m *ModLogHandler) SendMessage(s *discordgo.Session, message string) error {
	if !m.Enabled {
		return nil
	}
	_, err := s.ChannelMessageSend(m.Config.ModLogChannelId, message)
	return err
}