package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
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
				rMsg := &message.M2C_LearnSkillResult{}
				if hero != nil {
					config := handler.GameServer.GetSkillConfig()
					needSkillPoint := config.Datas[protoMsg.SkillId].LearnPoint
					if needSkillPoint == 0{
						println("11111")
					}
					if hero.HeroLearnSkill(protoMsg.SkillId, needSkillPoint) == false {
						rMsg.HeroId = -1
						log4g.Infof("[%d]学习技能[%d]失败!", innerMsg.RoleId, protoMsg.SkillId)
					} else {
						rMsg.HeroId = protoMsg.HeroId
						rMsg.SkillId = protoMsg.SkillId
						rMsg.SkillPoint = hero.SkillPoint
					}
				} else {
					rMsg.HeroId = -1
				}
				role.SetHero(*hero)
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5034, rMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_LearnSkill！")
		}
	}
}
