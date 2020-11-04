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
			ownerRole := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
			if ownerRole.IsInRooming(){
				//重复创建房间
				return
			}
			room := handler.GameServer.GetRoomManager().CreateRoom(roleId)
			if room != nil{
				//room.SetRoomName(ownerRole.GetNickName() + "的房间")
				room.SetRoomMaxPlayerNum(protoMsg.Room.MaxPlayerNum)
				room.SetGameType(protoMsg.Room.GameType)
				room.SetIsWarFow(protoMsg.Room.IsWarFow)
				room.SetIsOutsideMonster(protoMsg.Room.IsOutsideMonster)
				room.SetRoomOwnerCityId(protoMsg.Room.CityId)
				room.SetMapId(protoMsg.Room.MapId)
				rMsg := &message.M2C_JoinRoom{}
				if rMsg.Room == nil{
					rMsg.Room = &message.Room{}
					if rMsg.Room.Owner == nil{
						rMsg.Room.Owner = &message.RoomOwner{}
					}
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
				rMsg.Room.Owner.RoomId = room.GetRoomId()
				rMsg.Room.Owner.GameType = room.GetGameType()
				rMsg.Room.Owner.RoomOwnerId = room.GetRoomOwnerId()
				rMsg.Room.Owner.RoomOwnerName = room.GetRoomOwnerName()
				rMsg.Room.Owner.MaxPlayerNum = room.GetRoomMaxPlayerNum()
				rMsg.Room.Owner.CurPlayeNum = room.GetCurPlayerNum()
				rMsg.Room.Owner.GameType = room.GetGameType()
				rMsg.Room.Owner.IsOutsideMonster = room.GetIsOutsideMonster()
				rMsg.Room.Owner.MapId = room.GetMapId()
				rMsg.Room.Owner.RoomOwnerGroupId = room.GetRoomOwnerGroupId()
				rMsg.Room.Owner.RoomOwnerAvatarId = room.GetRoomOwnerAvatarId()
				rMsg.Room.Owner.RoomOwnerLevel = room.GetRoomOwnerLevel()
				rMsg.Room.Owner.CityId = room.GetRoomOwnerCityId()
				rMsg.Room.Owner.IsWarFow = room.GetIsWarFow()
				rMsg.Room.Owner.Arrower = room.GetArrowerData()
				rMsg.Room.Owner.Daodun = room.GetDaodunData()
				rMsg.Room.Owner.Fashi = room.GetFashiData()
				rMsg.Room.Owner.Spear = room.GetSpearData()
				rMsg.Member.Arrower = ownerRole.GetSoldierData(0)
				rMsg.Member.Daodun = ownerRole.GetSoldierData(1)
				rMsg.Member.Spear = ownerRole.GetSoldierData(2)
				rMsg.Member.Fashi = ownerRole.GetSoldierData(3)
				handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5004,rMsg)
			}
		}
	}
}
