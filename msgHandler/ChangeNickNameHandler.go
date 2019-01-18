package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"gopkg.in/mgo.v2/bson"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeNickNameHandler struct {
	GameServer face.IGameServer
}

func (handler *ChangeNickNameHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ChangeNickName); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetNickName(protoMsg.NickName)
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色昵称
					c := dbSession.DB("sanguozhizhan").C("Role")
					data := bson.M{"$set": bson.M{"nickname": protoMsg.NickName}}
					err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
					if err != nil {
						log4g.Errorf("更新NickName出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					msg := new(message.M2C_ChangeNickName)
					msg.ResultName = role.GetNickName()
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5024, msg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_ChangeNickName！")
		}
	}
}