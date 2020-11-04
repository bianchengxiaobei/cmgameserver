package roleManager

import (
	"cmgameserver/bean"
	"cmgameserver/face"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
	"cmgameserver/tool"
)

type RoleManager struct {
	lock               sync.RWMutex
	onlineRoles        map[int64]face.IOnlineRole                                    //[roleId]role
	onlineSessionRoles map[network.SocketSessionInterface]map[int64]face.IOnlineRole //[session][roleId]Role
	GameServer         face.IGameServer
}

func NewRoleManager(server face.IGameServer) *RoleManager {
	return &RoleManager{
		onlineRoles:        make(map[int64]face.IOnlineRole),
		onlineSessionRoles: make(map[network.SocketSessionInterface]map[int64]face.IOnlineRole),
		GameServer:         server,
	}
}
func (manager *RoleManager) GetOnlineRole(roleId int64) face.IOnlineRole {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.onlineRoles[roleId]
}
func (manager *RoleManager) AddOnlineRole(role face.IOnlineRole) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.onlineRoles[role.GetRoleId()] = role
	session := role.GetGateSession()
	roleMap := manager.onlineSessionRoles[session]
	if roleMap == nil {
		roleMap = make(map[int64]face.IOnlineRole)
		manager.onlineSessionRoles[session] = roleMap
	}
	roleMap[role.GetRoleId()] = role
}
func (manager *RoleManager) RemoveOnlineRole(role face.IOnlineRole) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.onlineRoles, role.GetRoleId())
	delete(manager.onlineSessionRoles[role.GetGateSession()], role.GetRoleId())
}
func (this *RoleManager)GetOnlineNum()int{
	return len(this.onlineRoles)
}
func (manager *RoleManager) NewOnlineRole(roleId int64) face.IOnlineRole {
	var err error
	if manager.GameServer.GetDBManager() != nil {
		dbSession := manager.GameServer.GetDBManager().Get()
		if dbSession != nil {
			role := bean.Role{}
			var heros []bean.Hero
			c := dbSession.DB("sanguozhizhan").C("Role")
			err = c.Find(bson.M{"roleid": roleId}).One(&role)
			if err != nil {
				return nil
			}
			c = dbSession.DB("sanguozhizhan").C("Hero")
			err = c.Find(bson.M{"roleid": roleId}).All(&heros)
			if err != nil {
				return nil
			}
			//初始化FreeSoldier
			if role.DunJiaDaoDun.TouKuiId == 0{
				role.DunJiaDaoDun.TouKuiId = 2001//盾甲官方头盔
				role.DunJiaDaoDun.PlayerType = 1
				role.DunJiaDaoDun.WeapId = 4001//官方刀
				role.DunJiaDaoDun.CarrierType = 1
				role.DunJiaArrow.TouKuiId = 2001
				role.DunJiaArrow.WeapId = 5001//官方弓箭
				role.DunJiaArrow.PlayerType = 2
				role.DunJiaArrow.CarrierType = 1
				role.DunJiaSpear.TouKuiId = 2001
				role.DunJiaSpear.WeapId = 8001//官方枪
				role.DunJiaSpear.PlayerType = 7
				role.DunJiaSpear.CarrierType = 1
				role.DunJiaFaShi.TouKuiId = 2001
				role.DunJiaFaShi.WeapId = 7001//官方法杖
				role.DunJiaFaShi.PlayerType = 10
				role.DunJiaFaShi.CarrierType = 1
			}else{
				role.DunJiaDaoDun.CarrierType = 1
				role.DunJiaArrow.CarrierType = 1
				role.DunJiaSpear.CarrierType = 1
				role.DunJiaFaShi.CarrierType = 1
			}
			if role.BuJiaDaoDun.TouKuiId == 0{
				role.BuJiaDaoDun.TouKuiId = 1//盾甲官方头盔
				role.BuJiaDaoDun.WeapId = 4001//官方刀
				role.BuJiaDaoDun.PlayerType = 1
				role.BuJiaDaoDun.CarrierType = 0
				role.BuJiaArrow.TouKuiId = 1
				role.BuJiaArrow.WeapId = 5001//官方弓箭
				role.BuJiaArrow.PlayerType = 2
				role.BuJiaArrow.CarrierType = 0
				role.BuJiaSpear.TouKuiId = 1
				role.BuJiaSpear.WeapId = 8001//官方枪
				role.BuJiaSpear.PlayerType = 7
				role.BuJiaSpear.CarrierType = 0
				role.BuJiaFaShi.TouKuiId = 1
				role.BuJiaFaShi.WeapId = 7001//官方法杖
				role.BuJiaFaShi.PlayerType = 10
				role.BuJiaFaShi.CarrierType = 0
			}else{
				role.BuJiaDaoDun.CarrierType = 0
				role.BuJiaArrow.CarrierType = 0
				role.BuJiaSpear.CarrierType = 0
				role.BuJiaFaShi.CarrierType = 0
			}
			if role.FreeSoldierData[0].TouKuiId != 0{

			}
			if role.MaxBagNum <= 32{
				role.MaxBagNum = 64
			}
			lenItem := len(role.Items)
			if lenItem == 0 {
				//说明游戏还没有背包
				for i := 0; i < int(role.MaxBagNum); i++ {
					item := bean.Item{}
					item.ItemId = 0
					role.Items = append(role.Items, item)
				}
			}else{
				if lenItem < role.MaxBagNum{
					for i := 0; i< role.MaxBagNum - lenItem;i++{
						item := bean.Item{}
						item.ItemId = 0
						role.Items = append(role.Items, item)
					}
				}
			}
			if len(heros) == 0 {
				//说明游戏还没有英雄，免费送，刚开始
				hero := bean.Hero{}
				hero.RoleId = roleId
				hero.HeroId = 1
				hero.Level = 1
				data := manager.GameServer.GetHeroConfig().Data[hero.HeroId]
				for i := 0; i < 3; i++ {
					item := new(bean.Item)
					hero.ItemIds[i] = *item
				}
				body := &hero.ItemIds[0]
				body.ItemId = data.GuanFangBodyId
				body.ItemNum = 1
				body.ItemSeed = 0
				body = &hero.ItemIds[1]
				body.ItemId = data.GuanFangWeapId
				body.ItemNum = 1
				body.ItemSeed = 0
				body = &hero.ItemIds[2]
				body.ItemId = data.GuanFangShoeId
				body.ItemNum = 1
				body.ItemSeed = 0
				err = c.Insert(&hero)
				if err != nil {
					return nil
				}
				heros = append(heros, hero)
			}
			now := time.Now()
			if tool.IsSameDay(role.LoginTime,now) == false{
				//清空每日任务
				for k, _ := range role.DayGetTask {
					role.DayGetTask[k] = 0
				}
				//role.LoginTime = now
				role.TaskSeed = int32(now.Unix())
				role.GetSign = false
			}
			onlineRole := OnlineRole{
				Role:       role,
				Heros:      make(map[int32]bean.Hero),
				BattleInfo: BattleInfo{},
				Connected:  true,
				PingTime:   time.Now(),
			}
			if len(heros) > 0 {
				for _, v := range heros {
					onlineRole.Heros[v.HeroId] = v
				}
			}
			return &onlineRole
		}
	}
	return nil
}
func (manager *RoleManager) GetAllOnlineRole(gateSession network.SocketSessionInterface) map[int64]face.IOnlineRole {
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.onlineSessionRoles[gateSession]
}
func (manager *RoleManager) AllRoleQuit() {
	defer func() {
		if err := recover(); err != nil {
			log4g.Info(err.(error).Error())
			return
		}
	}()
	log4g.Info("closeServerRoleQuit")
	for _, v := range manager.onlineRoles {
		if v != nil {
			v.UpdateDB(manager.GameServer.GetDBManager())
		}
	}
}
func (manager *RoleManager) RoleQuit(roleId int64) {
	defer func() {
		if err := recover(); err != nil {
			log4g.Info(err.(error).Error())
			return
		}
	}()
	//如果role正在战斗中，判断战斗的进程，比如房间中，或者是正在战斗
	role := manager.GetOnlineRole(roleId)
	//更新登录时间
	role.SetLoginTime(time.Now())
	//看在线获得钻石加入然后清空
	if role.GetOnlineDiam() > 0{
		role.AddDiam(role.GetOnlineDiam())
		role.SetOnlineDiam(0)
	}
	if role != nil {
		role.SetInSimulateBattle(false)
		role.SetLoadFinished(false)//还没有加载
		role.SetConnected(false)
		//更新数据库
		role.UpdateDB(manager.GameServer.GetDBManager())
		matchPlayer := role.GetMatchPlayer()
		if matchPlayer != nil {
			//说明正在匹配中，就退出
			manager.GameServer.GetMatchManager().RemovePlayerFromMatchTeam(matchPlayer)
		}
		if role.IsInBattling() {
			//退出战斗
			role.SetLoadFinished(false)
			battle := manager.GameServer.GetBattleManager().GetBattle(role.GetBattleId())
			if battle != nil {
				battle.CheckAllLeave()
			}
		} else {
			if role.IsInRooming() {
				//退出该房间，发送消息
				roomManger := manager.GameServer.GetRoomManager()
				roomId := role.GetRoomId()
				room := roomManger.GetRoomByRoomId(roomId)
				if room != nil {
					if room.GetInBattle() == false{
						if room.IsRoomOwner(roleId) {
							//如果是房主，删除房间，并且通知所有人
							if roomManger.DeleteRoom(roomId) {
								log4g.Infof("删除房间[%d]成功!", roomId)
							} else {
								log4g.Infof("删除房间[%d]失败!", roomId)
							}
						} else {
							//移除成员，并且通知房间所有人
							if roomManger.RemoveOneMemberByRoom(room, role) == false {
								log4g.Infof("移除房间[%d]成员[%d]失败!", roomId, roleId)
							}
						}
					}else{
						log4g.Infof("房间[%d]战斗中不退出!", roomId)
					}
				}
			}else{
				//移除缓存，如果是战斗中，就不移除，等过游戏结束如果玩家还么有连接上来再移除
				manager.RemoveOnlineRole(role)
			}
		}
		//发送给网关移除角色
		message := &message.M2G_RoleQuitGate{
			RoleId: roleId,
		}
		manager.GameServer.WriteInnerMsg(role.GetGateSession(), roleId, 10005, message)
	} else {
		log4g.Infof("不存在玩家[%d],退出游戏服!", roleId)
	}
}
//全服获得月卡
func (this *RoleManager) SendCardDiamToAllRole(){
	session := this.GameServer.GetDBManager().Get()
	now := time.Now()
	if session != nil {
		c := session.DB("sanguozhizhan").C("Role")
		if c != nil {
			var roles []bean.Role
			err := c.Find(nil).All(&roles)
			if err != nil {
				log4g.Infof("数据库出现错误:[%s]", err.Error())
				return
			}
			for _, v := range roles {
				online := this.onlineRoles[v.RoleId]
				if online != nil || v.IsGM {
					continue
				}
				//看是否月卡能发
				if len(v.Card) > 0{
					var add int32
					for k, card := range v.Card{
						if card.CardType == bean.NoneCard{
							continue
						}
						if card.IsTimeOut(now){
							//过期移除数据库
							err = c.Update(bson.M{"roleid": v.RoleId}, bson.M{"$pull": bson.M{
								"card": card,
							}})
							if err != nil {
								log4g.Infof("Update数据库出现错误:[%s]", err.Error())
							}
							continue
						}
						//发送每日奖励
						add += card.GetDayDiam(now,this.GameServer.GetCardInfoConfig())
						v.Card[k] = card
					}
					if add == 0{
						continue
					}
					err = c.Update(bson.M{"roleid": v.RoleId}, bson.M{"$inc": bson.M{
						"carddiamvalue": add,
					}})
					if err != nil {
						log4g.Infof("Update数据库出现错误:[%s]", err.Error())
						continue
					}
					err = c.Update(bson.M{"roleid": v.RoleId}, bson.M{"$set": bson.M{
						"card": v.Card,
					}})
					if err != nil {
						log4g.Infof("Update数据库出现错误:[%s]", err.Error())
						continue
					}
				}
			}
		}else{
			return
		}
	}else{
		return
	}
	for _, v := range this.onlineRoles {
		if v != nil {
			if v.GetIsGM() {
				continue
			}
			this.GetCardDiamByOnlineRole(now,v)
		}
	}
}
func (this *RoleManager)GetCardDiamByOnlineRole(now time.Time,role face.IOnlineRole){
	var add int32
	cards := *role.GetCards()
	if len(cards) > 0{
		for k, card := range cards{
			if card.CardType == bean.NoneCard{
				continue
			}
			if card.IsTimeOut(now){
				role.DeleteCard(k)
				continue
			}
			add += card.GetDayDiam(now,this.GameServer.GetCardInfoConfig())
			role.SetCard(card,k)
		}
		role.AddCardDayGetDiam(add)
		value := role.GetCardDayGetDiam()
		if value > 0{
			msg := new(message.M2C_CardAward)
			msg.CardDiamValue = value
			role.AddDiam(value)
			this.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5095, msg)
			role.SetCardDayGetDiam(0)
		}
	}
}
///向全服的玩家发送邮件
func (this *RoleManager) SendEmailToAllRole(title string,
	content string,awardList []int32)bool {
	time := time.Now().Unix()
	email := bean.Email{
		//EmailId:   emailId,
		Title:title,
		Content:content,
		EmailTime: time,
		AwardList:awardList,
		Get:       false,
	}
	session := this.GameServer.GetDBManager().Get()
	if session != nil {
		c := session.DB("sanguozhizhan").C("Role")
		if c != nil {
			var roles []bean.Role
			err := c.Find(nil).All(&roles)
			if err != nil {
				log4g.Infof("数据库出现错误:[%s]", err.Error())
				return false
			}
			for _, v := range roles {
				online := this.onlineRoles[v.RoleId]
				if online != nil || v.IsGM {
					continue
				}
				err = c.Update(bson.M{"roleid": v.RoleId}, bson.M{"$push": bson.M{
					"emails": email,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
			}
		}else{
			return false
		}
	}else{
		return false
	}
	for k, v := range this.onlineRoles {
		if v != nil {
			if v.GetIsGM() {
				continue
			}
			msg := new(message.M2C_AddEmail)
			msg.Email = email.ToMessageData()
			index := v.AddEmail(email)
			msg.Email.EmailIndex = index
			this.GameServer.WriteInnerMsg(v.GetGateSession(), k, 5051, msg)
		}
	}
	return true
}
func (this *RoleManager) SendRollInfoToAllRole(notify string)bool {
	msg := new(message.M2C_RollInfo)
	msg.RollType = message.M2C_RollInfo_Notify
	msg.RoleName = notify
	for k, v := range this.onlineRoles {
		if v != nil {
			if v.GetIsGM() {
				continue
			}
			this.GameServer.WriteInnerMsg(v.GetGateSession(), k, 5060, msg)
		}
	}
	return true
}
//通知全服玩家有人得到好的装备
func (this *RoleManager)SendRollInfoToAllRoleGetItem(itemId int32,roleName string){
	msg := new(message.M2C_RollInfo)
	msg.RollType = message.M2C_RollInfo_GetItem
	msg.RoleName = roleName
	msg.RollValue = itemId
	for k, v := range this.onlineRoles {
		if v != nil {
			if v.GetIsGM() {
				continue
			}
			this.GameServer.WriteInnerMsg(v.GetGateSession(), k, 5060, msg)
		}
	}
}
func (this *RoleManager) SendEmailToOneRole(roleId int64, title string,
	content string,awardList []int32)bool {
	time := time.Now().Unix()
	email := bean.Email{
		Title:   title,
		Content:content,
		AwardList:awardList,
		EmailTime: time,
		Get:       false,
	}
	session := this.GameServer.GetDBManager().Get()
	if session != nil {
		c := session.DB("sanguozhizhan").C("Role")
		if c != nil {
			var role bean.Role
			err := c.Find(bson.M{"roleid": roleId}).One(&role)
			if err != nil {
				log4g.Infof("数据库出现错误:[%s]", err.Error())
				return false
			}
			online := this.onlineRoles[role.RoleId]
			if online == nil{
				err = c.Update(bson.M{"roleid": roleId}, bson.M{"$push": bson.M{
					"emails": email,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
				return true
			}else{
				msg := new(message.M2C_AddEmail)
				index := online.AddEmail(email)
				email.EmailIndex = index
				msg.Email = email.ToMessageData()
				this.GameServer.WriteInnerMsg(online.GetGateSession(), roleId, 5051, msg)
				return true
			}
		}else{
			return false
		}
	}else{
		return false
	}
}
func (this *RoleManager)BuyCard(roleId int64,cardType int32)bool{
	role := this.onlineRoles[roleId]
	if role != nil && role.IsConnected(){
		if role.BuyCard(bean.CardType(cardType)){
			returnMsg := new(message.M2C_BuyCardResult)
			returnMsg.CardType = cardType
			returnMsg.Success = true
			//通知
			this.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5079,returnMsg)
			//通知获得月卡奖励
			this.GetCardDiamByOnlineRole(time.Now(),role)
			return true
		}
		return false
	}else{
		session := this.GameServer.GetDBManager().Get()
		if session != nil{
			c := session.DB("sanguozhizhan").C("Role")
			if c != nil{
				tempCard := bean.Card{
					CardType:bean.CardType(cardType),
					BuyTime:time.Now(),
				}
				err := c.Update(bson.M{"roleid": roleId}, bson.M{"$push": bson.M{
					"card": tempCard,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
			}else{
				return false
			}
		}else{
			return false
		}
	}
	return true
}
func (this *RoleManager) BuyDiamIOS(roleId int64,diamType int32,tranId string)bool{
	role := this.onlineRoles[roleId]
	tran := bean.Tran{
		TranId:tranId,
		TranType:bean.DiamTranType,
		TranValue:diamType,
	}
	if role != nil && role.IsConnected(){
		role.AddTran(tran)
		RMsg := new(message.M2C_BuyDiamResult)
		RMsg.DiamType = diamType
		RMsg.AddDiam = tool.GetDiamValueByDiamType(diamType)
		this.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5081,RMsg)
		return true
	}else{
		session := this.GameServer.GetDBManager().Get()
		if session != nil{
			c := session.DB("sanguozhizhan").C("Role")
			if c != nil{
				err := c.Update(bson.M{"roleid":roleId},bson.M{"$push":bson.M{
					"trans":tran,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
			}else{
				return false
			}
		}else{
			return false
		}
	}
	return true
}
func (this *RoleManager)BuyCardIOS(roleId int64,cardType int32,tranId string)bool{
	role := this.onlineRoles[roleId]
	tran := bean.Tran{
		TranId:tranId,
		TranType:bean.CardTranType,
		TranValue:cardType,
	}
	if role != nil && role.IsConnected(){
		role.AddTran(tran)
		returnMsg := new(message.M2C_BuyCardResult)
		returnMsg.CardType = cardType
		returnMsg.Success = true
		//通知
		this.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5079,returnMsg)
		return true
	}else{
		session := this.GameServer.GetDBManager().Get()
		if session != nil{
			c := session.DB("sanguozhizhan").C("Role")
			if c != nil{
				err := c.Update(bson.M{"roleid": roleId}, bson.M{"$push": bson.M{
					"trans": tran,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
			}else{
				return false
			}
		}else{
			return false
		}
	}
	return true
}
func (this *RoleManager)BuyDiam(roleId int64,diamType int32)bool{
	role := this.onlineRoles[roleId]
	addDiam := tool.GetDiamValueByDiamType(diamType)
	if role != nil && role.IsConnected(){
		RMsg := new(message.M2C_BuyDiamResult)
		RMsg.DiamType = diamType
		RMsg.AddDiam = addDiam
		role.AddDiam(addDiam)//增加钻石
		this.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5081,RMsg)
		return true
	}else{
		session := this.GameServer.GetDBManager().Get()
		if session != nil{
			c := session.DB("sanguozhizhan").C("Role")
			if c != nil{
				err := c.Update(bson.M{"roleid":roleId},bson.M{"$inc":bson.M{
					"diam":addDiam,
				}})
				if err != nil {
					log4g.Infof("Update数据库出现错误:[%s]", err.Error())
					return false
				}
			}else{
				return false
			}
		}else{
			return false
		}
	}
	return true
}

func (this *RoleManager)SendSaiJiAward(){
	//取得前几名
	dbSession := this.GameServer.GetDBManager().Get()
	if dbSession == nil{
		return
	}
	allRole := make([]bean.Role,20)
	c := dbSession.DB("sanguozhizhan").C("Role")
	err := c.Find(nil).Sort("-rankscore").Limit(20).All(&allRole)
	if err != nil {
		log4g.Error("更新RankLevel出错d")
		return
	}
	for k,v := range allRole{
		//前3名获得紫晶翅膀
		if k < 3{
			this.SendEmailToOneRole(v.RoleId,"恭喜获得1","fef",nil)
		}else if k >= 3 && k < 7{
			this.SendEmailToOneRole(v.RoleId,"恭喜获得2","fef",nil)
		}else{
			this.SendEmailToOneRole(v.RoleId,"恭喜获得3","fef",nil)
		}
	}
}
func (this *RoleManager)AddRoleHero(role face.IOnlineRole,hero bean.Hero){
	role.AddHero(hero)
	//更新数据库
	dbSession := this.GameServer.GetDBManager().Get()
	if dbSession != nil {
		c := dbSession.DB("sanguozhizhan").C("Hero")
		if c != nil{
			err := c.Insert(hero)
			if err != nil {
				log4g.Errorf("插入英雄出错[%s],RoleId:%d", err.Error(), role.GetRoleId())
			}
		}
	}
}
func (this *RoleManager)AddRoleExp(role face.IOnlineRole,expValue int32){
	ok,num := role.AddExp(expValue);if ok{
		now := time.Now()
		seed := int32(now.Unix())
		min := role.GetLevel() - num
		for i:= min + 1;i<=role.GetLevel();i++{
			if i > 10{
				if i % 10 == 0{
					role.AddDiam(300)
				}else{
					role.AddGold(1000)
				}
			}else {
				datas := this.GameServer.GetLevelUpgradeConfig().Datas
				data := datas[i]
				for _, v := range data.Awards {
					role.AddItemNoMsg(v, seed, now.Unix(), false)
				}
				if data.Diam > 0 {
					role.AddDiam(data.Diam)
				}
				if data.Gold > 0 {
					role.AddGold(data.Gold)
				}
				if data.HeroId > 0 {
					//看是否已经拥有该武将
					tempHero := role.GetHero(data.HeroId)
					tempItem := bean.Item{
						ItemNum: 1,
					}
					if tempHero == nil {
						//添加武将
						tempHero = &bean.Hero{
							Level:  1,
							HeroId: data.HeroId,
							RoleId: role.GetRoleId(),
						}
						heroData := this.GameServer.GetHeroConfig().Data[data.HeroId]
						tempItem.ItemId = heroData.GuanFangBodyId
						tempHero.ItemIds[0] = tempItem
						tempItem.ItemId = heroData.GuanFangWeapId
						tempHero.ItemIds[1] = tempItem
						tempItem.ItemId = heroData.GuanFangShoeId
						tempHero.ItemIds[2] = tempItem
						this.AddRoleHero(role, *tempHero)
					} else {
						//已经拥有则送钻石
						role.AddDiam(400)
					}
				}
			}
		}
	}
}