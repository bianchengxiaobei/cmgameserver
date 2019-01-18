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
					room := &message.Room{}
					if room.Owner == nil {
						room.Owner = &message.RoomOwner{}
					}
					room.Owner.RoomId = k
					room.Owner.RoomName = v.GetRoomName()
					room.Owner.RoomOwnerId = v.GetRoomOwnerId()
					room.Owner.RoomOwnerName = handler.GameServer.GetRoleManager().GetOnlineRole(v.GetRoomOwnerId()).GetNickName()
					room.Owner.MaxPlayerNum = v.GetRoomMaxPlayerNum()
					room.Owner.CurPlayeNum = v.GetCurPlayerNum()
					room.Owner.GameType = v.GetGameType()
					room.Owner.MapId = v.GetMapId()
					room.Owner.RoomOwnerGroupId = v.GetRoomOwnerGroupId()
					room.Owner.RoomOwnerAvatarId = v.GetRoomOwnerAvatarId()
					room.Owner.IsWarFow = v.GetIsWarFow()
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
								role := handler.GameServer.GetRoleManager().GetOnlineRole(mem)
								if role != nil {
									msgMem := &message.RoomMember{}
									msgMem.GroupId = groupId
									msgMem.JoinerLevel = role.GetLevel()
									msgMem.JoinerIconId = role.GetAvatarId()
									msgMem.JoinerName = role.GetNickName()
									msgMem.JoinerId = mem
									msgMem.Ready = ready
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
