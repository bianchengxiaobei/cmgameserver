package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ReqJoinRoomHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqJoinRoomHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ReqJoinRoom); ok {
			room := handler.GameServer.GetRoomManager().GetRoomByRoomId(protoMsg.RoomId)
			rMsg := &message.M2C_JoinRoom{}
			if room != nil{
				if rMsg.Room == nil{
					rMsg.Room = &message.Room{}
					if rMsg.Room.Owner == nil {
						rMsg.Room.Owner = &message.RoomOwner{}
						rMsg.Room.Owner.RoomId = -1//默认不存在房间
					}
					rMsg.Member = &message.RoomMember{}
				}
				roleId := innerMsg.RoleId
				if groupId,ok := room.JoinOneMember(roleId);ok == true{
					role := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
					role.SetRoomId(room.GetRoomId())
					role.SetInRooming(true)
					//先房主
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
					//成员
					//添加存在的队员
					if room.GetCurPlayerNum() > 1 {
						if rMsg.Room.Members == nil {
							rMsg.Room.Members= make([]*message.RoomMember, 0)
						}
						roleIds := room.GetRoomRoleIds()
						for _, mem := range roleIds {
							if mem > 0 {
								if mem == room.GetRoomOwnerId(){
									continue
								}
								groupId := room.GetRoomMemberGroupId(mem)
								ready := room.GetRoomMemberReady(mem)
								cityId := room.GetRoomMemberCityId(mem)
								role := handler.GameServer.GetRoleManager().GetOnlineRole(mem)
								if role != nil {
									msgMem := &message.RoomMember{}
									msgMem.GroupId = groupId
									msgMem.JoinerLevel = role.GetLevel()
									msgMem.JoinerIconId = role.GetAvatarId()
									msgMem.JoinerName = role.GetNickName()
									msgMem.JoinerId = mem
									msgMem.CityId = cityId
									msgMem.Ready = ready
									msgMem.Arrower = role.GetSoldierData(0)
									msgMem.Daodun = role.GetSoldierData(1)
									msgMem.Spear = role.GetSoldierData(2)
									msgMem.Fashi = role.GetSoldierData(3)
									rMsg.Room.Members = append(rMsg.Room.Members, msgMem)
								}
							}
						}
					}
					//通知房间内的客户端其他成员包括自己加入通知
					rMsg.Member.JoinerId = roleId
					rMsg.Member.JoinerName = role.GetNickName()
					rMsg.Member.JoinerIconId = role.GetAvatarId()
					rMsg.Member.GroupId = groupId
					rMsg.Member.CityId = 1
					rMsg.Member.JoinerLevel = role.GetLevel()
					rMsg.Member.Arrower = role.GetSoldierData(0)
					rMsg.Member.Daodun = role.GetSoldierData(1)
					rMsg.Member.Spear = role.GetSoldierData(2)
					rMsg.Member.Fashi = role.GetSoldierData(3)
					allRoomRoles := room.GetRoomRoleIds()
					for _,v := range allRoomRoles{
						if v > 0{
							role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
							if role != nil && role.IsConnected(){
								handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5004,rMsg)
							}
						}
					}
				}else{
					//加入失败
					role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
					rMsg.Room.Owner.RoomId = -1000//默认不存在房间
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5004,rMsg)
				}
			}else{
				//加入失败
				role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5004,rMsg)
			}
		}
	}
}
