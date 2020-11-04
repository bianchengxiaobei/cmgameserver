package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"cmgameserver/face"
)

type RoleQuitMatchTeamHandler struct {
	GameServer face.IGameServer
}
//玩家退出匹配队伍
func (handler *RoleQuitMatchTeamHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_QuitMatchTeam); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				matchPlayer := role.GetMatchPlayer()
				if matchPlayer != nil{
					if matchPlayer.GetBInBattleRoom() == false{
						handler.GameServer.GetMatchManager().RemovePlayerFromMatchTeam(matchPlayer)
					}else{
						handler.GameServer.GetMatchManager().OnePlayerQuitMatchBattleRoom(matchPlayer)
					}
				}
			}
		}
	}
}
