package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type BuyHeroHandler struct {
	GameServer face.IGameServer
}
func (handler *BuyHeroHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_BuyHero); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				rMsg := &message.M2C_BuyHeroResult{}
				if role.BuyHero(protoMsg.HeroId) == false{
					rMsg.HeroId = protoMsg.HeroId
					rMsg.ResultCode = 0
				}else{
					dbSession := handler.GameServer.GetDBManager().Get()
					if dbSession != nil {
						//更新角色经验和金钱
						c := dbSession.DB("sanguozhizhan").C("Hero")
						hero := role.GetHero(protoMsg.HeroId)
						err := c.Insert(hero)
						if err != nil {
							log4g.Errorf("插入英雄出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
							rMsg.HeroId = protoMsg.HeroId
							rMsg.ResultCode = 0
						}else{
							rMsg.HeroId = protoMsg.HeroId
							rMsg.ResultCode = 1
						}
						handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5031, rMsg)
					}
				}

			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_ChangeAvatarIcon！")
		}
	}
}
