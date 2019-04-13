package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type SellItemHandler struct {
	GameServer face.IGameServer
}

func (handler *SellItemHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_SellItem); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//通过配置文件读取
				item := role.GetItem(protoMsg.ItemIndex)
				if item.ItemId != 0{
					config := handler.GameServer.GetSoldierItemIdEquipConfig()
					soldierData := config.Data[item.ItemId]
					if &soldierData != nil{
						role.AddGold(soldierData.SellGold)
					}else{
						heroData := handler.GameServer.GetHeroItemIdEquipConfig().Data[item.ItemId]
						if &heroData != nil{
							role.AddGold(heroData.SellGold)
						}
					}
				}else{
					return
				}
				item.Clear()
				role.SetItem(protoMsg.ItemIndex,item)
				gold := role.GetGold()
				returnMsg := new(message.M2C_SellItemResult)
				returnMsg.ItemIndex = protoMsg.ItemIndex
				returnMsg.RoleGold = gold
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5047, returnMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetTask！")
		}
	}
}