package main

import (
	"github.com/BitCrackers/BitBot/config"
	"github.com/BitCrackers/BitBot/database"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

type ReactionRoleHandler struct {
	DB     *database.Database
	Config *config.Config
}

func (rh *ReactionRoleHandler) reactionRoleAddHandler(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	e, err := rh.DB.ReactionRoleExist(m.MessageID)
	if err != nil {
		logrus.Errorf("Error while getting reaction role entry from database %v", err)
		return
	}

	if !e {
		return
	}

	r, err := rh.DB.GetReaction(m.MessageID)
	if err != nil {
		logrus.Errorf("Error while getting reaction role entry from database %v", err)
		return
	}

	err = s.GuildMemberRoleAdd(rh.Config.GuildID, m.UserID, r.Role)
	if err != nil {
		logrus.Errorf("Error while adding reaction role to user %v", err)
		return
	}
}

func (rh *ReactionRoleHandler) reactionRoleRemoveHandler(s *discordgo.Session, m *discordgo.MessageReactionRemove) {
	e, err := rh.DB.ReactionRoleExist(m.MessageID)
	if err != nil {
		logrus.Errorf("Error while getting reaction role entry from database %v", err)
		return
	}

	if !e {
		return
	}

	r, err := rh.DB.GetReaction(m.MessageID)
	if err != nil {
		logrus.Errorf("Error while getting reaction role entry from database %v", err)
		return
	}

	err = s.GuildMemberRoleRemove(rh.Config.GuildID, m.UserID, r.Role)
	if err != nil {
		logrus.Errorf("Error while removing reaction role from user %v", err)
		return
	}
}
