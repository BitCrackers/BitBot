package database

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"time"
)

func (d *Database) startJanitor() chan struct{} {
	c := make(chan struct{})
	go func() {
		t := time.NewTicker(time.Duration(d.config.JanitorCycleDuration) * time.Second)
		defer t.Stop()
		for {
			select {
			case <-t.C:
				logrus.Debugf("Executing janitor cycle")
			case <-c:
				return
			}
			records, err := d.AllUserRecords()
			if err != nil {
				logrus.Errorf("janitor error while getting all user records: %v", err)
				return
			}

			guild, err := d.session.Guild(d.config.GuildID)
			if err != nil {
				logrus.Errorf("janitor error while getting all guild: %v", err)
				return
			}

			for _, record := range records {
				inGuild := false
				for _, member := range guild.Members {
					if member.User.ID != record.ID {
						continue
					}
					inGuild = true
					break
				}
				if !inGuild {
					continue
				}

				if !record.Mute.Empty() {
					member, err := d.session.GuildMember(d.config.GuildID, record.ID)
					if err != nil {
						logrus.Errorf("Janitor error while getting guild member roles: %v", err)
						continue
					}
					var hasMuteRole bool
					for _, role := range member.Roles {
						if role == d.config.MuteRoleID {
							hasMuteRole = true
						}
					}
					if !hasMuteRole {
						logrus.Debugf("Janitor found muted user without muted role")
						err = d.session.GuildMemberRoleAdd(d.config.GuildID, record.ID, d.config.MuteRoleID)
						if err != nil {
							logrus.Errorf("Janitor error while adding muted role to user: %v", err)
						}
					}
				}

				if !record.Mute.Empty() &&
					record.Mute.Length != -1 &&
					record.Mute.Date.Add(record.Mute.Length).Before(time.Now()) {
					logrus.Debugf("Janitor removing expired mute")
					if err = d.UnmuteRecord(record, true); err != nil {
						logrus.Errorf("%v", err)
					}
				}

				if !record.Ban.Empty() &&
					record.Ban.Length != -1 &&
					record.Ban.Date.Add(record.Ban.Length).Before(time.Now()) {
					logrus.Debugf("Janitor removing expired ban")
					if err = d.UnbanRecord(record, true); err != nil {
						logrus.Errorf("%v", err)
					}
				}
			}
		}
	}()
	return c
}

func (d *Database) UnbanRecord(ur UserRecord, log bool) error {
	if err := d.session.GuildBanDelete(d.config.GuildID, ur.ID); err != nil {
		return fmt.Errorf("error while removing user ban: %v", err)
	}
	
	ur.Ban = Punishment{Type: -1}
	if err := d.SetUserRecord(ur); err != nil {
		logrus.Errorf("error while deleting user record from database: %v", err)
	}

	if !log {
		return nil
	}

	user, err := d.session.User(ur.ID)
	if err != nil {
		return fmt.Errorf("error while getting discordgo user from id: %v", err)
	}

	err = d.modlog.SendEmbed(d.session, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:         fmt.Sprintf("[UNMUTE] %s#%s", user.Username, user.Discriminator),
			IconURL:      user.AvatarURL("256"),
		},
		Description: "**Unbanned because ban expired**",
		Timestamp: time.Now().Format(time.RFC3339),
		Color: 3574686,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  fmt.Sprintf("<@%s>", user.ID),
				Inline: true,
			},
			{
				Name:   "Moderator",
				Value:  "**BitBot Janitor**",
				Inline: true,
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error logging unban: %v", err)
	}

	return nil
}

func (d *Database) UnmuteRecord(ur UserRecord, log bool) error {
	err := d.session.GuildMemberRoleRemove(d.config.GuildID, ur.ID, d.config.MuteRoleID)
	if err != nil {
		return fmt.Errorf("error while removing user role: %v", err)
	}
	ur.Mute = Punishment{Type: -1}
	if err = d.SetUserRecord(ur); err != nil {
		return fmt.Errorf("error while deleting user record from database: %v", err)
	}

	if !log {
		return nil
	}

	user, err := d.session.User(ur.ID)
	if err != nil {
		return fmt.Errorf("error while getting discordgo user from id: %v", err)
	}

	err = d.modlog.SendEmbed(d.session, &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:         fmt.Sprintf("[UNMUTE] %s#%s", user.Username, user.Discriminator),
			IconURL:      user.AvatarURL("256"),
		},
		Description: "**Unmuted because mute expired**",
		Timestamp: time.Now().Format(time.RFC3339),
		Color: 3574686,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "User",
				Value:  fmt.Sprintf("<@%s>", user.ID),
				Inline: true,
			},
			{
				Name:   "Moderator",
				Value:  "**BitBot Janitor**",
				Inline: true,
			},
		},
	})
	if err != nil {
		logrus.Errorf("Error logging unmute: %v", err)
	}

	return nil
}
