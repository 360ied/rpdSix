package extendeddiscordobjects

import (
	"errors"
	"fmt"

	"github.com/bwmarrin/discordgo"

	"rpdSix/helpers"
)

type ExtendedGuild struct {
	*discordgo.Guild
	session *discordgo.Session
}

func ExtendGuild(guild *discordgo.Guild, session *discordgo.Session) *ExtendedGuild {
	return &ExtendedGuild{
		Guild:   guild,
		session: session,
	}
}

func (guild *ExtendedGuild) GetRole(roleID string) (*discordgo.Role, error) {
	for _, role := range guild.Roles {
		if role.ID == roleID {
			return role, nil
		}
	}
	return nil, errors.New(fmt.Sprint(helpers.RoleNotFoundErrorTemplate,
		"Role ", roleID, " not found in guild ", guild.ID))
}

// Better than GetRole as this finds all the Roles in 1 pass
func (guild *ExtendedGuild) GetRoles(roleIDs map[string]struct{}) (roles []*discordgo.Role) {
	for _, role := range guild.Roles {
		var _, contains = roleIDs[role.ID]
		if contains {
			roles = append(roles, role)
		}
	}
	return
}

// GetRoles but your slice is converted into a map
func (guild *ExtendedGuild) GetRolesSlice(roleIDs []string) (roles []*discordgo.Role) {
	var roleIDMap = make(map[string]struct{})
	for _, roleID := range roleIDs {
		roleIDMap[roleID] = struct{}{}
	}
	return guild.GetRoles(roleIDMap)
}
