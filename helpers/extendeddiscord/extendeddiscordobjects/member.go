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

	if memberGuild.OwnerID == member.User.ID {
		return true, nil
	}

	var extendedGuild = ExtendGuild(memberGuild, member.session)

	if len(member.Roles) == 0 {
		var reGetMember, reGetMemberErr = member.session.GuildMember(member.GuildID, member.User.ID)
		if reGetMemberErr != nil {
			return false, reGetMemberErr
		}
		member.Roles = reGetMember.Roles
	}

	var roles, rolesErr = extendedGuild.GetRolesSlice(member.Roles)
	if rolesErr != nil {
		return false, rolesErr
	}

	var combinedPermissionInteger = 0

	for _, role := range roles {
		// Bitwise OR
		combinedPermissionInteger |= role.Permissions
	}

	return extendeddiscordpermissions.IsPermittedAll(combinedPermissionInteger, permissions...), nil
}
