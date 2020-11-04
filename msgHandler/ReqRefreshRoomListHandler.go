package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ReqRefreshRoomListHandle struct {
	GameServer face.IGameServer
}

func (handler *ReqRefreshRoomListHandle) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_ReqRefreshRoomList); ok {
			rMsg := &message.M2C_RefreshRoomList{
				RoomList: make([]*message.Room, 0),
			}
			roomMap := handler.GameServer.GetRoomManager().GetAllRoom()
			if len(roomMap) > 0 {
				for k, v := range roomMap {
					if v.GetInBattle(){
						continue
					}
					room := &message.Room{}
					if room.Owner == nil {
						room.Owner = &message.RoomOwner{}
					}
					room.Owner.RoomId = k
					room.Owner.GameType = v.GetGameType()
					room.Owner.RoomOwnerId = v.GetRoomOwnerId()
					room.Owner.RoomOwnerName = handler.GameServer.GetRoleManager().GetOnlineRole(v.GetRoomOwnerId()).GetNickName()
					room.Owner.MaxPlayerNum = v.GetRoomMaxPlayerNum()
					room.Owner.CurPlayeNum = v.GetCurPlayerNum()
					room.Owner.GameType = v.GetGameType()
					room.Owner.IsOutsideMonster = v.GetIsOutsideMonster()
					room.Owner.MapId = v.GetMapId()
					room.Owner.RoomOwnerGroupId = v.GetRoomOwnerGroupId()
					room.Owner.RoomOwnerAvatarId = v.GetRoomOwnerAvatarId()
					room.Owner.RoomOwnerLevel = v.GetRoomOwnerLevel()
					room.Owner.CityId = v.GetRoomOwnerCityId()
					room.Owner.IsWarFow = v.GetIsWarFow()
					room.Owner.Arrower = v.GetArrowerData()
					room.Owner.Daodun = v.GetDaodunData()
					room.Owner.Fashi = v.GetFashiData()
					room.Owner.Spear = v.GetSpearData()
					//添加存在的队员
					if v.GetCurPlayerNum() > 1 {
						if room.Members == nil {
							room.Members = make([]*message.RoomMember, 0)
						}
						roleIds := v.GetRoomRoleIds()
						for _, mem := range roleIds {
							if mem > 0 {
								if mem == v.GetRoomOwnerId(){
									continue
								}
								groupId := v.GetRoomMemberGroupId(mem)
								ready := v.GetRoomMemberReady(mem)
								cityId := v.GetRoomMemberCityId(mem)
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
									room.Members = append(room.Members, msgMem)
								}
							}
						}
					}
					rMsg.RoomList = append(rMsg.RoomList, room)
				}
			}

			handler.GameServer.WriteInnerMsg(session, innerMsg.RoleId, 5002, rMsg)
		}
	}
}
