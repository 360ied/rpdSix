package extendeddiscordpermissions

func IsPermitted(permissionsInteger, permission int) bool {
	return permissionsInteger&permission == permission || permissionsInteger&ADMINISTRATOR == ADMINISTRATOR
}

func IsPermittedAll(permissionInteger int, permissions ...int) bool {
	if IsPermitted(permissionInteger, ADMINISTRATOR) {
		return true
	}
	for _, permission := range permissions {
		if !IsPermitted(permissionInteger, permission) {
			return false
		}
	}
	return true
}
