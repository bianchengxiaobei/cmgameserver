package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
)

type BattleLoadFinishedHandler struct {
	GameServer face.IGameServer
}

func (handler *BattleLoadFinishedHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_LoadFinished); ok {
			//检测房间内所有玩家是否加载完成
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				role.SetLoadFinished(true)
				if handler.GameServer.GetRoomManager().CheckAllRoomMemberLoadFinished(role.GetRoomId()){
					rMsg := new(message.M2C_StartBattle)
					battle := handler.GameServer.GetBattleManager().GetBattleInFree()
					//开始战斗（发送给所有玩家，关闭加载界面）
					room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
					roleIds := room.GetRoomRoleIds()
					var roles [4]face.IOnlineRole
					if battle != nil{
						index := 0
						rMsg.BattleId = battle.GetBattleId()
						for _,v := range roleIds{
							if v > 0{
								role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
								if role != nil{
									roles[index] = role
									index++
									role.SetInRooming(false)
									role.SetInBattling(true)
									role.SetBattleId(battle.GetBattleId())
									if  role.IsConnected(){
										handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5011,rMsg)
									}
								}
							}
						}
						battle.ReStart(&roles,face.FreeRoomBattleType)
					}else{
						index := 0
						for _,v := range roleIds{
							if v > 0{
								role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
								if role != nil{
									roles[index] = role
									index++
									role.SetInRooming(false)
									role.SetInBattling(true)
								}
							}
						}
						battle := handler.GameServer.GetBattleManager().CreateBattle(&roles,face.FreeRoomBattleType)
						if battle != nil{
							battle.Start()
							rMsg.BattleId = battle.GetBattleId()
						}
						for _,v := range roles{
							if v != nil && v.IsConnected(){
								v.SetBattleId(battle.GetBattleId())
								handler.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5011,rMsg)
							}
						}
					}
				}
			}
		}
	}
}
