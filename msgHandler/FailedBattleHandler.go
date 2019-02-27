package msgHandler

import (
	"cmgameserver/face"
	"gopkg.in/mgo.v2/bson"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
)

type FailedBattleHandler struct {
	GameServer face.IGameServer
}
func (handler *FailedBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_FailedBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.AddExp(10)
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色经验和金钱
					c := dbSession.DB("sanguozhizhan").C("Role")
					data := bson.M{"$set": bson.M{"exp": role.GetExp()}}
					err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
					if err != nil {
						log4g.Errorf("更新s出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					msg := new(message.M2C_BattleResult)
					msg.Award1 = 1
					msg.Award2 = 2
					msg.Award3 = 3
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5019, msg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_FailedBattle！")
		}
	}
}
