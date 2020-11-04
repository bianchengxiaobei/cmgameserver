package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/tool"
)

type UpgradeHeroLevelHandler struct {
	GameServer face.IGameServer
}


func (handler *UpgradeHeroLevelHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_UpgradeHeroLevel); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				hero := role.GetHero(protoMsg.HeroId)
				msg := new(message.M2C_UpgradeHeroLevelResult)
				if hero == nil{
					msg.HeroId = -1
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5091,msg)
				return
					}
				//升级需要的金币
				needGold := tool.GetHeroUpgradeNeedGold(hero.Level)
				if needGold == 0{
					println("22222")
				}
				if role.GetGold() < needGold{
					msg.HeroId = -1
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5091,msg)
				return
					}
				item := role.GetItem(protoMsg.ItemIndex)
				if item.ItemId != protoMsg.ItemId || item.ItemNum < 10{
					msg.HeroId = -1
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5091,msg)
					return
				}
				role.AddGold(-needGold)
				level,skillPoint,ok := role.UpgradeHeroLevel(protoMsg.HeroId)
				if ok == false{
					msg.HeroId = -1
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5091,msg)
					return
				}
				item.ItemNum -= 10
				if item.ItemNum == 0{
					item.Clear()
				}
				role.SetItem(protoMsg.ItemIndex,item)
				msg.HeroId = protoMsg.HeroId
				msg.HeroLevel = level
				msg.SkillPoint = skillPoint
				msg.ItemIndex = protoMsg.ItemIndex
				msg.ItemNum = item.ItemNum
				msg.ItemId = item.ItemId
				msg.Gold = role.GetGold()
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5091,msg)
			}
		} else {
			log4g.Error("不是C2M_UpgradeHeroLevel！")
		}
	}
}
