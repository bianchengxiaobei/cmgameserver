package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/face"
	"time"
	"cmgameserver/bean"
	"math/rand"
)

type GetGiftHandler struct {
	GameServer face.IGameServer
}
//玩家领取礼品码处理
func (handler *GetGiftHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_GetGift); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//加载配置文件，
				config := handler.GameServer.GetGiftConfig()
				heroData := handler.GameServer.GetHeroConfig().Data
				item := bean.Item{
					ItemNum:1,
				}
				// 看角色是否已经领取过了
				msg := new(message.M2C_GetGiftResult)
				if &config != nil{
					if data,ok := config.Datas[protoMsg.Code];ok{
						//礼品码正确，看是否已经领取过了
						if role.GetGiftCode(protoMsg.Code){
							//发送给客户端成功
							msg.Success = 1
							msg.Gold = data.Gold
							role.AddGold(data.Gold)
							msg.Diam = data.Diam
							role.AddDiam(data.Diam)
							msg.Exp = data.Exp
							handler.GameServer.GetRoleManager().AddRoleExp(role,data.Exp)
							msg.HeroId = data.HeroId
							getHero := bean.Hero{
								HeroId:msg.HeroId,
								Level:1,
								RoleId:role.GetRoleId(),
							}
							//看玩家是已经有该武将，如果已经存在就随机获得一个未获得的武将
							hero := role.GetHero(msg.HeroId)
							if hero == nil{
								tempHeroData := heroData[msg.HeroId]
								//为拥有
								item.ItemId = tempHeroData.GuanFangBodyId
								getHero.ItemIds[0] = item
								item.ItemId = tempHeroData.GuanFangWeapId
								getHero.ItemIds[1] = item
								item.ItemId = tempHeroData.GuanFangShoeId
								getHero.ItemIds[2] = item
								//插入服务器
								handler.GameServer.GetRoleManager().AddRoleHero(role,getHero)
							}else{
								//已经有了
								list := make([]int32,1)//未获得的武将
								for _,v := range heroData{
									has := role.GetHero(v.HeroId)
									if has == nil{
										list = append(list, v.HeroId)
									}
								}
								Llen := len(list)
								if Llen > 0{
									//随机
									rHeroIndex := rand.Intn(Llen)
									rHeroId := list[rHeroIndex]
									//添加
									if rHeroId != 0{
										getHero.HeroId = rHeroId
										tempHeroData := heroData[rHeroId]
										//tempHeroData.HeroId = rHeroId
										item.ItemId = tempHeroData.GuanFangBodyId
										getHero.ItemIds[0] = item
										item.ItemId = tempHeroData.GuanFangWeapId
										getHero.ItemIds[1] = item
										item.ItemId = tempHeroData.GuanFangShoeId
										getHero.ItemIds[2] = item
										msg.HeroId = rHeroId
										handler.GameServer.GetRoleManager().AddRoleHero(role,getHero)
									}else{
										msg.HeroId = 0
										msg.Diam += 400
										role.AddDiam(400)
									}
								}
							}
							now := time.Now()
							for _,v := range data.Awards{
								item := message.Item{
									ItemId:v,
									ItemSeed:int32(now.Unix()),
									ItemTime:now.Unix(),
								}
								index := role.AddItem(item,false)
								item.Index = index
								msg.AwardItem = append(msg.AwardItem, &item)
							}
						}else{
							//发送给客户端失败
							msg.Success = 0
						}
					}
				}else{
					msg.Success = 0
				}
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5098, msg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetSign！")
		}
	}
}