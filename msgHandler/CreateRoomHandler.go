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
				room.SetIsWarFow(protoMsg.Room.IsWarFow)
				room.SetMapId(protoMsg.Room.MapId)
				rMsg := &message.M2C_JoinRoom{}
				if rMsg.RoomSetting == nil{
					rMsg.RoomSetting = &message.RoomSetting{}
				}
				ownerRole := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
				if rMsg.Member == nil{
					rMsg.Member = &message.RoomMember{}
				}
				rMsg.Member.JoinerId = roleId
				rMsg.Member.JoinerName = ownerRole.GetNickName()
				rMsg.Member.JoinerIconId = ownerRole.GetAvatarId()
				rMsg.Member.GroupId = room.GetRoomOwnerGroupId()
				rMsg.RoomId = room.GetRoomId()
				rMsg.RoomSetting.IsWarFow = protoMsg.Room.IsWarFow
				rMsg.RoomSetting.MapId = protoMsg.Room.MapId
				handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5004,rMsg)
			}
		}
	}
}
