package permission

type Permission string

const (
	PermissionView   Permission = "VIEW"
	PermissionRead   Permission = "READ"
	PermissionWrite  Permission = "WRITE"
	PermissionEdit   Permission = "EDIT"
	PermissionCreate Permission = "CREATE"
	PermissionAdd    Permission = "ADD"
	PermissionAdmin  Permission = "ADMIN"
	PermissionDelete Permission = "DELETE"
)

var rolePermissions = map[string][]Permission{
	"USER": {PermissionView, PermissionRead, PermissionWrite, PermissionEdit},
	"ADMIN": {PermissionView, PermissionRead, PermissionWrite, PermissionAdmin,
		PermissionDelete, PermissionCreate, PermissionView,
		PermissionAdd, PermissionEdit},
}
