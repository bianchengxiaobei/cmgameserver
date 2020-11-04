package msgHandler

import (
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
)

type InviteRoomHandler struct {
	GameServer face.IGameServer
}

func (handler *InviteRoomHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_InviteRoom); ok {
			msg := new(message.M2C_InviteRoomResult)
			msg.RoomId = -1
			role := handler.GameServer.GetRoleManager().GetOnlineRole(protoMsg.RoleId)
			my := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			room := handler.GameServer.GetRoomManager().GetRoomByRoomId(protoMsg.RoomId)
			if role != nil && role.IsConnected() && room != nil && room.GetInBattle() == false{
				if role.GetInSimulateBattle() || role.IsInRooming() || role.IsInBattling(){
					handler.GameServer.WriteInnerMsg(my.GetGateSession(), my.GetRoleId(), 5105, msg)
					return
				}else{
					matchPlayer := role.GetMatchPlayer()
					if matchPlayer != nil && (matchPlayer.GetBInMatching() || matchPlayer.GetBInBattleRoom()){
						handler.GameServer.WriteInnerMsg(my.GetGateSession(), my.GetRoleId(), 5105, msg)
						return
					}
				}
				msg.RoomId = room.GetRoomId()
				msg.Name = my.GetNickName()
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5105, msg)
			} else {
				handler.GameServer.WriteInnerMsg(my.GetGateSession(), my.GetRoleId(), 5105, msg)
			}
		} else {
			log4g.Error("InviteRoomHandler！")
		}
	}
}
