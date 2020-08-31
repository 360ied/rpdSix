package extendeddiscordpermissions

func HasPermission(permissionsInteger, permission int) bool {
	return permissionsInteger&permission == permission
}

func HasAllPermissions(permissionInteger int, permissions ...int) bool {
	for _, permission := range permissions {
		if !HasPermission(permissionInteger, permission) {
			return false
		}
	}
	return true
}
