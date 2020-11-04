package msgHandler

import (
	"cmgameserver/message"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type RoomChatHandler struct {
	GameServer face.IGameServer
}

func (handler *RoomChatHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M2C_Chat); ok {
			roleId := innerMsg.RoleId
			role := handler.GameServer.GetRoleManager().GetOnlineRole(roleId)
			if role != nil{
				if protoMsg.ChatType == message.C2M2C_Chat_LobbyWorld{
					if protoMsg.RoomId > 0{
						//说明是大厅发送组队信息
						//创建房间,自由房间还是显示，如果有玩家加入直接跳转
						if role.IsInRooming(){
							//重复创建房间,客户端约定好
							protoMsg.RoomId = -1000
							handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5032,protoMsg)
							return
						}
						room := handler.GameServer.GetRoomManager().CreateDefaultRoom(role.GetRoleId())
						if room != nil{
							protoMsg.RoomId = room.GetRoomId()
						}else{
							protoMsg.RoomId = -1001
							handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5032,protoMsg)
							return
						}
					}
					handler.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5032,protoMsg)
					roles := handler.GameServer.GetRoleManager().GetAllOnlineRole(session)
					for _,v := range roles{
						//这里应该也要过滤进入模拟联机或者战役
						if v.GetInSimulateBattle() || v.IsInRooming() || v.IsInBattling() || v.GetRoleId() == role.GetRoleId(){
							continue
						}
						if v != nil && v.IsConnected(){
							handler.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5032,protoMsg)
						}
					}
					//添加聊天进入缓存
					handler.GameServer.AddChatInfo(protoMsg.AvatarId,protoMsg.Name,protoMsg.Chat,
						protoMsg.Time,protoMsg.Level,protoMsg.RankScore,protoMsg.Sex,protoMsg.RoomId)
				}else if protoMsg.ChatType == message.C2M2C_Chat_LobbyFriend{

				}else if protoMsg.ChatType == message.C2M2C_Chat_RoomBattle{
					battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
					if battle != nil{
						roleIds := battle.GetBattleMember()
						for _,v := range roleIds{
							if v != nil && v.IsConnected(){
								handler.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5032,protoMsg)
							}
						}
					}else{
						room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
						if room != nil{
							roleIds := room.GetRoomRoleIds()
							for _,v := range roleIds{
								if v != 0{
									role := handler.GameServer.GetRoleManager().GetOnlineRole(v)
									if role != nil && role.IsConnected(){
										handler.GameServer.WriteInnerMsg(role.GetGateSession(),v,5032,protoMsg)
									}
								}
							}
						}else{
							log4g.Infof("没有房间[%d]",role.GetRoomId())
						}
					}
				}
			}
		}
	}
}