package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"time"
)

type QuickGetAllNoReadEmailHandler struct {
	GameServer face.IGameServer
}


func (handler *QuickGetAllNoReadEmailHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_GetAllNoReadEmail); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				returnMsg := new(message.M2C_GetAllNoReadEmailResult)
				//config := handler.GameServer.GetEmailConfig()
				nowTime := time.Now().UnixNano()
				seed := int32(nowTime)
				for i:=0;i<role.GetMaxEmailNum();i++{
					email := role.GetEmailByIndex(i)
					if &email != nil{
						if email.Get || email.Valid == false{
							//已经领取过额不要
							continue
						}
						if role.SetEmail(i,true) == false{
							log4g.Infof("领取邮件奖励失败:[%d]",i)
							continue
						}
						//index
						returnMsg.EmailIndex = append(returnMsg.EmailIndex, int32(i))
						//奖励
						if len(email.AwardList) > 0{
							for _,v := range email.AwardList{
								if v > 0{
									//宝箱之类的
									itemIndex := role.AddItemNoMsg(v,seed,nowTime,true)
									item := new(message.Item)
									item.ItemId = v
									item.ItemTime = nowTime
									item.ItemSeed =seed
									item.Index = itemIndex
									returnMsg.AwardItems = append(returnMsg.AwardItems, item)
								}
							}
						}
						//通过配置文件读取
						//data := config.Data[email.EmailId]
						//if &data == nil{
						//	continue
						//}
						//if len(data.Awards) > 0{
						//	for _,v := range data.Awards{
						//		if v > 0{
						//			//宝箱之类的
						//			itemIndex := role.AddItemNoMsg(v,seed,nowTime,true)
						//			item := new(message.Item)
						//			item.ItemId = v
						//			item.ItemTime = nowTime
						//			item.ItemSeed =seed
						//			item.Index = itemIndex
						//			returnMsg.AwardItems = append(returnMsg.AwardItems, item)
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
					}
				}
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5093, returnMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_GetAllNoReadEmail！")
		}
	}
}
