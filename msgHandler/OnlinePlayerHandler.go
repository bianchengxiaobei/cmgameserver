package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type OnlinePlayerHandler struct {
	GameServer face.IGameServer
}

func (handler *OnlinePlayerHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_OnlinePlayer); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				msg := new(message.M2C_OnlinePlayerResult)
				if msg.Players == nil{
					msg.Players = make([]*message.OnlinePlayer,0)
				}
				allRole := handler.GameServer.GetRoleManager().GetAllOnlineRole(session)
				if len(allRole) > 0{
					for _,v := range allRole{
						if v.GetRoleId() == role.GetRoleId(){
							continue
						}
						online := new(message.OnlinePlayer)
						online.RoleId = v.GetRoleId()
						online.Name  = v.GetNickName()
						online.AvatarId = v.GetAvatarId()
						online.State = 0
						if v.GetInSimulateBattle(){
							online.State = 1
						}else if v.IsInRooming(){
							online.State = 2
						}else{
							matchPlayer := v.GetMatchPlayer()
							if matchPlayer.GetBInMatching() || matchPlayer.GetBInBattleRoom(){
								online.State = 3
							}
						}
						msg.Players = append(msg.Players, online)
					}
				}
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5107, msg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_LearnSkill！")
		}
	}
}
