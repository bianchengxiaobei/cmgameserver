package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"time"
)

type GetEmailAwardHandler struct {
	GameServer face.IGameServer
}


func (handler *GetEmailAwardHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_GetEmailAward); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				email,_ := role.GetEmail(protoMsg.EmailId)
				if email == nil{
					return
				}
				if email.Get{
					return
				}
				email.Get = true
				//通过配置文件读取
				config := handler.GameServer.GetEmailConfig()
				data := config.Data[email.EmailId]
				if &data == nil{
					return
				}
				returnMsg := new(message.M2C_GetEmailResult)
				nowTime := time.Now().UnixNano()
				seed := int32(nowTime)
				for _,v := range data.Awards{
					if v > 0{
						itemIndex := role.AddItemNoMsg(v,seed,nowTime,true)
						item := new(message.Item)
						item.ItemId = v
						item.ItemTime = nowTime
						item.ItemSeed =seed
						item.Index = itemIndex
						returnMsg.AwardItem = append(returnMsg.AwardItem, item)
					}
				}
				returnMsg.EmailId = protoMsg.EmailId
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5049, returnMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetAchievement！")
		}
	}
}