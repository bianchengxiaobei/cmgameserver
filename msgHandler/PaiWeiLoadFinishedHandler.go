package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
)

type PaiWeiLoadFinishedHandler struct {
	GameServer face.IGameServer
}

func (handler *PaiWeiLoadFinishedHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_PaiWeiLoadFinished); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetLoadFinished(true)
				matchPlayer := role.GetMatchPlayer()
				if matchPlayer != nil {
					allLoad, allPlayer := handler.GameServer.GetMatchManager().
						CheckAllPlayerLoadFinished(matchPlayer.GetMatchRoomId())
					if allLoad && allPlayer != nil {
						//如果所有人加载好了
						battle := handler.GameServer.GetBattleManager().GetBattleInFree()
						if battle != nil {
							battle.ReStart(allPlayer,face.PaiWeiBattleType)
							log4g.Infof("重新利用战斗[%d]",battle.GetBattleId())
						}else {
							battle := handler.GameServer.GetBattleManager().CreateBattle(allPlayer,face.PaiWeiBattleType)
							if battle != nil {
								battle.Start()
							}
						}
					}
				}
			}
		} else {
			log4g.Error("不是C2M_PaiWeiLoadFinished！")
		}
	}
}
