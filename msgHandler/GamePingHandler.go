package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"time"
)

type GamePingHandler struct {
	GameServer face.IGameServer
}
func (handler *GamePingHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.M2C2M_GamePing); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.SetPingTime(time.Now())
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是M2C2M_GamePing！")
		}
	}
}
