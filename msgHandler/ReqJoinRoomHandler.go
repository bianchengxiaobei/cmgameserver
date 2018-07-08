package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ReqjoinRoomHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqjoinRoomHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ReqJoinRoom); ok {
			room := handler.GameServer.GetRoomManager().GetRoomByRoomId(protoMsg.RoomId)
			if room != nil{
				roleId := innerMsg.RoleId
				if groupId,ok := room.JoinOneMember(roleId);ok == true{
					role := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
					role.SetRoomId(room.GetRoomId())
					//通知房间内的客户端其他成员包括自己加入通知
					rMsg := &message.M2C_JoinRoom{}
					rMsg.JoinerId = roleId
					rMsg.JoinerName = handler.GameServer.GetRoleManager().GetOnlineRole(roleId).GetNickName()
					rMsg.JoinerIconId = 0
					rMsg.GroupId = groupId
					rMsg.RoomId = room.GetRoomId()
					allRoomRoles := room.GetRoomRoleIds()
					for _,v := range allRoomRoles{
						if v > 0{
							role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
							if role != nil && role.IsConnected(){
								handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5004,rMsg)
							}
						}
					}
				}
			}
		}
	}
}
