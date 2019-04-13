package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ReadyHandler struct {
	GameServer face.IGameServer
}

func (handler *ReadyHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ReqReady); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room != nil{
					if room.SetRoomMemberReady(protoMsg.Ready,innerMsg.RoleId) == true{
						//发送准备消息，房间内玩家
						allRoomRoles := room.GetRoomRoleIds()
						rMsg := new(message.M2C_ReadySuccess)
						rMsg.RoleId = role.GetRoleId()
						rMsg.Ready = protoMsg.Ready
						for _,v := range allRoomRoles{
							if v > 0{
								role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
								if role != nil && role.IsConnected(){
									handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5007,rMsg)
								}
							}
						}
					}
				}
			}
		}
	}
}
