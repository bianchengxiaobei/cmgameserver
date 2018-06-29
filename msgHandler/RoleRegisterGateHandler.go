package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
)

type RoleRegisterGateHandler struct {
	GameServer IGameServer
}

func (handler *RoleRegisterGateHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg,ok:=msg.(network.InnerWriteMessage);ok {
		if protoMsg,ok := innerMsg.MsgData.(*message.G2M_RoleRegisterToGateSuccess);ok{
			rMsg := new(message.M2C_EnterLobby)
			rMsg.IsInBattle = true
			handler.GameServer.WriteInnerMsg(session,protoMsg.RoleId,5000,rMsg)
		}
	}
}

