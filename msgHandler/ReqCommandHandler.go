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
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			battleId := role.GetBattleId()
			if battleId > 0{
				battle := handler.GameServer.GetBattleManager().GetBattle(battleId)
				if battle != nil{
					if protoMsg.Cmd.CommandType == 26{
						battle.RemovePlayer(innerMsg.RoleId)
						//不移除房间，移除role的房间信息
						role.QuitBattle()
						return
					}else if protoMsg.Cmd.CommandType == 27{
						//删除战斗
						roomId := role.GetRoomId()
						handler.GameServer.GetRoomManager().DeleteRoom(roomId)
						battle.Finish()
						log4g.Infof("战斗移除[%d]",battleId)
					}
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
