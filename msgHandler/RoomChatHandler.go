package msgHandler

import (
	"cmgameserver/message"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
)

type RoomChatHandler struct {
	GameServer face.IGameServer
}

func (handler *RoomChatHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_Chat); ok {
			roleId := innerMsg.RoleId
			role := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
			if role != nil{
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room != nil{
					roleIds := room.GetRoomRoleIds()
					for _,v := range roleIds{
						if v > 0{
							role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
							if role != nil && role.IsConnected(){
								handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5032,protoMsg)
							}
						}
					}
				}
			}
		}
	}
}