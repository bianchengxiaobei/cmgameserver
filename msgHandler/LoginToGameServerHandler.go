package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
)

type LoginToGameServerHandler struct {
	GameServer IGameServer
}
func (handler *LoginToGameServerHandler)Action(session network.SocketSessionInterface,msg interface{}){
	if protoMsg,ok := msg.(*message.G2M_LoginToGameServer);ok{
		roleId := protoMsg.RoleId
	}
}
