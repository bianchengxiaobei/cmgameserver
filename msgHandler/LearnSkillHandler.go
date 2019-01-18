package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type LearnSkillHandler struct {
	GameServer face.IGameServer
}

func (handler *LearnSkillHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_LearnSkill); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				hero := role.GetHero(protoMsg.HeroId)
				if hero != nil{
					if hero.HeroLearnSkill(protoMsg.SkillId) == false{
						log4g.Infof("[%d]学习技能[%d]失败!",innerMsg.RoleId,protoMsg.SkillId)
					}else{
						//发送消息
						rMsg := &message.M2C_LearnSkillResult{}
						rMsg.HeroId = protoMsg.HeroId
						rMsg.SkillId = protoMsg.SkillId
						rMsg.SkillPoint =hero.SkillPoint
						handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5034,rMsg)
					}
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_LearnSkill！")
		}
	}
}