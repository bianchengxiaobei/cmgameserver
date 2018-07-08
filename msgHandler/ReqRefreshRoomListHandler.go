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
				RoomList:make([]*message.Room,0),
			}
			roomMap := handler.GameServer.GetRoomManager().GetAllRoom()
			if len(roomMap) > 0{
				for k,v := range roomMap{
					room := &message.Room{
						RoomId:k,
						RoomName:v.GetRoomName(),
						RoomOwnerId:v.GetRoomOwnerId(),
						RoomOwnerName:handler.GameServer.GetRoleManager().GetOnlineRole(v.GetRoomOwnerId()).GetNickName(),
						MaxPlayerNum:v.GetRoomMaxPlayerNum(),
						CurPlayeNum:v.GetCurPlayerNum(),
						GameType:v.GetGameType(),
						MapId:v.GetMapId(),
						RoomOwnerGroupId:v.GetRoomOwnerGroupId(),
					}
					rMsg.RoomList = append(rMsg.RoomList, room)
				}
			}
			handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5002,rMsg)
		}
	}
}
