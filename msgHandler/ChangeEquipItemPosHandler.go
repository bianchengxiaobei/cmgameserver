package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"gopkg.in/mgo.v2/bson"
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
					//看背包是否id>0，如果>0说明是替换，否则是移除英雄，增加背包
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
						heroItem.ItemId = 0
						heroItem.ItemNum = 0
						heroItem.ItemSeed = 0
					}
					//更新英雄数据库
					dbSession := handler.GameServer.GetDBManager().Get()
					if dbSession != nil {
						c := dbSession.DB("sanguozhizhan").C("Hero")
						data := bson.M{"$set": bson.M{"itemids": hero.ItemIds}}
						c.Update(bson.M{"roleid": innerMsg.RoleId, "heroid": protoMsg.HeroId}, data)
					}
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
					bagItem := role.GetItem(protoMsg.Index1)
					heroItem := &hero.ItemIds[int(protoMsg.Index2)]
					if heroItem.ItemId > 0 {
						//替换
						role.SetItem(protoMsg.Index1,*heroItem)
						heroItem.ItemId = bagItem.ItemId
						heroItem.ItemSeed = bagItem.ItemSeed
						heroItem.ItemNum = bagItem.ItemNum
					} else {
						nullItem := new(bean.Item)
						nullItem.ItemNum = 0
						nullItem.ItemId = 0
						nullItem.ItemSeed = 0
						role.SetItem(protoMsg.Index1,*nullItem)
						heroItem.ItemId = bagItem.ItemId
						heroItem.ItemSeed = bagItem.ItemSeed
						heroItem.ItemNum = bagItem.ItemNum
					}
					//更新英雄数据库
					dbSession := handler.GameServer.GetDBManager().Get()
					if dbSession != nil {
						c := dbSession.DB("sanguozhizhan").C("Hero")
						data := bson.M{"$set": bson.M{"itemids": hero.ItemIds}}
						c.Update(bson.M{"roleid": innerMsg.RoleId, "heroid": protoMsg.HeroId}, data)
					}
				case message.C2M2C_ChangeEquipItemPos_SoldierToBag:
					bagItem := role.GetItem(protoMsg.Index2)
					soldierIndex := int(protoMsg.HeroId)
					equipIndex := int(protoMsg.Index1)
					if bagItem.ItemId > 0 {
						//背包>0是替换
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId)
						soldierItem := new(bean.Item)
						soldierItem.ItemId = tempEquipId
						soldierItem.ItemNum = 1
						soldierItem.ItemSeed = 0
						role.SetItem(protoMsg.Index2,*soldierItem)
					} else {
						//让士兵换上官方id
						guangfanId,_ := role.GetFreeSoldierGuangFanEquipId(soldierIndex,equipIndex)
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,guangfanId)
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

					soldierEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
					_,isGuangFan := role.GetFreeSoldierGuangFanEquipId(soldierIndex,equipIndex)
					//判断替换的士兵上是否是官方额，如果是官方的就不放在背包上
					if isGuangFan == false{
						//替换
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId)
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
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItem.ItemId)
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