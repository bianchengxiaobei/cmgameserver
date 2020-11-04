package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeSexHandler struct {
	GameServer face.IGameServer
}

func (handler *ChangeSexHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ChangeSex); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetSex(protoMsg.Sex)
				msg := new(message.M2C_ChangeSexResult)
				msg.Sex = role.GetSex()
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5041, msg)
				//dbSession := handler.GameServer.GetDBManager().Get()
				//if dbSession != nil {
				//	//更新角色昵称
				//	c := dbSession.DB("sanguozhizhan").C("Role")
				//	data := bson.M{"$set": bson.M{"sex": protoMsg.Sex}}
				//	err := c.Update(bson.M{"roleid": innerMsg.RoleId}, data)
				//	if err != nil {
				//		log4g.Errorf("更新Sex出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
				//		return
				//	}
				//
				//}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_ChangeNickName！")
		}
	}
}
