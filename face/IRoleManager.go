package face

import "github.com/bianchengxiaobei/cmgo/network"

type IRoleManager interface {
	GetOnlineRole(roleId int64) IOnlineRole
	AddOnlineRole(role IOnlineRole)
	NewOnlineRole(roleId int64) IOnlineRole
	GetAllOnlineRole(gateSession network.SocketSessionInterface)map[int64]IOnlineRole
	RoleQuit(roleId int64)
}