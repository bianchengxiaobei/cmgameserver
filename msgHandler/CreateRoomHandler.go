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
			ownerRole := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
			if room != nil{
				//room.SetRoomName(ownerRole.GetNickName() + "的房间")
				room.SetRoomMaxPlayerNum(protoMsg.Room.MaxPlayerNum)
				room.SetGameType(protoMsg.Room.GameType)
				room.SetIsWarFow(protoMsg.Room.IsWarFow)
				room.SetIsOutsideMonster(protoMsg.Room.IsOutsideMonster)
				room.SetRoomOwnerCityId(protoMsg.Room.CityId)
				room.SetMapId(protoMsg.Room.MapId)
				rMsg := &message.M2C_JoinRoom{}
				if rMsg.RoomSetting == nil{
					rMsg.RoomSetting = &message.RoomSetting{}
				}

				if rMsg.Member == nil{
					rMsg.Member = &message.RoomMember{}
				}
				rMsg.Member.JoinerId = roleId
				rMsg.Member.JoinerName = ownerRole.GetNickName()
				rMsg.Member.JoinerIconId = ownerRole.GetAvatarId()
				rMsg.Member.JoinerLevel = ownerRole.GetLevel()
				rMsg.Member.CityId = protoMsg.Room.CityId
				rMsg.Member.GroupId = room.GetRoomOwnerGroupId()
				rMsg.RoomId = room.GetRoomId()
				rMsg.RoomSetting.IsWarFow = protoMsg.Room.IsWarFow
				rMsg.RoomSetting.MapId = protoMsg.Room.MapId
				rMsg.RoomSetting.GameType = protoMsg.Room.GameType
				rMsg.RoomSetting.BOutsideMonster = protoMsg.Room.IsOutsideMonster
				rMsg.Member.Arrower = ownerRole.GetSoldierData(0)
				rMsg.Member.Daodun = ownerRole.GetSoldierData(1)
				rMsg.Member.Spear = ownerRole.GetSoldierData(2)
				rMsg.Member.Fashi = ownerRole.GetSoldierData(3)
				handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5004,rMsg)
			}
		}
	}
}
