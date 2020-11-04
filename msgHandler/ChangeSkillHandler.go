package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ChangeSkillHandler struct {
	GameServer face.IGameServer
}
func (handler *ChangeSkillHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ChangeSkill); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				hero := role.GetHero(protoMsg.HeroId)
				if hero != nil{
					if protoMsg.Index == 1{
						hero.Skill1 = protoMsg.SkillId
					}else{
						hero.Skill2 = protoMsg.SkillId
					}
					role.SetHero(*hero)
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5035,protoMsg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_ChangeSkill！")
		}
	}
}