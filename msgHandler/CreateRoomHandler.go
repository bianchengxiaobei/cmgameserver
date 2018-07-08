package msgHandler

import (
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/face"
)

type CreateRoomHandler struct {
	GameServer face.IGameServer
}

func (handler *CreateRoomHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_CreateRoom); ok {
			roleId := innerMsg.RoleId
			room := handler.GameServer.GetRoomManager().CreateRoom(roleId)
			if room != nil{
				room.SetRoomName(protoMsg.Room.RoomName)
				room.SetRoomMaxPlayerNum(protoMsg.Room.MaxPlayerNum)
				room.SetGameType(protoMsg.Room.GameType)
				rMsg := &message.M2C_JoinRoom{}
				rMsg.JoinerId = roleId
				rMsg.JoinerName = handler.GameServer.GetRoleManager().GetOnlineRole(roleId).GetNickName()
				rMsg.JoinerIconId = 0
				rMsg.GroupId = room.GetRoomOwnerGroupId()
				rMsg.RoomId = room.GetRoomId()
				handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5004,rMsg)
			}
		}
	}
}
