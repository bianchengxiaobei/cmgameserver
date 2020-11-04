package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ForceEnterBattleHandler struct {
	GameServer face.IGameServer
}

func (handler *ForceEnterBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_ForceEnterBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				if role.GetBattleId() > 0{
					//说明已经在战斗
					return
				}
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room != nil{
					roleIds := room.GetRoomRoleIds()
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
					battle := handler.GameServer.GetBattleManager().GetBattleInFree()
					if battle != nil{
						battle.ReStart(&roles,face.FreeRoomBattleType)
					}else{
						battle := handler.GameServer.GetBattleManager().CreateBattle(&roles,face.FreeRoomBattleType)
						if battle != nil{
							battle.Start()
						}
					}
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_GetBoxAwardItem！")
		}
	}
}