package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeRoomCityIdHandler struct {
	GameServer face.IGameServer
}

func (handler *ChangeRoomCityIdHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ChangeRoomCityId); ok {
			roleManager := handler.GameServer.GetRoleManager()
			role := roleManager.GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				if role.IsInRooming() == false{
					return
				}
				//如果已经在战斗中不能改变
				if role.IsInBattling(){
					return
				}
				if role.GetRoomId() != protoMsg.RoomId{
					return
				}
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(protoMsg.RoomId)
				if room == nil || room.GetInBattle(){
					return
				}
				if room.IsRoomOwner(protoMsg.RoleId){
					//如果是房主
					room.SetRoomOwnerCityId(protoMsg.CityId)
				}else{
					//如果不是则修改成员
					room.SetRoomMemberCityId(protoMsg.CityId,innerMsg.RoleId)
				}
				allRole := room.GetRoomRoleIds()
				for _,v := range allRole{
					if v > 0{
						temp := roleManager.GetOnlineRole(v)
						if temp != nil && temp.IsConnected(){
							handler.GameServer.WriteInnerMsg(temp.GetGateSession(),temp.GetRoleId(),5099,protoMsg)
						}
					}
				}
			}
		} else {
			log4g.Error("不是C2M_TeamStartMatch！")
		}
	}
}
