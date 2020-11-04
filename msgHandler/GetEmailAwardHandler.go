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
				email := role.GetEmailByIndex(int(protoMsg.EmailId))
				if &email == nil{
					return
				}
				if email.Get{
					return
				}
				if role.SetEmail(int(protoMsg.EmailId),true) == false{
					log4g.Infof("领取邮件奖励失败:Index[%d]",protoMsg.EmailId)
					return
				}
				//通过配置文件读取
				//config := handler.GameServer.GetEmailConfig()
				//data := config.Data[email.EmailId]
				//if &data == nil{
					//return
				//}
				returnMsg := new(message.M2C_GetEmailResult)
				nowTime := time.Now().UnixNano()
				seed := int32(nowTime)
				var itemIndex int32
				if len(email.AwardList) > 0{
					for _,v := range email.AwardList{
						if v > 0{
							//看是否是赛季结算奖励，这里应该读取配置文件的
							//还要看是否金币，钻石等等
							itemIndex = role.AddItemNoMsg(v,seed,nowTime,true)
							item := new(message.Item)
							item.ItemId = v
							item.ItemTime = nowTime
							item.ItemSeed =seed
							item.Index = itemIndex
							returnMsg.AwardItem = append(returnMsg.AwardItem, item)
						}
					}
				}
				//if len(data.Awards) > 0{
				//	for _,v := range data.Awards{
				//		if v > 0{
				//			//如果是赛季结算奖励
				//			if email.EmailId == 16 || email.EmailId == 17 || email.EmailId == 18{
				//				itemIndex = 0
				//			}else{
				//				itemIndex = role.AddItemNoMsg(v,seed,nowTime,true)
				//			}
				//			item := new(message.Item)
				//			item.ItemId = v
				//			item.ItemTime = nowTime
				//			item.ItemSeed =seed
				//			item.Index = itemIndex
				//			returnMsg.AwardItem = append(returnMsg.AwardItem, item)
				//		}
				//	}
				//}
				//if len(data.OtherAwards) > 0{
				//	for k,v := range data.OtherAwards{
				//		awardType := bean.AwardType(k)
				//		switch awardType {
				//		case bean.Gold:
				//			role.AddGold(v)
				//			break
				//		case bean.Diam:
				//			role.AddDiam(v)
				//			break
				//		case bean.Exp:
				//			handler.GameServer.GetRoleManager().AddRoleExp(role,v)
				//			break
				//		}
				//	}
				//}
				returnMsg.EmailId = protoMsg.EmailId
				returnMsg.EmailIndex = protoMsg.EmailId
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