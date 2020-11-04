package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
)

type ReStartPauseBattleHandler struct {
	GameServer face.IGameServer
}
///从暂停中恢复比赛
func (handler *ReStartPauseBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ReStartPauseBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				battleId := role.GetBattleId()
				if battleId != protoMsg.BattleId{
return
				}
				battle := handler.GameServer.GetBattleManager().GetBattle(battleId)
				if battle != nil{
					battle.RestartFormPause(protoMsg)
				}else{
					log4g.Error("No Battle")
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_WinBattle！")
		}
	}
}
