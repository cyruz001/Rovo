package permission

type Permission string

const (
	PermissionRead   Permission = "READ"
	PermissionWrite  Permission = "WRITE"
	PermissionEdit   Permission = "EDIT"
	PermissionAdmin  Permission = "ADMIN"
	PermissionDelete Permission = "DELETE"
)

var rolePermissions = map[string][]Permission{
	"USER":  {PermissionRead, PermissionWrite, PermissionEdit},
	"ADMIN": {PermissionRead, PermissionWrite, PermissionAdmin},
}
