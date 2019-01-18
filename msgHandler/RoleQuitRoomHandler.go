package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type RoleQuitRoomHandler struct {
	GameServer 		face.IGameServer
}

func (handler *RoleQuitRoomHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_QuitRoom); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			roomManager := handler.GameServer.GetRoomManager()
			room := roomManager.GetRoomByRoomId(role.GetRoomId())
			roomId := room.GetRoomId()
			if room != nil{
				//如果是房主退出就地解散
				if room.IsRoomOwner(role.GetRoleId()){
					if roomManager.DeleteRoom(roomId){
						log4g.Infof("删除房间[%d]成功!",roomId)
					}else{
						log4g.Errorf("删除房间[%d]失败!",roomId)
					}
				}else{
					if roomManager.RemoveOneMemberByRoom(room,role) == false{
						log4g.Errorf("移除房间[%d]成员[%d]失败!",roomId,role.GetRoleId())
					}
				}
			}
		}
	}
}