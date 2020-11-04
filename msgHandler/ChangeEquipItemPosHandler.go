package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/bean"
)

type ChangeEquipItemPosHandler struct {
	GameServer face.IGameServer
}

func (handler *ChangeEquipItemPosHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_ChangeEquipItemPos); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			hero := role.GetHero(protoMsg.HeroId)
			if role != nil {
				switch protoMsg.ChangeType {
				case message.C2M2C_ChangeEquipItemPos_HeroToBag:
					//看背包是否id>0，如果>0说明是替换，否则是英雄变成官方，增加背包
					if hero == nil{
						log4g.Infof("不存在武将[%d]!",protoMsg.HeroId)
						return
					}
					if protoMsg.Index1 > 2 {
						log4g.Info("数据错误!")
						return
					}
					bagItem := role.GetItem(protoMsg.Index2)
					heroItem := &hero.ItemIds[int(protoMsg.Index1)]
					if bagItem.ItemId > 0 {
						hero.ItemIds[int(protoMsg.Index1)] = bagItem
						role.SetItem(protoMsg.Index2,*heroItem)
					} else {
						role.SetItem(protoMsg.Index2,*heroItem)
						//官方
						data := handler.GameServer.GetHeroConfig().Data[protoMsg.HeroId]
						var guanFangId int32
						if protoMsg.Index1 == 0{
							guanFangId = data.GuanFangBodyId
						}else if protoMsg.Index1 == 1{
							guanFangId = data.GuanFangWeapId
						}else if protoMsg.Index1 == 2{
							guanFangId = data.GuanFangShoeId
						}
						heroItem.ItemId = guanFangId
						heroItem.ItemNum = 1
						heroItem.ItemSeed = 0
					}
					role.SetHero(*hero)
				case message.C2M2C_ChangeEquipItemPos_BagToBag:
					//直接双方替换
					bag1 := role.GetItem(protoMsg.Index1)
					bag2 := role.GetItem(protoMsg.Index2)
					role.SetItem(protoMsg.Index1,bag2)
					role.SetItem(protoMsg.Index2,bag1)
				case message.C2M2C_ChangeEquipItemPos_BagToHero:
					if protoMsg.Index2 > 2 {
						log4g.Info("数据错误!")
						return
					}
					heroIndex := int(protoMsg.Index2)
					bagItem := role.GetItem(protoMsg.Index1)
					heroItem := &hero.ItemIds[heroIndex]
					//如果heroItem == null说明要设置为官方
					if heroItem == nil{
						heroData := handler.GameServer.GetHeroConfig().Data[protoMsg.HeroId]
						var tempItemId int32
						if heroIndex == 0{
							tempItemId = heroData.GuanFangBodyId
						}else if heroIndex == 1{
							tempItemId = heroData.GuanFangWeapId
						}else if heroIndex == 2{
							tempItemId = heroData.GuanFangShoeId
						}
						tempHeroItem := bean.Item{
							ItemId:tempItemId,
							ItemTime:0,
							ItemSeed:0,
							ItemNum:1,
						}
						hero.ItemIds[heroIndex] = tempHeroItem
					}
					bGaunFang := handler.GameServer.GetHeroItemIdEquipConfig().Data[heroItem.ItemId].GuanFang
					if heroItem.ItemId == 0{
						bGaunFang = true
					}
					if bGaunFang == false{
						//如果现在武将身上的不是官方的，就替换背包的
						role.SetItem(protoMsg.Index1,*heroItem)
						heroItem.ItemId = bagItem.ItemId
						heroItem.ItemSeed = bagItem.ItemSeed
						heroItem.ItemNum = bagItem.ItemNum
					} else {
						//如果现在武将身上的是官方的，就直接设置背包为空
						nullItem := new(bean.Item)
						nullItem.ItemNum = 0
						nullItem.ItemId = 0
						nullItem.ItemSeed = 0
						role.SetItem(protoMsg.Index1,*nullItem)
						heroItem.ItemId = bagItem.ItemId
						heroItem.ItemSeed = bagItem.ItemSeed
						heroItem.ItemNum = bagItem.ItemNum
					}
					role.SetHero(*hero)
				case message.C2M2C_ChangeEquipItemPos_SoldierToBag:
					bagItem := role.GetItem(protoMsg.Index2)
					soldierIndex := int(protoMsg.HeroId)
					equipIndex := int(protoMsg.Index1)
					carryType := protoMsg.OtherValue
					if bagItem.ItemId > 0 {
						//背包>0是替换
						//判断背包是否是士兵装备
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
						oS := handler.GameServer.GetSoldierItemIdEquipConfig().Data[tempEquipId]
						if soldierConfig,ok := handler.GameServer.GetSoldierItemIdEquipConfig().Data[bagItem.ItemId];ok == false{
							log4g.Infof("不是士兵装备[%d]",bagItem.ItemId)
							return
						}else{
							if soldierConfig.ItemType == bean.Weap{
								if soldierConfig.PlayerType != oS.PlayerType{
									return
								}
							}
						}
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId,carryType)
						soldierItem := new(bean.Item)
						soldierItem.ItemId = tempEquipId
						soldierItem.ItemNum = 1
						soldierItem.ItemSeed = 0
						role.SetItem(protoMsg.Index2,*soldierItem)
					} else {
						//让士兵换上官方id
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)//现在士兵的装备
						guangfanId,_ := role.GetFreeSoldierGuangFanEquipId(soldierIndex,equipIndex)//需要替换的士兵官方装备
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,guangfanId,carryType)
						soldierItem := new(bean.Item)
						soldierItem.ItemId = tempEquipId
						soldierItem.ItemNum = 1
						soldierItem.ItemSeed = 0
						//然后背包设置上原先士兵的装备id
						role.SetItem(protoMsg.Index2,*soldierItem)
					}
				case message.C2M2C_ChangeEquipItemPos_BagToSoldier:
					bagItem := role.GetItem(protoMsg.Index1)
					soldierIndex := int(protoMsg.HeroId)
					equipIndex := int(protoMsg.Index2)
					carryType := protoMsg.OtherValue
					soldierEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
					oS := handler.GameServer.GetSoldierItemIdEquipConfig().Data[soldierEquipId]
					if soldierConfig,ok := handler.GameServer.GetSoldierItemIdEquipConfig().Data[bagItem.ItemId];ok == false{
						log4g.Infof("不是士兵装备[%d]",bagItem.ItemId)
						return
					}else{
						if soldierConfig.ItemType == bean.Weap{
							if soldierConfig.PlayerType != oS.PlayerType{
								return
							}
						}
					}
					_,isGuangFan := role.GetFreeSoldierGuangFanEquipId(soldierIndex,equipIndex)
					//判断替换的士兵上是否是官方额，如果是官方的就不放在背包上
					if isGuangFan == false{
						//替换
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId,carryType)
						soldierItem := new(bean.Item)
						soldierItem.ItemId = soldierEquipId
						soldierItem.ItemSeed = 0
						soldierItem.ItemNum = 1
						role.SetItem(protoMsg.Index1,*soldierItem)
					}else{
						soldierItem := new(bean.Item)
						soldierItem.ItemId = 0
						soldierItem.ItemSeed = 0
						soldierItem.ItemNum = 0
						role.SetItem(protoMsg.Index1,*soldierItem)
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId,carryType)
					}
				}
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5026, protoMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_ChangeEquipItemPos！")
		}
	}
}