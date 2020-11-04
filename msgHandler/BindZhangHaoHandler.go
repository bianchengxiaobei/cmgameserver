package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"strconv"
)

type BindZhangHaoHandler struct {
	GameServer face.IGameServer
}

func (handler *BindZhangHaoHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_BangDingZhang); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil{
				if protoMsg.ZhangHaoType == message.ZhangHaoType_EEmail{
					role.SetMailAddresss(protoMsg.ZhangHao)
				}else if protoMsg.ZhangHaoType == message.ZhangHaoType_Phone{
					num,_ := strconv.Atoi(protoMsg.ZhangHao)
					role.SetPhone(int32(num))
				}else if protoMsg.ZhangHaoType == message.ZhangHaoType_QQ{
					num,_ := strconv.Atoi(protoMsg.ZhangHao)
					role.SetQQ(int32(num))
				}else if protoMsg.ZhangHaoType == message.ZhangHaoType_WeiXin{
					role.SetWeiXin(protoMsg.ZhangHao)
				}
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5061, protoMsg)
				//发送给网关更新User
				rMsg := new(message.M2G_BindZhangHao)
				rMsg.ZhangHaoType = protoMsg.ZhangHaoType
				rMsg.UserId = role.GetUserId()
				rMsg.Value = protoMsg.ZhangHao
				handler.GameServer.WriteInnerMsg(session,0,10007,rMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_AgreePauseBattle！")
		}
	}
}