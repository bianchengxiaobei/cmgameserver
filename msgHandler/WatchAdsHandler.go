package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"gopkg.in/mgo.v2/bson"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type WatchAdsHandler struct {
	GameServer face.IGameServer
}

func (handler *WatchAdsHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_WatchAds); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				var addGold int32
				if protoMsg.IsBanner{
					addGold = 20
					role.AddGold(20)
				}else{
					addGold = 50
					role.AddGold(50)
				}
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色经验和金钱
					c := dbSession.DB("sanguozhizhan").C("Role")
					data := bson.M{"$set": bson.M{"gold": role.GetGold()}}
					err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
					if err != nil {
						log4g.Errorf("更新出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					msg := new(message.M2C_WatchAdsResult)
					msg.Gold = addGold
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5021, msg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_WinBattle！")
		}
	}
}