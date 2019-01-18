package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"gopkg.in/mgo.v2/bson"
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
					bagId := role.GetItem(protoMsg.Index2).ItemId
					heroId := hero.ItemIds[int(protoMsg.Index1)]
					if bagId > 0 {
						hero.ItemIds[int(protoMsg.Index1)] = bagId
						role.GetItem(protoMsg.Index2).ItemId = heroId
					} else {
						hero.ItemIds[int(protoMsg.Index1)] = 0
						role.GetItem(protoMsg.Index2).ItemId = heroId
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
					item1Id := bag1.ItemId
					item2Id := bag2.ItemId
					bag1.ItemId = item2Id
					bag2.ItemId = item1Id
				case message.C2M2C_ChangeEquipItemPos_BagToHero:
					if protoMsg.Index2 > 2 {
						log4g.Info("数据错误!")
						return
					}
					bag := role.GetItem(protoMsg.Index1)
					itemId := bag.ItemId
					heroId := hero.ItemIds[int(protoMsg.Index2)]
					if heroId > 0 {
						//替换
						bag.ItemId = heroId
						hero.ItemIds[int(protoMsg.Index2)] = itemId
					} else {
						bag.ItemId = 0
						hero.ItemIds[int(protoMsg.Index2)] = itemId
					}
					//更新英雄数据库
					dbSession := handler.GameServer.GetDBManager().Get()
					if dbSession != nil {
						c := dbSession.DB("sanguozhizhan").C("Hero")
						data := bson.M{"$set": bson.M{"itemids": hero.ItemIds}}
						c.Update(bson.M{"roleid": innerMsg.RoleId, "heroid": protoMsg.HeroId}, data)
					}
				case message.C2M2C_ChangeEquipItemPos_SoldierToBag:
					bagId := role.GetItem(protoMsg.Index2).ItemId
					soldierIndex := int(protoMsg.HeroId)
					equipIndex := int(protoMsg.Index1)
					if bagId > 0 {
						//背包>0是替换
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagId)
						role.GetItem(protoMsg.Index2).ItemId = tempEquipId
					} else {
						tempEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
						role.GetItem(protoMsg.Index2).ItemId = tempEquipId
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,0)
					}
				case message.C2M2C_ChangeEquipItemPos_BagToSoldier:
					bag := role.GetItem(protoMsg.Index1)
					bagItemId := bag.ItemId
					soldierIndex := int(protoMsg.HeroId)
					equipIndex := int(protoMsg.Index2)
					soldierEquipId := role.GetFreeSoldierEquipId(soldierIndex,equipIndex)
					if soldierEquipId > 0{
						//替换
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItemId)
						bag.ItemId = soldierEquipId
					}else{
						bag.ItemId = 0
						role.ChangeFreeSoldierEquipId(soldierIndex,equipIndex,bagItemId)
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
