package checkedrun

import (
	"errors"
	"fmt"

	"github.com/ztrue/tracerr"

	"rpdSix/commands"
	"rpdSix/helpers/extendeddiscord/extendeddiscordobjects"
	"rpdSix/helpers/extendeddiscord/extendeddiscordpermissions"
)

// Builds a function that checks whether the command caller meets your permission requirements
func Builder(
	callBack func(commands.CommandContext) error,
	requiredPermissions ...int,
) func(commands.CommandContext) error {

	return func(ctx commands.CommandContext) error {
		var authorMember, authorMemberErr = ctx.Message.AuthorMember()
		if authorMemberErr != nil {
			return tracerr.Wrap(authorMemberErr)
		}
		var extendedAuthorMember = extendeddiscordobjects.ExtendMember(authorMember, ctx.Session)

		extendedAuthorMember.GuildID = ctx.Message.GuildID // fix

		var hasAllPermissions, hasAllPermissionsErr = extendedAuthorMember.HasAllPermissions(requiredPermissions...)

		if hasAllPermissionsErr != nil {
			return tracerr.Wrap(hasAllPermissionsErr)
		}

		if hasAllPermissions {
			return callBack(ctx)
		} else {

			var requiredPermissionNames []string

			for _, permission := range requiredPermissions {
				var permissionName = extendeddiscordpermissions.ValueWithName[permission]
				requiredPermissionNames = append(requiredPermissionNames, permissionName)
			}

			return tracerr.Wrap(errors.New(fmt.Sprint(
				"permissions error, author does not have required permissions\n",
				"required permissions are: ", requiredPermissionNames)))
		}
	}
}
