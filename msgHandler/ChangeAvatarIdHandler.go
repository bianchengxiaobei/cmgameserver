package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"gopkg.in/mgo.v2/bson"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeAvatarIdHandler struct {
	GameServer 		face.IGameServer
}

func (handler *ChangeAvatarIdHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ChangeAvatarIcon); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetAvatarId(protoMsg.AvatarId)
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession != nil {
					//更新角色经验和金钱
					c := dbSession.DB("sanguozhizhan").C("Role")
					data := bson.M{"$set": bson.M{"avatarid": protoMsg.AvatarId}}
					err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
					if err != nil {
						log4g.Errorf("更新AvatarId出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
						return
					}
					msg := new(message.M2C_ChangeAvatarIcon)
					msg.AvatarId = role.GetAvatarId()
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5025, msg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_ChangeAvatarIcon！")
		}
	}
}
