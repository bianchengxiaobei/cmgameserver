package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"time"
	"github.com/bianchengxiaobei/cmgo/tsrandom"
)

type WinBattleHandler struct {
	GameServer face.IGameServer
	getBoxQuality [3]int
}

func (handler *WinBattleHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_WinBattle); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				role.WinLevel(protoMsg.BattleId)
				msg := new(message.M2C_BattleResult)
				nowTime := time.Now().UnixNano()
				if protoMsg.Seed != 0{
					msg.Seed = protoMsg.Seed
				}else{
					msg.Seed =  int32(nowTime)
				}
				//钱袋
				value := 0
				random := tsrandom.New(int(msg.Seed))
				if handler.getBoxQuality[0] == 0 {
					handler.getBoxQuality = [3]int{5,20,75}
				}
				gold := random.RangeInt(0,4)
				if gold == 1{
					gold = 0
				}
				if gold == 0{
					msg.Award1 = 200006
					msg.Index1 = role.AddItemNoMsg(msg.Award1,msg.Seed,nowTime,false)
				}else if gold == 2{
					msg.Award1 = 200007
					msg.Index1 = role.AddItemNoMsg(msg.Award1,msg.Seed,nowTime,false)
				}else if gold == 3{
					msg.Award1 = 200008
					msg.Index1 = role.AddItemNoMsg(msg.Award1,msg.Seed,nowTime,false)
				}
				//经验
				exp := random.RangeInt(0,4)
				if exp == 1{
					exp = 0
				}
				if exp == 0{
					msg.Award2 = 200003
					msg.Index2 = role.AddItemNoMsg(msg.Award2,msg.Seed,nowTime,false)
				}else if exp == 2{
					msg.Award2 = 200004
					msg.Index2 = role.AddItemNoMsg(msg.Award2,msg.Seed,nowTime,false)
				}else if exp == 3{
					msg.Award2 = 200005
					msg.Index2 = role.AddItemNoMsg(msg.Award2,msg.Seed,nowTime,false)
				}
				//宝箱
				value = 0
				for _, v := range handler.getBoxQuality {
					value += v
				}
				value += 1
				for k, v := range handler.getBoxQuality {
					r := random.RangeInt(0, value)
					if r < v {
						if k == 0 {
							msg.Award3 = 200002
							msg.Index3 = role.AddItemNoMsg(msg.Award3,msg.Seed,nowTime,true)
							break
						}else if k== 1{
							msg.Award3 = 200001
							msg.Index3 = role.AddItemNoMsg(msg.Award3,msg.Seed,nowTime,true)
							break
						}else if k == 2{
							msg.Award3 = 200000
							msg.Index3 = role.AddItemNoMsg(msg.Award3,msg.Seed,nowTime,true)
							break
						}
					}else{
						value -= v
					}
				}
				log4g.Infof("1:[%d],2:[%d],3:[%d]",msg.Award1,msg.Award2,msg.Award3)
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5019, msg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_FailedBattle！")
		}
	}
}
