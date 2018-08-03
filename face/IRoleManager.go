package face

type IRoleManager interface {
	GetOnlineRole(roleId int64) IOnlineRole
	AddOnlineRole(role IOnlineRole)
	NewOnlineRole(roleId int64) IOnlineRole
	GetAllOnlineRole()map[int64]IOnlineRole
}