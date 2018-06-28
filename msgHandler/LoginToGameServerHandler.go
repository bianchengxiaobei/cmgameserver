package msgHandler

import (
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type LoginToGameServerHandler struct {
	GameServer IGameServer
}

func (handler *LoginToGameServerHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if protoMsg, ok := msg.(*message.G2M_LoginToGameServer); ok {
		roleId := protoMsg.RoleId
		log4g.Infof("id:%d",roleId)
	}
}
