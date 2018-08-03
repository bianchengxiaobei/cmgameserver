package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"cmgameserver/face"
)

type RoleRegisterGateHandler struct {
	GameServer face.IGameServer
}

func (handler *RoleRegisterGateHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg,ok:=msg.(network.InnerWriteMessage);ok {
		if protoMsg,ok := innerMsg.MsgData.(*message.G2M_RoleRegisterToGateSuccess);ok{
			rMsg := new(message.M2C_EnterLobby)
			if rMsg.RoleBasicInfo == nil{
				rMsg.RoleBasicInfo = new(message.Role)
			}
			rMsg.IsInBattle = false
			rMsg.RoleBasicInfo.RoleId = protoMsg.RoleId
			role := handler.GameServer.GetRoleManager().GetOnlineRole(protoMsg.RoleId)
			rMsg.RoleBasicInfo.NickName = role.GetNickName()
			rMsg.RoleBasicInfo.Gold = role.GetGold()
			rMsg.RoleBasicInfo.AvatarId = role.GetAvatarId()
			rMsg.RoleBasicInfo.Exp = role.GetExp()
			rMsg.RoleBasicInfo.Diam = role.GetDiam()
			rMsg.RoleBasicInfo.Level = role.GetLevel()
			heroMap := role.GetAllHero()
			if len(heroMap) > 0{
				for _,v := range heroMap{
					hero := &message.Hero{}
					hero.HeroId = v.HeroId
					hero.Level = v.Level
					for _,v1 := range v.ItemIds{
						if v1 > 0{
							hero.ItemIds = append(hero.ItemIds, v1)
						}
					}
					rMsg.HeroInfo = append(rMsg.HeroInfo, hero)
				}
			}
			handler.GameServer.WriteInnerMsg(session,protoMsg.RoleId,5000,rMsg)
		}
	}
}

