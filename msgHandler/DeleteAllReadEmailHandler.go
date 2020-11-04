package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type DeleteAllReadEmailHandler struct {
	GameServer face.IGameServer
}

func (handler *DeleteAllReadEmailHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_DeleteAllReadEmail); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				for i:=0;i<role.GetMaxEmailNum();i++ {
					email := role.GetEmailByIndex(i)
					if &email != nil{
						if email.Get == false || email.Valid == false{
							continue
						}
						if role.DeleteEmail(i) == false{
							log4g.Infof("删除Index[%d]邮件Index[%d]出错",i,email.EmailIndex)
						}
					}
				}
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5094, protoMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_DeleteAllReadEmail！")
		}
	}
}