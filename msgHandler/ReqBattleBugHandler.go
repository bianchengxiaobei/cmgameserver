package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type ReqBattleBugHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqBattleBugHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if _, ok := innerMsg.MsgData.(*message.C2M_ReqBattleBugData); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role.IsInBattling() || role.GetBattleId() > 0 {
				battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
				if battle == nil {
					log4g.Info("Battle == null")
					return
				}
				log4g.Info("BUg")
				//发送重连消息
				if battle.GetBattleType() == face.FreeRoomBattleType {
					reconnectMsg := &message.M2C_ReBattleBugData{}
					if reconnectMsg.RoomInfo == nil {
						reconnectMsg.RoomInfo = &message.Room{}
					}
					if reconnectMsg.Commands == nil {
						reconnectMsg.Commands = make([]*message.SaveCommandList, 0)
					}
					room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
					if room == nil {
						log4g.Infof("出错:房间id[%d]", role.GetRoomId())
						return
					}
					//房间信息
					handler.InitRoomInfo(reconnectMsg.RoomInfo, room)
					//战斗帧信息
					saveFrames := battle.GetSaveFrames()
					for k, v := range saveFrames {
						sv := &message.SaveCommandList{}
						if sv.Cmds == nil {
							sv.Cmds = make([]*message.Command, 0)
						}
						sv.FrameIndex = k
						sv.Cmds = v
						reconnectMsg.Commands = append(reconnectMsg.Commands, sv)
					}
					reconnectMsg.Seed = room.GetSeed()
					reconnectMsg.FrameCount = battle.GetFrameCount()
					handler.GameServer.WriteInnerMsg(session, innerMsg.RoleId, 5102, reconnectMsg)
				}
			}
		}
	}
}
func (handler *ReqBattleBugHandler) InitRoomInfo(room *message.Room, roomData face.IRoom) {
	if room.Owner == nil {
		room.Owner = &message.RoomOwner{}
	}
	room.Owner.RoomId = roomData.GetRoomId()
	room.Owner.RoomOwnerName = roomData.GetRoomOwnerName()
	room.Owner.RoomOwnerId = roomData.GetRoomOwnerId()
	room.Owner.CityId = roomData.GetRoomOwnerCityId()
	room.Owner.GameType = roomData.GetGameType()
	room.Owner.IsWarFow = roomData.GetIsWarFow()
	room.Owner.RoomOwnerAvatarId = roomData.GetRoomOwnerAvatarId()
	room.Owner.CurPlayeNum = roomData.GetCurPlayerNum()
	room.Owner.MapId = roomData.GetMapId()
	room.Owner.MaxPlayerNum = roomData.GetRoomMaxPlayerNum()
	room.Owner.GameType = roomData.GetGameType()
	room.Owner.RoomOwnerGroupId = roomData.GetRoomOwnerGroupId()
	room.Owner.Arrower = roomData.GetArrowerData()
	room.Owner.Daodun = roomData.GetDaodunData()
	room.Owner.Spear = roomData.GetSpearData()
	room.Owner.Fashi = roomData.GetFashiData()
	//添加存在的队员
	if roomData.GetCurPlayerNum() > 1 {
		if room.Members == nil {
			room.Members = make([]*message.RoomMember, 0)
		}
		roleIds := roomData.GetRoomRoleIds()
		for _, mem := range roleIds {
			if mem > 0 {
				if mem == roomData.GetRoomOwnerId() {
					continue
				}
				groupId := roomData.GetRoomMemberGroupId(mem)
				ready := roomData.GetRoomMemberReady(mem)
				role := handler.GameServer.GetRoleManager().GetOnlineRole(mem)
				if role != nil {
					msgMem := &message.RoomMember{}
					msgMem.GroupId = groupId
					msgMem.JoinerLevel = role.GetLevel()
					msgMem.JoinerIconId = role.GetAvatarId()
					msgMem.JoinerName = role.GetNickName()
					msgMem.CityId = roomData.GetRoomMemberCityId(mem)
					msgMem.JoinerId = mem
					msgMem.Arrower = role.GetSoldierData(0)
					msgMem.Daodun = role.GetSoldierData(1)
					msgMem.Spear = role.GetSoldierData(2)
					msgMem.Fashi = role.GetSoldierData(3)
					msgMem.Ready = ready
					room.Members = append(room.Members, msgMem)
				}
			}
		}
	}
}
