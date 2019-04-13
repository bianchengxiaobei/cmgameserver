package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeFreeSoldierDataHandler struct {
	GameServer face.IGameServer
}
func (handler *ChangeFreeSoldierDataHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ChangeFreeSoldierData); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				if role.ChangeFreeSoldierData(int(protoMsg.PlayerTypeIndex),*protoMsg.Data){
					//回送
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5039, protoMsg)
				}else{
					log4g.Infof("更改自由士兵出错:%d",innerMsg.RoleId)
				}
			} else {
				log4g.Infof("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Info("不是C2M2C_ChangeFreeSoldierData！")
		}
	}
}
