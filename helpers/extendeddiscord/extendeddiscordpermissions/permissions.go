package extendeddiscordpermissions

func IsPermittedAll(permissionInteger int, permissions ...int) bool {
	if permissionInteger&ADMINISTRATOR == ADMINISTRATOR {
		return true
	}
	for _, permission := range permissions {
		if permissionInteger&permission != permission {
			return false
		}
	}
	return true
}
