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
					//失败
					if protoMsg.Cmd.CommandType == 26{
						//该玩家被占领了
						battle.RemovePlayer(innerMsg.RoleId)
						//不移除房间，移除role的房间信息
						role.QuitBattle()
						return
					}else if protoMsg.Cmd.CommandType == 27{
						//胜利
						if battle.GetBattleType() == face.FreeRoomBattleType{
							//删除战斗
							roomId := role.GetRoomId()
							if roomId <= 0{
								log4g.Info("房间id为0")
								return
							}
							handler.GameServer.GetRoomManager().DeleteRoom(roomId)
						}else if battle.GetBattleType() == face.PaiWeiBattleType{
							room := handler.GameServer.GetMatchManager().GetMatchBattleRoom(role.GetRoomId())
							if room != nil{
								room.Clear()
							}else {
								log4g.Infof("不存在匹配房间[%d]",role.GetRoomId())
							}
						}
						battle.Finish()
						log4g.Infof("战斗移除[%d]",battleId)
					}else if protoMsg.Cmd.CommandType == 47{
						//该玩家被占领了
						battle.RemovePlayer(innerMsg.RoleId)
						//不移除房间，移除role的房间信息
						role.QuitBattle()
					}
					battle.AddFrameCommand(protoMsg.Cmd.PlayerId,protoMsg.Cmd.CommandType,protoMsg.Cmd.Param)
				}else{
					log4g.Infof("不存在该战斗:[%d]",battleId)
				}
			}else{
				log4g.Infof("战斗Id无效:[%d],命令[%d]",battleId,protoMsg.Cmd.CommandType)
			}
		}
	}
}
