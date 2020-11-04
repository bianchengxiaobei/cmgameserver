package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"gopkg.in/mgo.v2/bson"
)

type ChangePasswordHandler struct {
	GameServer face.IGameServer
}

func (handler *ChangePasswordHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ChangePassword); ok {
			roleManager := handler.GameServer.GetRoleManager()
			role := roleManager.GetOnlineRole(innerMsg.RoleId)
			msg := new(message.M2C_ChangePasswordResult)
			msg.Result = -1
			msg.Pass = protoMsg.Password
			if role != nil {
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色
					c := dbSession.DB("sanguozhizhan").C("User")
					err := c.Update(bson.M{"userid": role.GetUserId()},
					bson.M{"$set":bson.M{"password":protoMsg.Password}})
					if err != nil {
						log4g.Infof("更新Role出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					msg.Result = 1
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5110, msg)
				} else {
					log4g.Info("dbSession == nil")
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5110, msg)
				}
			}
		} else {
			log4g.Error("不是C2M_TeamStartMatch！")
		}
	}
}