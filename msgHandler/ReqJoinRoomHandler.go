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
			if room != nil{
				roleId := innerMsg.RoleId
				if groupId,ok := room.JoinOneMember(roleId);ok == true{
					role := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
					role.SetRoomId(room.GetRoomId())
					role.SetInRooming(true)
					//通知房间内的客户端其他成员包括自己加入通知
					rMsg := &message.M2C_JoinRoom{}
					if rMsg.RoomSetting == nil{
						rMsg.RoomSetting = &message.RoomSetting{}
					}
					rMsg.Member = &message.RoomMember{}
					rMsg.Member.JoinerId = roleId
					rMsg.Member.JoinerName = role.GetNickName()
					rMsg.Member.JoinerIconId = role.GetAvatarId()
					rMsg.Member.GroupId = groupId
					rMsg.Member.CityId = 1
					rMsg.Member.JoinerLevel = role.GetLevel()
					rMsg.RoomId = room.GetRoomId()
					rMsg.RoomSetting.MapId = room.GetMapId()
					rMsg.RoomSetting.IsWarFow = room.GetIsWarFow()
					rMsg.RoomSetting.GameType = room.GetGameType()
					rMsg.RoomSetting.BOutsideMonster = room.GetIsOutsideMonster()
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
				}
			}
		}
	}
}
