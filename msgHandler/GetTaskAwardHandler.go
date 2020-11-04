package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"time"
	"cmgameserver/bean"
)

type GetTaskAwardHandler struct {
	GameServer face.IGameServer
}

func (handler *GetTaskAwardHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_GetTask); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				//通过配置文件读取
				//然后加入角色已经领取的列表
				if role.AddGetTaskAward(protoMsg.TaskId){
					config := handler.GameServer.GetTaskConfig()
					data := config.Tasks[protoMsg.TaskId]
					for k,v := range data.Award{
						if k == bean.Gold{
							role.AddGold(int32(v))
						}else if k == bean.Exp{
							handler.GameServer.GetRoleManager().AddRoleExp(role,int32(v))
						}else if k == bean.Diam{
							role.AddDiam(int32(v))
						}else if k == bean.ItemType{
							//宝箱,增加到背包
							time := time.Now().Unix()
							role.AddItemNoMsg(int32(v),int32(time), time,false)
						}
					}
					//回送
					handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5037, protoMsg)
				}
			} else {
				log4g.Errorf("不存在RoleId:%d", innerMsg.RoleId)
			}
		} else {
			log4g.Error("不是C2M2C_GetTask！")
		}
	}
}