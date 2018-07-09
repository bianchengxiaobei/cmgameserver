package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ReqCommandHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqCommandHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_Command); ok {
			battleId := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId).GetBattleId()
			if battleId > 0{
				battle := handler.GameServer.GetBattleManager().GetBattle(battleId)
				if battle != nil{
					battle.AddFrameCommand(protoMsg.Cmd.PlayerId,protoMsg.Cmd.CommandType,protoMsg.Cmd.Param)
				}else{
					log4g.Infof("不存在该战斗:[%d]",battleId)
				}
			}else{
				log4g.Infof("战斗Id无效:[%d]",battleId)
			}
		}
	}
}
