package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"strconv"
	"strings"
)

type ECommandType int32

const(
	SendAllUserEmail ECommandType = 0//向全服的玩家发送邮件奖励
	BuyCard  ECommandType = 1//某个玩家购买月卡等卡片
	SendOneUseEmail ECommandType = 2//向某个玩家发送邮件奖励
	BuyDiam ECommandType = 3//某个玩家充值钻石
	SendAllUserRollInfo ECommandType = 4//向全服玩家发送滚动信息
	CloseServer ECommandType = 5//延时关闭服务器
	GetOnlineNum ECommandType = 6//取得在线的人数
)

type GMCommandHandler struct {
	GameServer face.IGameServer
}

func (handler *GMCommandHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_GMCommand); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role != nil {
				if role.GetIsGM() == false{
					//如果不是管理员就返回，防止黑客恶意刷装备
					return
				}
				returnMsg := new (message.M2C_GMCommandResult)
				//判断指令id
				switch ECommandType(protoMsg.CommandId) {
				case SendAllUserEmail:
					//emailId,err := strconv.Atoi(protoMsg.CommandParam)
					var err error
					var id int
					content := strings.Split(protoMsg.CommandParam,";")
					title := content[0]
					emailContent := content[1]
					list := strings.Split(content[2],",")
					awardList := make([]int32,len(list) - 1)
					for _,v := range list{
						if  v != ""{
							id,err = strconv.Atoi(v)
							awardList = append(awardList, int32(id))
						}
					}
					if err != nil{
						returnMsg.Result = false
						log4g.Info(err.Error())
					}else{
						returnMsg.Result = handler.GameServer.
							GetRoleManager().SendEmailToAllRole(title,emailContent,awardList)
					}
					break
				case SendAllUserRollInfo:
					returnMsg.Result = handler.GameServer.GetRoleManager().SendRollInfoToAllRole(protoMsg.CommandParam)
					break
				case SendOneUseEmail:
					content := strings.Split(protoMsg.CommandParam,";")
					roleId,err := strconv.ParseInt(content[0],10,64)
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					title,emailContent,awardList,err := GetEmailCommand(content[1])
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					returnMsg.Result = handler.GameServer.GetRoleManager().
						SendEmailToOneRole(roleId,title,emailContent,awardList)
					break
				case BuyCard:
					content := strings.Split(protoMsg.CommandParam,";")
					roleId,err := strconv.ParseInt(content[0],10,64)
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					cardType,err := strconv.Atoi(content[1])
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					returnMsg.Result = handler.GameServer.GetRoleManager().BuyCard(roleId,int32(cardType))
					break
				case BuyDiam:
					content := strings.Split(protoMsg.CommandParam,";")
					roleId,err := strconv.ParseInt(content[0],10,64)
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					diamType,err := strconv.Atoi(content[1])
					if err != nil{
						log4g.Info(err.Error())
						returnMsg.Result = false
						break
					}
					returnMsg.Result = handler.GameServer.GetRoleManager().BuyDiam(roleId,int32(diamType))
					break
				case CloseServer:
					handler.GameServer.SetNeedClose()
					//发送滚动信息
					handler.GameServer.GetRoleManager().SendRollInfoToAllRole("尊敬的玩家,服务器即将关闭,以免造成游戏数据丢失,请及时退出游戏,谢谢配合!")
					break
				case GetOnlineNum:
					returnMsg.ResultValue = int32(handler.GameServer.GetRoleManager().GetOnlineNum())
					break
				}
				returnMsg.CommandId = protoMsg.CommandId
				handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),100002,returnMsg)
			}
		} else {
			log4g.Error("不是C2M_GMCommand！")
		}
	}
}
func GetEmailCommand(commandParam string) (string,string,[]int32,error) {
	var id int
	var err error
	content := strings.Split(commandParam,";")
	title := content[0]
	emailContent := content[1]
	list := strings.Split(content[2],",")
	awardList := make([]int32,len(list) - 1)
	for _,v := range list{
		if  v != ""{
			id,err = strconv.Atoi(v)
			awardList = append(awardList, int32(id))
		}
	}
	return title,emailContent,awardList,err
}
