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
				items,notifyId := handler.GetBoxEquipQuality(int(boxItem.ItemSeed),boxItem.ItemId,role)
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
				//向全服发送取得该物品
				if notifyId > 0{
					//说明得到的是传奇以上的装备,通知全服玩家
					handler.GameServer.GetRoleManager().SendRollInfoToAllRoleGetItem(notifyId,role.GetNickName())
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
func (handler *GetBoxAwardHandler)GetBoxEquipQuality(seed int, boxItemId int32,onlineRole face.IOnlineRole) ([]message.Item,int32){
	var boxChange [4]int
	var qulity bean.EItemQuality
	var notifyId int32
	notifyId = -1
	returnMsg := make([]message.Item,0)
	random := tsrandom.New(seed)
	boxList := handler.GameServer.GetBoxItemInfoConfig().BoxList
	box := boxList[boxItemId]
	if &box == nil{
		return nil,notifyId
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
		return nil,notifyId
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
		Get:
		index := random.RangeInt(0,len(all))
		//如果是官方重新获取
		if all[index].GuanFang{
			goto Get
		}
		equipId := all[index].ItemId
		if qulity == bean.Orange{
			notifyId = equipId
		}
		item := &message.Item{
			ItemId:equipId,
			ItemSeed:int32(seed),
			ItemNum:1,
			ItemTime:nowTime,
		}
		itemIndex := onlineRole.AddItem(*item,false)
		item.Index = itemIndex
		returnMsg = append(returnMsg,*item)
		//如果是紫色宝箱，必得紫色装备
		if boxItemId == 200010{
			ziseAll := data[bean.ZiSe]
			ziseIndex := random.RangeInt(0,len(ziseAll))
			ziseId := ziseAll[ziseIndex].ItemId
			notifyId = ziseId
			ziseItem := &message.Item{
				ItemId:ziseId,
				ItemSeed:int32(seed),
				ItemNum:1,
				ItemTime:nowTime,
			}
			ziseItem.Index = onlineRole.AddItem(*ziseItem,false)
			returnMsg = append(returnMsg,*ziseItem)
		}
	}else{
		data := handler.GameServer.GetSoldierQualityEquipConfig().Data
		all := data[qulity]
		HeroStart:
		index := random.RangeInt(0,len(all))
		equipId := all[index].ItemId
		//如果是官方重新获取
		if handler.IsGuangFangId(equipId){
			goto HeroStart
		}
		//如果是传奇装备的话，就通知
		if qulity == bean.Orange{
			notifyId = equipId
		}
		item := &message.Item{
			ItemId:equipId,
			ItemSeed:int32(seed),
			ItemNum:1,
			ItemTime:nowTime,
		}
		itemIndex := onlineRole.AddItem(*item,false)
		item.Index = itemIndex
		returnMsg = append(returnMsg,*item)
		//如果是紫色宝箱，必得紫色装备
		if boxItemId == 200010{
			ziseAll := data[bean.ZiSe]
			ziseIndex := random.RangeInt(0,len(ziseAll))
			ziseId := ziseAll[ziseIndex].ItemId
			notifyId = ziseId
			ziseItem := &message.Item{
				ItemId:ziseId,
				ItemSeed:int32(seed),
				ItemNum:1,
				ItemTime:nowTime,
			}
			ziseItem.Index = onlineRole.AddItem(*ziseItem,false)
			returnMsg = append(returnMsg,*ziseItem)
		}
	}
	return returnMsg,notifyId
}
func (handler *GetBoxAwardHandler)IsGuangFangId(id int32)bool{
	if id == 1 || id == 2001 || id == 4001 ||
		id == 5001 || id == 7001 || id == 8001{
		//如果是官方的装备
		return true
	}else{
		return false
	}
}
