package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type RoleRegisterGateHandler struct {
	GameServer face.IGameServer
}

func (handler *RoleRegisterGateHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.G2M_RoleRegisterToGateSuccess); ok {
			rMsg := new(message.M2C_EnterLobby)
			if rMsg.RoleBasicInfo == nil {
				rMsg.RoleBasicInfo = new(message.Role)
			}
			rMsg.RoleBasicInfo.RoleId = protoMsg.RoleId
			role := handler.GameServer.GetRoleManager().GetOnlineRole(protoMsg.RoleId)
			rMsg.RoleBasicInfo.NickName = role.GetNickName()
			rMsg.RoleBasicInfo.Gold = role.GetGold()
			rMsg.RoleBasicInfo.AvatarId = role.GetAvatarId()
			rMsg.RoleBasicInfo.Exp = role.GetExp()
			rMsg.RoleBasicInfo.Diam = role.GetDiam()
			rMsg.RoleBasicInfo.Level = role.GetLevel()
			rMsg.RoleBasicInfo.MaxBagNum = role.GetMaxBagNum()
			rMsg.RoomId = role.GetRoomId()
			rMsg.BattleId = role.GetBattleId()
			rMsg.Version = handler.GameServer.GetGameVersion()
			rMsg.TaskSeed = role.GetTaskSeed()
			heroMap := role.GetAllHero()
			if len(heroMap) > 0 {
				for _, v := range heroMap {
					hero := &message.Hero{}
					hero.HeroId = v.HeroId
					hero.Level = v.Level
					for _, v1 := range v.ItemIds {
						if v1 > 0 {
							item := &message.Item{}
							item.ItemId = v1
							hero.ItemIds = append(hero.ItemIds, item)
						}
					}
					rMsg.HeroInfo = append(rMsg.HeroInfo, hero)
				}
			}
			for i := 0; i < int(role.GetMaxBagNum()); i++ {
				itemIndex := int32(i)
				te := role.GetItem(itemIndex)
				if te.ItemId > 0 {
					item := &message.Item{}
					item.Index = itemIndex
					item.ItemId = te.ItemId
					rMsg.Items = append(rMsg.Items, item)
				}
			}
			//初始化自由士兵配置
			rMsg.Arrower = role.GetSoldierData(0)
			rMsg.Daodun = role.GetSoldierData(1)
			rMsg.Spear = role.GetSoldierData(2)
			rMsg.Fashi = role.GetSoldierData(3)
			handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5000, rMsg)
			if role.IsInBattling() || role.GetBattleId() > 0 {
				//发送重连消息
				reconnectMsg := &message.M2C_ReBattleConnect{}
				if reconnectMsg.RoomInfo == nil {
					reconnectMsg.RoomInfo = &message.Room{}
				}
				if reconnectMsg.Commands == nil {
					reconnectMsg.Commands = make([]*message.SaveCommandList, 0)
				}
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
				if room == nil || battle == nil {
					log4g.Info("出错")
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
				handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5028, reconnectMsg)
			}else if role.IsInRooming(){
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room == nil{
					log4g.Info("出错1")
					return
				}
				roomReconnect := &message.M2C_ReRoomConnect{}
				if roomReconnect.RoomInfo == nil{
					roomReconnect.RoomInfo = &message.Room{}
					if roomReconnect.RoomInfo.Owner == nil{
						roomReconnect.RoomInfo.Owner = &message.RoomOwner{}
					}
					if roomReconnect.RoomInfo.Members == nil{
						roomReconnect.RoomInfo.Members = make([]*message.RoomMember,0)
					}
				}
				handler.InitRoomInfo(roomReconnect.RoomInfo,room)
				handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5029, roomReconnect)
			}
		}
	}
}
func (handler *RoleRegisterGateHandler) InitRoomInfo(room *message.Room, roomData face.IRoom) {
	if room.Owner == nil {
		room.Owner = &message.RoomOwner{}
	}
	room.Owner.RoomId = roomData.GetRoomId()
	room.Owner.RoomOwnerName = roomData.GetRoomOwnerName()
	room.Owner.RoomOwnerId = roomData.GetRoomOwnerId()
	room.Owner.RoomName = roomData.GetRoomName()
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
