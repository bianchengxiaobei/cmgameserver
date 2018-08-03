package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type ReqStartBattleHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqStartBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_StartBattle); ok {
			//判断是否可以开始游戏
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room != nil{
					if room.CheckRoomReady() == true{
						//发送给所有客户端开始游戏
						rMsg := new(message.M2C_StartBattleLoad)
						rMsg.AllReady = true
						allRoomRoles := room.GetRoomRoleIds()
						for _,v := range allRoomRoles{
							if v > 0{
								role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
								if role != nil && role.IsConnected(){
									role.SetInRooming(false)
									role.SetInBattling(true)
									handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5009,rMsg)
								}
							}
						}
					}else{
						//发送给客户端有玩家没有准备
						rMsg := new(message.M2C_StartBattleLoad)
						rMsg.AllReady = false
						handler.GameServer.WriteInnerMsg(session,innerMsg.RoleId,5009,rMsg)
					}
				}
			}
		}
	}

}
