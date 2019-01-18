package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type GetAchieveHandler struct {
	GameServer face.IGameServer
}

func (handler *GetAchieveHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_GetAchievement); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//通过配置文件读取
				//然后加入角色已经领取的列表
				if role.AddGetAchieveAward(protoMsg.AchievenId){

				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetAchievement！")
		}
	}
}