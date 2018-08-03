package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
)

type RoleQuitHandler struct {
	GameServer 		face.IGameServer
}

func (handler *RoleQuitHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.G2M_RoleQuitGameServer); ok {
			handler.GameServer.GetRoleManager().RoleQuit(protoMsg.RoleId)
		}
	}
}

