package msgHandler

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"time"
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
			role := handler.GameServer.GetRoleManager().GetOnlineRole(protoMsg.RoleId)
			//看是否有订单没有完成
			role.DoNoCompleteTran()
			//看登录时间和现在登录时间差，满足1小时获得2钻石
			now := time.Now()
			lastLoginTime := role.GetLoginTime()
			delatTime := now.Sub(lastLoginTime)
			allHours := delatTime.Hours()
			if allHours >= 1{
				//如果超过三天没登陆就没离线奖励
				if allHours >= 3 * 24{
					rMsg.NotOnlineDiam = 0
				}else{
					addDiam := int32(allHours)
					//println(addDiam)
					rMsg.NotOnlineDiam = addDiam
					role.AddDiam(addDiam)//增加离线收益
				}
			}
			if role.GetOnlineDiam() > 0{
				role.AddDiam(role.GetOnlineDiam())
				role.SetOnlineDiam(0)
			}
			role.SetLoginTime(now)
			rMsg.RoleBasicInfo.RoleId = protoMsg.RoleId
			rMsg.RoleBasicInfo.GuideFinished = role.GetGuideFinised()
			rMsg.RoleBasicInfo.NickName = role.GetNickName()
			rMsg.RoleBasicInfo.Qq = role.GetQQ()
			rMsg.RoleBasicInfo.Weixin = role.GetWeiXin()
			rMsg.RoleBasicInfo.Phone = role.GetPhone()
			rMsg.RoleBasicInfo.MaillAddress = role.GetMailAddresss()
			rMsg.RoleBasicInfo.Gold = role.GetGold()
			rMsg.RoleBasicInfo.AvatarId = role.GetAvatarId()
			rMsg.RoleBasicInfo.Exp = role.GetExp()
			rMsg.RoleBasicInfo.Diam = role.GetDiam()
			rMsg.RoleBasicInfo.Level = role.GetLevel()
			rMsg.RoleBasicInfo.Sex = role.GetSex()
			rMsg.RoleBasicInfo.Sign = role.GetSign()
			rMsg.RoleBasicInfo.MaxBagNum = int32(role.GetMaxBagNum())
			rMsg.RoleBasicInfo.RankScore = role.GetRankScore()
			rMsg.Achieve = role.GetAchieveMsgData()
			rMsg.RoomId = role.GetRoomId()
			rMsg.BattleId = role.GetBattleId()
			rMsg.Version = handler.GameServer.GetGameVersion()
			rMsg.TaskSeed = role.GetTaskSeed()
			rMsg.RoleBasicInfo.SaijiId = role.GetSaiJiId()
			_,rMsg.OnlineDiamLeft = handler.GameServer.SubAllOnlineDiam(0)
			//Hero
			heroMap := role.GetAllHero()
			if len(heroMap) > 0 {
				for _, v := range heroMap {
					hero := &message.Hero{}
					hero.HeroId = v.HeroId
					hero.Level = v.Level
					hero.SkillPoint = v.SkillPoint
					hero.SkillId1 = v.Skill1
					hero.SkillId2 = v.Skill2
					for _, v1 := range v.ItemIds {
						if v1.ItemId > 0 {
							item := &message.Item{}
							item.ItemId = v1.ItemId
							item.ItemSeed = v1.ItemSeed
							item.ItemNum = 1
							hero.ItemIds = append(hero.ItemIds, item)
						}
					}
					for _,v2 := range v.LearnSkill{
						if v2 > 0{
							hero.LearnSkillIds = append(hero.LearnSkillIds, v2)
						}
					}
					rMsg.HeroInfo = append(rMsg.HeroInfo, hero)
				}
			}
			//更新武将数量
			role.SetHeroCount(int32(len(rMsg.HeroInfo)))
			//背包
			for i := 0; i < role.GetMaxBagNum(); i++ {
				itemIndex := int32(i)
				te := role.GetItem(itemIndex)
				if te.ItemId > 0 {
					item := &message.Item{}
					item.Index = itemIndex
					item.ItemId = te.ItemId
					item.ItemSeed = te.ItemSeed
					item.ItemNum = te.ItemNum
					rMsg.Items = append(rMsg.Items, item)
				}
			}
			//初始化邮件
			if role.GetCurMaxEmailNum() > role.GetMaxEmailNum(){
				//说明有新邮件
				for i:=role.GetMaxEmailNum();i<role.GetCurMaxEmailNum();i++{
					email := role.GetEmailByIndex(i)
					//if email != nil{
					//	value := role.AddEmail(email)
					//	if value < 0{
					//		//已经满了
					//		break
					//	}
					//}
					value := role.AddEmail(email)
					if value < 0{
						//已经满了
						break
					}
				}
				role.DeleteEmailReally()
			}
			for i := 0; i < role.GetMaxEmailNum(); i++ {
				email := role.GetEmailByIndex(i)
				if &email == nil{
					role.AddEmptyEmail(i)
				}else{
					//if email.EmailId > 0{
					//	emailMsg := email.ToMessageData()
					//	rMsg.Emails = append(rMsg.Emails, emailMsg)
					//}
					if email.Valid{
						emailMsg := email.ToMessageData()
						rMsg.Emails = append(rMsg.Emails, emailMsg)
					}
				}
			}
			//初始化每日任务领取过的id
			taskList := role.GetDayTaskIdList()
			rMsg.DayHasGetTaskId = taskList
			//成就
			rMsg.HasGetAchievementId = role.HasGetAchievementId()
			//初始化功能卡
			allCard := role.GetCardIds()
			if allCard != nil{
				rMsg.CardType = *allCard
			}
			cardDiam := role.GetCardDayGetDiam()
			if cardDiam > 0{
				//发送
				msg := new(message.M2C_CardAward)
				msg.CardDiamValue = cardDiam
				role.AddDiam(cardDiam)
				handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5095, msg)
				role.SetCardDayGetDiam(0)
			}
			//聊天记录
			chat := *handler.GameServer.GetChatInfo()
			for i:=0;i<10;i++{
				temp := &chat[i]
				if temp.AvatarId != 0{
					rMsg.Chats = append(rMsg.Chats, temp)
				}
			}
			//初始化自由士兵配置
			rMsg.Arrower = role.GetSoldierData(0)
			rMsg.Daodun = role.GetSoldierData(1)
			rMsg.Spear = role.GetSoldierData(2)
			rMsg.Fashi = role.GetSoldierData(3)
			rMsg.DunJiaFreeSoldier = *role.GetDunJiaSoldierData()
			rMsg.BuJiaFreeSoldier = *role.GetBuJiaSoldierData()
			handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5000, rMsg)
			if role.IsInBattling() || role.GetBattleId() > 0 {
				battle := handler.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
				if battle == nil{
					log4g.Info("Battle == null")
					return
				}
				//log4g.Info("重连战斗")
				//发送重连消息
				if battle.GetBattleType() == face.FreeRoomBattleType{
					reconnectMsg := &message.M2C_ReBattleConnect{}
					if reconnectMsg.RoomInfo == nil {
						reconnectMsg.RoomInfo = &message.Room{}
					}
					if reconnectMsg.Commands == nil {
						reconnectMsg.Commands = make([]*message.SaveCommandList, 0)
					}
					room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
					if room == nil{
						log4g.Infof("出错:房间id[%d]",role.GetRoomId())
						role.SetInRooming(false)
						role.SetInBattling(false)
						role.SetRoomId(0)
						role.SetBattleId(0)
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
				}else{
					reconnectMsg := &message.M2C_RePaiWeiBattleConnect{}
					if reconnectMsg.Commands == nil {
						reconnectMsg.Commands = make([]*message.SaveCommandList, 0)
					}
					if reconnectMsg.Players == nil{
						reconnectMsg.Players = make([]*message.MatchRoomPlayer,0)
					}
					room := handler.GameServer.GetMatchManager().GetMatchBattleRoom(role.GetRoomId())
					if room != nil{
						reconnectMsg.RoomId = room.GetRoomId()
						reconnectMsg.Seed = room.GetRoomSeed()
						allPlayer := room.GetAllPlayer()
						for _,v := range allPlayer{
							if v != nil{
								tempRole := v.GetOnlineRole()
								msgPlayer := new(message.MatchRoomPlayer)
								msgPlayer.RoleId = v.GetRoleId()
								msgPlayer.GroupId = v.GetGroupId()
								msgPlayer.PosIndex = v.GetPosIndex()
								msgPlayer.AvatarId = tempRole.GetAvatarId()
								msgPlayer.NickName = tempRole.GetNickName()
								msgPlayer.CityId = 1
								msgPlayer.Arrower = tempRole.GetSoldierData(0)
								msgPlayer.Daodun = tempRole.GetSoldierData(1)
								msgPlayer.Spear = tempRole.GetSoldierData(2)
								msgPlayer.Fashi = tempRole.GetSoldierData(3)
								reconnectMsg.Players = append(reconnectMsg.Players, msgPlayer)
							}
						}
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
						reconnectMsg.FrameCount = battle.GetFrameCount()
						handler.GameServer.WriteInnerMsg(session, protoMsg.RoleId, 5075, reconnectMsg)
					}else{
						log4g.Infof("不存在匹配房间[%d]",role.GetRoomId())
					}
				}
			}else if role.IsInRooming(){
				room := handler.GameServer.GetRoomManager().GetRoomByRoomId(role.GetRoomId())
				if room == nil{
					log4g.Info("出错1")
					return
				}
				//log4g.Info("重连房间")
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
				roomReconnect.Load = room.GetInBattle()
				roomReconnect.Seed = room.GetSeed()
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
