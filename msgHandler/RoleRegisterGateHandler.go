package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"cmgameserver/face"
)

type RoleRegisterGateHandler struct {
	GameServer face.IGameServer
}

func (handler *RoleRegisterGateHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg,ok:=msg.(network.InnerWriteMessage);ok {
		if protoMsg,ok := innerMsg.MsgData.(*message.G2M_RoleRegisterToGateSuccess);ok{
			rMsg := new(message.M2C_EnterLobby)
			if rMsg.RoleBasicInfo == nil{
				rMsg.RoleBasicInfo = new(message.RoleBasicInfo)
			}
			rMsg.IsInBattle = true
			rMsg.RoleBasicInfo.RoleId = protoMsg.RoleId
			rMsg.RoleBasicInfo.NickName = handler.GameServer.GetRoleManager().GetOnlineRole(protoMsg.RoleId).GetNickName()
			rMsg.RoleBasicInfo.Gold = 0
			handler.GameServer.WriteInnerMsg(session,protoMsg.RoleId,5000,rMsg)
		}
	}
}

