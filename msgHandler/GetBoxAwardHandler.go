package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/tsrandom"
	"cmgameserver/bean"
	"time"
)

type GetBoxAwardHandler struct {
	GameServer face.IGameServer
}

func (handler *GetBoxAwardHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_GetBoxAwardItem); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				boxItem := role.GetItem(protoMsg.ItemIndex)
				if boxItem.ItemId == 0{
					return
				}
				items := handler.GetBoxEquipQuality(int(boxItem.ItemSeed),boxItem.ItemId,role)
				if items == nil{
					return
				}
				returnMsg := new(message.M2C_GetBoxAwardResult)
				if returnMsg.ItemList == nil{
					returnMsg.ItemList = make([]*message.Item,0)
				}
				for _,v := range items{
					if &v != nil{
						returnMsg.ItemList = append(returnMsg.ItemList, &v)
						//log4g.Infof("item:[%d]",v.ItemId)
					}
				}
				//设置宝箱消失
				boxItem.ItemId = 0
				boxItem.ItemSeed = 0
				boxItem.ItemNum = 0
				boxItem.ItemTime = 0
				role.SetItem(protoMsg.ItemIndex,boxItem)
				returnMsg.BoxItemIndex = protoMsg.ItemIndex
				//回送
				handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5045, returnMsg)
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M_GetBoxAwardItem！")
		}
	}
}
func (handler *GetBoxAwardHandler)GetBoxEquipQuality(seed int, boxItemId int32,onlineRole face.IOnlineRole) []message.Item{
	var boxChange [4]int
	var qulity bean.EItemQuality
	returnMsg := make([]message.Item,0)
	random := tsrandom.New(seed)
	boxList := handler.GameServer.GetBoxItemInfoConfig().BoxList
	box := boxList[boxItemId]
	if &box == nil{
		return nil
	}
	boxQuality := box.Quality
	has := false
	nowTime := time.Now().UnixNano()
	for _,v := range boxList{
		if v.Quality == boxQuality {
			boxChange = v.BoxChange
			has = true
			break
		}
	}
	if has == false{
		return nil
	}
	value := 0
	for i := 0; i < len(boxChange); i++{
		value += boxChange[i]                                   //获取总数
	}
	value += 1
	for i := 0; i < len(boxChange); i++ {
		r := random.RangeInt(0, value)
		if r < boxChange[i]{
			if i == 0{
				qulity = bean.White
				break
			}else if i == 1{
				qulity = bean.Green
				break
			}else if i == 2{
				qulity = bean.Blue
				break
			}else if i == 3{
				qulity = bean.Orange
				break
			}
		}else {
			value -= boxChange[i]
		}
	}
	soldier := random.RandomBool()
	if soldier == false{
		data := handler.GameServer.GetHeroQualityEquipConfig().Data
		all := data[qulity]
		index := random.RangeInt(0,len(all))
		equipId := int32(all[index].ItemId)
		item := &message.Item{
			ItemId:equipId,
			ItemSeed:int32(seed),
			ItemNum:1,
			ItemTime:nowTime,
		}
		itemIndex := onlineRole.AddItem(*item,false)
		item.Index = itemIndex
		returnMsg = append(returnMsg,*item)
	}else{
		data := handler.GameServer.GetSoldierQualityEquipConfig().Data
		all := data[qulity]
		index := random.RangeInt(0,len(all))
		equipId := int32(all[index].ItemId)
		item := &message.Item{
			ItemId:equipId,
			ItemSeed:int32(seed),
			ItemNum:1,
			ItemTime:nowTime,
		}
		itemIndex := onlineRole.AddItem(*item,false)
		item.Index = itemIndex
		returnMsg = append(returnMsg,*item)
	}
	return returnMsg
}
