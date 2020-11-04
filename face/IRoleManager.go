package face

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/bean"
)

type IRoleManager interface {
	GetOnlineRole(roleId int64) IOnlineRole
	AddOnlineRole(role IOnlineRole)
	NewOnlineRole(roleId int64) IOnlineRole
	GetAllOnlineRole(gateSession network.SocketSessionInterface)map[int64]IOnlineRole
	RoleQuit(roleId int64)
	AllRoleQuit()
	SendEmailToAllRole(title string,content string,awardList []int32)bool
	SendRollInfoToAllRole(notify string)bool
	SendRollInfoToAllRoleGetItem(itemId int32,roleName string)
	SendEmailToOneRole(roleId int64, title string,content string,awardList []int32)bool
	BuyCard(roleId int64,cardType int32)bool
	BuyDiam(roleId int64,diamType int32)bool
	GetOnlineNum()int
	BuyDiamIOS(roleId int64,diamType int32,tranId string)bool
	BuyCardIOS(roleId int64,cardType int32,tranId string)bool
	SendCardDiamToAllRole()
	SendSaiJiAward()
	AddRoleHero(role IOnlineRole,hero bean.Hero)
	AddRoleExp(role IOnlineRole,expValue int32)
}