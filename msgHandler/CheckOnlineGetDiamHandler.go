package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"time"
	"cmgameserver/tool"
)

type CheckOnlineGetDiamHandler struct {
	GameServer face.IGameServer
}

func (handler *CheckOnlineGetDiamHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_CheckOnlineGetDiam); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			rMsg := new(message.M2C_CheckOnlineGetDiam)
			if role != nil {
				now := time.Now()
				loginTime := role.GetLoginTime()
				min := int32(now.Sub(loginTime).Minutes())
				bei := int32(0)
				rankLevel := tool.GetRankLevelFromRankSocre(role.GetRankScore())
				if rankLevel == 0{
					bei = min / 3//测试10s
				}else if rankLevel == 1{
					bei = min / 2
				}else if rankLevel >= 2{
					bei = min
				}
				canAdd := false
				allLeft := int32(0)
				if role.GetOnlineDiam() != bei{
					canAdd,allLeft = handler.GameServer.SubAllOnlineDiam(bei - role.GetOnlineDiam())
				}
				if canAdd{
					role.SetOnlineDiam(bei)
					rMsg.Diam = bei
					rMsg.LeftDiam = allLeft
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5112, rMsg)
				}else{
					rMsg.Diam = role.GetOnlineDiam()
					rMsg.LeftDiam = 0
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5112, rMsg)
				}
			}
		} else {
			log4g.Error("不是C2M_CheckBuyShopDiam！")
		}
	}
}
