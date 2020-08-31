package extendeddiscordobjects

import (
	"github.com/bwmarrin/discordgo"

	"rpdSix/helpers/extendeddiscord/extendeddiscordpermissions"
)

type ExtendedMember struct {
	*discordgo.Member
	session *discordgo.Session
}

func ExtendMember(member *discordgo.Member, session *discordgo.Session) *ExtendedMember {
	return &ExtendedMember{
		Member:  member,
		session: session,
	}
}

func (member *ExtendedMember) Guild() (*discordgo.Guild, error) {
	return member.session.Guild(member.GuildID)
}

func (member *ExtendedMember) HasAllPermissions(permissions ...int) (bool, error) {
	var memberGuild, memberGuildErr = member.Guild()
	if memberGuildErr != nil {
		return false, memberGuildErr
	}
	var extendedGuild = ExtendGuild(memberGuild, member.session)

	var roles = extendedGuild.GetRolesSlice(member.Roles)

	var combinedPermissionInteger = 0

	for _, role := range roles {
		// Bitwise OR
		combinedPermissionInteger |= role.Permissions
	}

	return extendeddiscordpermissions.HasAllPermissions(combinedPermissionInteger, permissions...), nil
}
