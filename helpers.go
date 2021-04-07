package main

import "github.com/bwmarrin/discordgo"

func CheckPermissionForUser(s *discordgo.Session, m discordgo.Member, p int64) (bool, error) {
	for _, r := range m.Roles {
		role, err := s.State.Role(m.GuildID, r)

		if err != nil {
			return false, err
		}

		if role.Permissions&p == p {
			return true, nil
		}
	}

	return false, nil
}
