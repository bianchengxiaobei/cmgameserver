package msgHandler

import (
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type LoginToGameServerHandler struct {
	GameServer IGameServer
}

func (handler *LoginToGameServerHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg,ok:=msg.(network.InnerWriteMessage);ok{
		if protoMsg, ok := innerMsg.MsgData.(*message.G2M_LoginToGameServer); ok {
			roleId := protoMsg.RoleId
			roleManager := handler.GameServer.GetRoleManager()
			onlineRole := roleManager.GetOnlineRole(roleId)
			if onlineRole == nil{
				onlineRole = roleManager.NewOnlineRole(roleId)
				if onlineRole == nil{
					log4g.Infof("数据库载入OnlineRole[%d]失败!",roleId)
					return
				}
			}
			//初始化在线角色
			onlineRole.GateId = protoMsg.GateId
			onlineRole.UserName = protoMsg.UserName

			roleManager.AddOnlineRole(onlineRole)
			//通知网关服务器登录游戏逻辑服成功
			rMsg := new(message.M2G_LoginSuccessNotifyGate)
			rMsg.RoleId = onlineRole.Role.RoleId
			rMsg.UserId = protoMsg.UserId
			rMsg.ServerId = protoMsg.ServerId
			handler.GameServer.WriteInnerMsg(session,0,10002,rMsg)
		}
	}
}
