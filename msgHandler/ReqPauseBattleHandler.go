package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ReqPauseBattleHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqPauseBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ReqPauseBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
				if battle != nil{
					state := battle.GetBattleState()
					//如果不在战斗中，不能暂停
					if state != face.InBattling{
						return
					}
					battle.SetPauseRoleId(innerMsg.RoleId)
					battle.SetAgreePause(true,innerMsg.RoleId)
					//回送给每个战斗成员
					players := battle.GetBattleMember()
					for _,v := range players{
						if v != nil && v.GetRoleId() > 0{
							if v.GetRoleId() != role.GetRoleId(){
								handler.GameServer.WriteInnerMsg(v.GetGateSession(), v.GetRoleId(), 5054, protoMsg)
							}
						}
					}
				}else{
					log4g.Info("No vattle")
				}
			}
		}
	}
}

