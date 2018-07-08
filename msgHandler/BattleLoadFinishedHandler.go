package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type BattleLoadFinsihedHandler struct {
	GameServer face.IGameServer
}

func (handler *BattleLoadFinsihedHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_LoadFinished); ok {
			//检测房间内所有玩家是否加载完成
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				role.SetLoadFinished(true)
				if handler.GameServer.GetRoomManager().CheckAllRoomMemberLoadFinished(role.GetRoomId()){
					//开始战斗（发送给所有玩家，关闭加载界面）
					roleIds := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId()).GetRoomRoleIds()
					rMsg := new(message.M2C_StartBattle)
					var roles [4]face.IOnlineRole
					index := 0
					for _,v := range roleIds{
						if v > 0{
							role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
							if role != nil{
								roles[index] = role
								index++
								if  role.IsConnected(){
									handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5011,rMsg)
								}
							}
						}
					}
					battle := handler.GameServer.GetBattleManager().CreateBattle(roles)
					if battle != nil{
						battle.Start()
					}
				}
			}
		}
	}
}
