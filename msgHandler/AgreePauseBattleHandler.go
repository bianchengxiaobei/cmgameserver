package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type AgreePauseBattleHandler struct {
	GameServer face.IGameServer
}

func (handler *AgreePauseBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_AgreePauseBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
			if role != nil  && battle != nil{
				pause := battle.SetAgreePause(protoMsg.Agree,innerMsg.RoleId)
				if pause == 1{
					//发送暂停消息
					returnMsg := new(message.M2C_StartPause)
					returnMsg.Pause = true
					returnMsg.RoleId = battle.GetPauseRoleId()
					for _,v := range battle.GetBattleMember(){
						if v != nil && v.GetRoleId() > 0{
							handler.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5056,returnMsg)
						}
					}
				}else if pause == 0{
					//不能暂停
					returnMsg := new(message.M2C_StartPause)
					returnMsg.Pause = false
					returnMsg.RoleId = battle.GetPauseRoleId()
					for _,v := range battle.GetBattleMember(){
						if v != nil && v.GetRoleId() > 0{
							handler.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5056,returnMsg)
						}
					}
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_AgreePauseBattle！")
		}
	}
}