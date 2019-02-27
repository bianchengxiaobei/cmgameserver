package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"gopkg.in/mgo.v2/bson"
)

type WinBattleHandler struct {
	GameServer face.IGameServer
}

func (handler *WinBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_WinBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.AddGold(300)
				role.AddExp(100)
				role.WinLevel(protoMsg.BattleId)
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色经验和金钱
					c := dbSession.DB("sanguozhizhan").C("Role")
					data := bson.M{"$set": bson.M{"gold": role.GetGold(), "exp": role.GetExp()}}
					err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
					if err != nil {
						log4g.Errorf("更新AvatarId出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					//更新英雄经验
					//if len(protoMsg.HeroIds) > 0 {
					//	c := dbSession.DB("sanguozhizhan").C("Hero")
					//	for _, v := range protoMsg.HeroIds {
					//		exp, has := role.AddHeroExp(v, 50)
					//		if has {
					//			data := bson.M{"$set": bson.M{"exp": exp}}
					//			c.Update(bson.M{"roleid": innerMsg.RoleId, "heroid": v}, data)
					//		}
					//	}
					//}
					//获得宝箱
					msg := new(message.M2C_BattleResult)
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5019, msg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_WinBattle！")
		}
	}
}
