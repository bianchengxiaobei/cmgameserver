package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/bean"
)

type UseItemHandler struct {
	GameServer face.IGameServer
}

func (handler *UseItemHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_UseItem); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//通过配置文件读取
				item := role.GetItem(protoMsg.ItemIndex)
				itemId := item.ItemId
				itemNum := item.ItemNum
				if itemId != 0{
					config := handler.GameServer.GetMaterialConfig()
					materialData := config.Data[item.ItemId]
					if &materialData != nil{
						if materialData.MatType == bean.EGold{
							value := int32(materialData.MatValue * int(itemNum))
							role.AddGold(value)
						}else if materialData.MatType == bean.EExp{
							value := int32(materialData.MatValue * int(itemNum))
							handler.GameServer.GetRoleManager().AddRoleExp(role,value)
						}
					}else{
						//log4g.Infof("没有该物品：[%d]",item.ItemId)
						return
					}
				}else{
					//log4g.Infof("物品id为0：[%d]",protoMsg.ItemIndex)
					return
				}
				item.Clear()
				role.SetItem(protoMsg.ItemIndex,item)
				returnMsg := new(message.M2C_UseItemResult)
				returnMsg.ItemIndex = protoMsg.ItemIndex
				returnMsg.ItemNum = itemNum
				returnMsg.ItemId = itemId
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5053, returnMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_UseItem！")
		}
	}
}
