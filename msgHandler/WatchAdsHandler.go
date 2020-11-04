package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"time"
)

type WatchAdsHandler struct {
	GameServer face.IGameServer
}

func (handler *WatchAdsHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_WatchAds); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				msg := new(message.M2C_WatchAdsResult)
				nowTime := time.Now().UnixNano()
				if protoMsg.IsBanner{
					role.AddGold(30)
					msg.Index = role.GetGold()
				}else{
					//视频观看，说明得随机数
					//seed := 0
					//if protoMsg.Seed > 0{
					//	seed = int(protoMsg.Seed)
					//}else{
					//	seed = int(time.Now().UnixNano())
					//}
					//random := tsrandom.New(seed)
					//gold := random.RangeInt(0,4)
					//if gold == 1{
					//	gold = 0
					//}
					//if gold == 0{
					//	msg.BoxId = 200006
					//	msg.Index = role.AddItemNoMsg(msg.BoxId,msg.Seed,nowTime,true)
					//}else if gold == 2{
					//	msg.BoxId = 200007
					//	msg.Index = role.AddItemNoMsg(msg.BoxId,msg.Seed,nowTime,true)
					//}else if gold == 3{
					//	msg.BoxId = 200008
					//	msg.Index = role.AddItemNoMsg(msg.BoxId,msg.Seed,nowTime,true)
					//}
					msg.BoxId = 200006
					msg.Index = role.AddItemNoMsg(msg.BoxId,msg.Seed,nowTime,true)
				}
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5021, msg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_WinBattle！")
		}
	}
}