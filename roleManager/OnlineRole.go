package roleManager

import (
	"cmgameserver/bean"
	"cmgameserver/face"
	"cmgameserver/matchManager"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"gopkg.in/mgo.v2/bson"
	"sync"
	"time"
	"cmgameserver/tool"
)

type OnlineRole struct {
	Role         bean.Role
	Heros        map[int32]bean.Hero
	GateId       int32
	UserName     string
	BattleInfo   BattleInfo
	Connected    bool
	InSimulateBattle bool
	NextLevelExp int32
	GateSession  network.SocketSessionInterface
	PingTime     time.Time
	//ArrowerSoldierData message.FreeSoldierData
	//DaoDunSoldierData message.FreeSoldierData
	//SpearSoldierData message.FreeSoldierData
	//FashiSoldierData message.FreeSoldierData
	MatchPlayer face.IMatchPlayer
	sync.RWMutex
}
type BattleInfo struct {
	RoomId         int32
	BattleId       int32
	BattleSeed     int32 //战斗Seed
	AgreePause     bool
	IsInBattling   bool
	IsInRooming    bool
	IsLoadFinished bool
}

func (role *OnlineRole) GetRole() *bean.Role {
	return &role.Role
}
func (role *OnlineRole) SetGateSession(session network.SocketSessionInterface) {
	role.GateSession = session
}
func (role *OnlineRole) GetGateSession() network.SocketSessionInterface {
	return role.GateSession
}
func (role *OnlineRole) GetRoleId() int64 {
	return role.Role.RoleId
}
func (role *OnlineRole) GetServerId() int32 {
	return role.Role.ServerId
}
func (role *OnlineRole) GetUserId() int64 {
	return role.Role.UserId
}
func (role *OnlineRole) GetUseName() string {
	return role.UserName
}
func (role *OnlineRole) GetGateId() int32 {
	return role.GateId
}
func (role *OnlineRole)GetOnlineDiam() int32{
	return role.Role.OnlineGetDiam
}
func (role *OnlineRole)SetOnlineDiam(value int32){
	role.Role.OnlineGetDiam = value
}
func (role *OnlineRole) SetRoleId(roleId int64) {
	role.Role.RoleId = roleId
}
func (role *OnlineRole) SetServerId(serverId int32) {
	role.Role.ServerId = serverId
}
func (role *OnlineRole) SetUserId(userId int64) {
	role.Role.UserId = userId
}
func (role *OnlineRole) SetUseName(name string) {
	role.UserName = name
}
func (role *OnlineRole)SetHeroCount(count int32){
	role.Role.HeroCount = count
}
func (role *OnlineRole)GetSaiJiId()int32{
	return role.Role.SaiJiId
}
func (role *OnlineRole)GetInSimulateBattle()bool{
	return role.InSimulateBattle
}
func (role *OnlineRole)SetInSimulateBattle(bValue bool){
	role.InSimulateBattle = bValue
}
func (role *OnlineRole)SetSaiJiId(id int32){
	role.Role.SaiJiId = id
}
func (role *OnlineRole) SetGateId(gateId int32) {
	role.GateId = gateId
}
func (role *OnlineRole) GetNickName() string {
	return role.Role.NickName
}
func (role *OnlineRole) SetNickName(nickName string) {
	role.Role.NickName = nickName
}
func (role *OnlineRole) GetSex() int32 {
	return role.Role.Sex
}
func (role *OnlineRole) SetSex(sex int32) {
	role.Role.Sex = sex
}
func (role *OnlineRole) GetRoomId() int32 {
	return role.BattleInfo.RoomId
}
func (role *OnlineRole) SetRoomId(roomId int32) {
	role.BattleInfo.RoomId = roomId
}
func (role *OnlineRole) SetBattleSeed(seed int32) {
	role.BattleInfo.BattleSeed = seed
}
func (role *OnlineRole) GetBattleSeed() int32 {
	return role.BattleInfo.BattleSeed
}
func (role *OnlineRole) GetLevel() int32 {
	return role.Role.Level
}
func (role *OnlineRole) SetLevel(level int32) {
	role.Role.Level = level
}
func (role *OnlineRole) GetDiam() int32 {
	return role.Role.Diam
}
func (role *OnlineRole) SetDiam(diam int32) {
	role.Role.Diam = diam
}
func (role *OnlineRole) GetGold() int32 {
	return role.Role.Gold
}
func (role *OnlineRole) SetGold(gold int32) {
	role.Role.Gold = gold
}
func (role *OnlineRole) GetExp() int32 {
	return role.Role.Exp
}
func (role *OnlineRole) SetExp(exp int32) {
	role.Role.Exp = exp
}

//取得排位分数
func (role *OnlineRole) GetRankScore() int32 {
	return role.Role.RankScore
}
func (role *OnlineRole) AddRankScore(value int32) int32 {
	role.Role.RankScore += value
	if role.Role.RankScore <= 0 {
		role.Role.RankScore = 0
	}
	return role.Role.RankScore
}
func (role *OnlineRole) GetAvatarId() int32 {
	return role.Role.AvatarId
}
func (role *OnlineRole) SetAvatarId(avatarId int32) {
	role.Role.AvatarId = avatarId
}
func (role *OnlineRole) GetAllHero() map[int32]bean.Hero {
	return role.Heros
}
func (role *OnlineRole) GetMailAddresss() string {
	return role.Role.MailAddress
}
func (role *OnlineRole) SetMailAddresss(address string) {
	role.Role.MailAddress = address
}
func (role *OnlineRole) GetQQ() int32 {
	return role.Role.QQ
}
func (role *OnlineRole) SetQQ(qq int32) {
	role.Role.QQ = qq
}
func (role *OnlineRole) GetWeiXin() string {
	return role.Role.WeiXin
}
func (role *OnlineRole) SetWeiXin(weixin string) {
	role.Role.WeiXin = weixin
}
func (role *OnlineRole) GetPhone() int32 {
	return role.Role.Phone
}
func (role *OnlineRole) SetPhone(phone int32) {
	role.Role.Phone = phone
}

//是否掉线或者连接
func (role *OnlineRole) IsConnected() bool {
	role.RLock()
	defer role.RUnlock()
	return role.Connected
}
func (role *OnlineRole) SetConnected(conn bool) {
	role.Lock()
	defer role.Unlock()
	role.Connected = conn
}
func (role *OnlineRole) IsInBattling() bool {
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInBattling
}
func (role *OnlineRole) SetInBattling(value bool) {
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInBattling = value
}
func (role *OnlineRole)GetGuideFinised()bool{
	return role.Role.GuideFinished
}
func (role *OnlineRole)SetGuideFinised(value bool){
	role.Role.GuideFinished = value
}
func (role *OnlineRole) IsInRooming() bool {
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInRooming
}
func (role *OnlineRole) SetInRooming(value bool) {
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInRooming = value
}
func (role *OnlineRole) IsLoadFinished() bool {
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsLoadFinished
}
func (role *OnlineRole) SetLoadFinished(finished bool) {
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsLoadFinished = finished
}
func (role *OnlineRole) GetBattleId() int32 {
	return role.BattleInfo.BattleId
}
func (role *OnlineRole) SetBattleId(id int32) {
	role.BattleInfo.BattleId = id
}
func (role *OnlineRole) GetAgreePause() bool {
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.AgreePause
}
func (role *OnlineRole) SetAgreePause(agree bool) {
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.AgreePause = agree
}
func (role *OnlineRole) GetSign() string {
	return role.Role.Sign
}
func (role *OnlineRole) SetSign(sign string) {
	role.Role.Sign = sign
}
func (role *OnlineRole) GetMaxBagNum() int {
	return role.Role.MaxBagNum
}
func (role *OnlineRole) SetMaxBagNum(num int) {
	role.Role.MaxBagNum = num
}
func (role *OnlineRole) GetMaxEmailNum() int {
	return role.Role.MaxEmailNum
}
func (role *OnlineRole) SetMaxEmailNum(num int) {
	role.Role.MaxEmailNum = num
}
func (role *OnlineRole)GetLoginTime()time.Time{
	return role.Role.LoginTime
}
func (role *OnlineRole)SetLoginTime(time time.Time){
	role.Role.LoginTime = time
}
func (role *OnlineRole) GetPingTime() time.Time {
	return role.PingTime
}
func (role *OnlineRole) SetPingTime(time time.Time) {
	role.PingTime = time
}
func (role *OnlineRole) AddGold(gold int32) {
	role.Role.Gold += gold
}
func (role *OnlineRole) AddDiam(diam int32) {
	role.Role.Diam += diam
}
func (role *OnlineRole) AddExp(exp int32) (bool,int32){
	role.Role.Exp += exp
	temp := role.Role.Level
	role.UpdateNextExp()
	if role.Role.Exp >= role.NextLevelExp {
		Start:
		role.Role.Level++
		role.UpdateNextExp()
		if role.Role.Exp >= role.NextLevelExp{
			goto Start
		}
		num := role.Role.Level - temp
		return true,num
	}
	return false,0
}
func (role *OnlineRole) UpdateNextExp() {
	role.NextLevelExp = 75*role.Role.Level*role.Role.Level - 25*role.Role.Level
}
func (role *OnlineRole) UpgradeHeroLevel(heroId int32) (level int32,skillPoint int32,value bool) {
	if hero, ok := role.Heros[heroId]; ok {
		level,skillPoint = hero.UpgradeLevel()
		role.Heros[heroId] = hero
		return level,skillPoint, ok
	} else {
		log4g.Infof("Role[%d]不包含该英雄[%d]", role.GetRoleId(), heroId)
		return -1, -1,false
	}
}

func (role *OnlineRole) AddItem(item message.Item, hasNum bool) int32 {
	allNum := int(role.GetMaxBagNum())
	if hasNum == false {
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == 0 {
				roleItem.ItemId = item.ItemId
				roleItem.ItemSeed = item.ItemSeed
				roleItem.ItemNum = item.ItemNum
				roleItem.ItemTime = item.ItemTime
				return int32(i)
			}
		}
	} else {
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == item.ItemId {
				roleItem.ItemSeed = item.ItemSeed
				roleItem.ItemNum++
				roleItem.ItemTime = item.ItemTime
				return int32(i)
			}
		}
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == 0 {
				roleItem.ItemId = item.ItemId
				roleItem.ItemSeed = item.ItemSeed
				roleItem.ItemNum = item.ItemNum
				roleItem.ItemTime = item.ItemTime
				return int32(i)
			}
		}
	}
	return -1
}
func (role *OnlineRole) AddItemNoMsg(itemId int32, itemSeed int32, itemTime int64, hasNum bool) int32 {
	allNum := int(role.GetMaxBagNum())
	if hasNum == false {
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == 0 {
				roleItem.ItemId = itemId
				roleItem.ItemSeed = itemSeed
				roleItem.ItemNum = 1
				roleItem.ItemTime = itemTime
				return int32(i)
			}
		}
	} else {
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == itemId {
				roleItem.ItemId = itemId
				roleItem.ItemSeed = itemSeed
				roleItem.ItemNum++
				roleItem.ItemTime = itemTime
				return int32(i)
			}
		}
		for i := 0; i < allNum; i++ {
			roleItem := &role.Role.Items[i]
			if roleItem.ItemId == 0 {
				roleItem.ItemId = itemId
				roleItem.ItemSeed = itemSeed
				roleItem.ItemNum = 1
				roleItem.ItemTime = itemTime
				return int32(i)
			}
		}
	}
	return -1
}
func (role *OnlineRole) GetItem(index int32) bean.Item {
	if int(index) >= role.GetMaxBagNum() || int(index) < 0{
		nilItem := new(bean.Item)
		return *nilItem
	} else {
		return role.Role.Items[index]
	}
}
func (role *OnlineRole) SetItem(index int32, item bean.Item) {
	if int(index) > role.GetMaxBagNum() {
		return
	} else {
		roleItem := &role.Role.Items[index]
		roleItem.ItemId = item.ItemId
		roleItem.ItemSeed = item.ItemSeed
		roleItem.ItemNum = item.ItemNum
		roleItem.ItemTime = item.ItemTime
	}
}
func (role *OnlineRole)DeleteItem(index int32){
	if int(index) > role.GetMaxBagNum() {
		return
	} else {
		roleItem := &role.Role.Items[index]
		roleItem.ItemId = 0
		roleItem.ItemSeed = 0
		roleItem.ItemNum = 0
		roleItem.ItemTime = 0
		role.Role.Items[index] = *roleItem
	}
}
func (role *OnlineRole)GetGiftCode(code string)bool{
	for _,v := range role.Role.GiftCode{
		if v == code{
			//已经存在，领取失败
			return false
		}
	}
	//添加已经领取过的
	role.Role.GiftCode = append(role.Role.GiftCode, code)
	return true
}
func (role *OnlineRole) GetHero(heroId int32) *bean.Hero {
	if hero, ok := role.Heros[heroId]; ok {
		return &hero
	} else {
		return nil
	}
}
func (role *OnlineRole)SetHero(hero bean.Hero){
	role.Heros[hero.HeroId] = hero
}
func (role *OnlineRole)AddHero(hero bean.Hero){
	role.Heros[hero.HeroId] = hero
	role.Role.HeroCount++
}
func (role *OnlineRole) UpdateDB(manager *db.MongoBDManager) {
	dbSession := manager.Get()
	if dbSession != nil {
		//更新角色
		c := dbSession.DB("sanguozhizhan").C("Role")
		err := c.Update(bson.M{"roleid": role.Role.RoleId}, role.Role)
		if err != nil {
			log4g.Infof("更新Role出错[%s],RoleId:%d", err.Error(), role.Role.RoleId)
			return
		} else {
			log4g.Infof("更新Role数据库成功:[%d]", role.Role.RoleId)
		}
		c = dbSession.DB("sanguozhizhan").C("Hero")
		for _, v := range role.Heros {
			err = c.Update(bson.M{"roleid": role.Role.RoleId, "heroid": v.HeroId}, v)
			if err != nil {
				log4g.Infof("更新Hero数据库出错:[%d]", role.Role.RoleId)
			}
		}
	} else {
		log4g.Info("dbSession == nil")
	}
}
func (role *OnlineRole) QuitBattle() {
	role.BattleInfo.RoomId = 0
	role.BattleInfo.BattleId = 0
	role.SetLoadFinished(false)//战斗结束设置加载没有成功
	role.SetInBattling(false)
	role.SetBattleId(0)
	role.SetInRooming(false)
	role.SetInSimulateBattle(false)
}
func (role *OnlineRole) HasHero(heroId int32) bool {
	if _,ok := role.Heros[heroId];ok{
		return true
	}
	return false
}
//buyType 金币=0,钻石=1
func (role *OnlineRole) BuyHero(heroId int32,buyType int32) bool {
	hasHero := role.HasHero(heroId)
	if hasHero{
		return false
	}
	if buyType == 0{
		if role.GetGold() < 8000 {
			return false
		}
		role.AddGold(-8000)
	}else{
		if role.GetDiam() < 800{
			return false
		}
		role.AddDiam(-800)
	}
	hero := bean.Hero{}
	hero.RoleId = role.Role.RoleId
	hero.HeroId = heroId
	hero.Level = 1
	//这里基础的装备id得添加
	role.Heros[heroId] = hero
	role.Role.HeroCount++
	return true
}
func (role *OnlineRole) WinLevel(level int32) bool {
	if level == 0 {
		return false
	}
	has := false
	for _, v := range role.Role.WinLevel {
		if level == v {
			has = true
			break
		}
	}
	if has == false {
		role.Role.WinLevel = append(role.Role.WinLevel, level)
		return true
	} else {
		return false
	}
}
func (role *OnlineRole) AddGetTaskAward(taskId int32) bool {
	if taskId <= 0 {
		return false
	}
	for k, v := range role.Role.DayGetTask {
		if v == 0 {
			role.Role.DayGetTask[k] = taskId
			return true
		} else if v == taskId {
			return false
		}
	}
	role.Role.DayGetTask = append(role.Role.DayGetTask, taskId)
	return true
}
func (role *OnlineRole) AddGetAchieveAward(achieveId int32) bool {
	if achieveId <= 0 {
		return false
	}
	for _, v := range role.Role.Achievement {
		if v == achieveId {
			log4g.Infof("已经领取过了")
			return false
		}
	}
	role.Role.Achievement = append(role.Role.Achievement, achieveId)
	return true
}
func (role *OnlineRole)HasGetAchievementId()[]int32{
	list := make([]int32,0)
	if len(role.Role.Achievement) > 0{
		for _,v := range role.Role.Achievement{
			list = append(list, v)
		}
	}
	return list
}
//获取签到奖励
func (role *OnlineRole) GetSignAward() bool {
	if role.Role.GetSign == false {
		role.AddGold(100)
		role.Role.GetSign = true
		return true
	} else {
		return false
	}
}
func (role *OnlineRole) GetTaskSeed() int32 {
	return role.Role.TaskSeed
}
func (role *OnlineRole) GetSoldierData(index int) *message.FreeSoldierData {
	data := &message.FreeSoldierData{}
	arrowerData := role.Role.FreeSoldierData[index]
	data.CarrierType = arrowerData.CarrierType
	data.PlayerType = arrowerData.PlayerType
	data.BodyId = arrowerData.BodyId
	data.TouKuiId = arrowerData.TouKuiId
	data.WeapId = arrowerData.WeapId
	return data
}
func (role *OnlineRole)GetDunJiaSoldierData()*[]*message.FreeSoldierData{
	list := make([]*message.FreeSoldierData,0)
	data := &message.FreeSoldierData{}
	dunjiaArrow := &role.Role.DunJiaArrow
	dunjiaArrow.PlayerType = 2
	data.BodyId = dunjiaArrow.BodyId
	data.WeapId = dunjiaArrow.WeapId
	data.TouKuiId = dunjiaArrow.TouKuiId
	data.PlayerType = dunjiaArrow.PlayerType
	list = append(list, data)
	data2 := &message.FreeSoldierData{}
	dunjiaDaoDun := &role.Role.DunJiaDaoDun
	dunjiaDaoDun.PlayerType = 1
	data2.BodyId = dunjiaDaoDun.BodyId
	data2.WeapId = dunjiaDaoDun.WeapId
	data2.TouKuiId = dunjiaDaoDun.TouKuiId
	data2.PlayerType = dunjiaDaoDun.PlayerType
	list = append(list, data2)
	data3 := &message.FreeSoldierData{}
	dunjiaSpear := &role.Role.DunJiaSpear
	dunjiaSpear.PlayerType = 7
	data3.BodyId = dunjiaSpear.BodyId
	data3.WeapId = dunjiaSpear.WeapId
	data3.TouKuiId = dunjiaSpear.TouKuiId
	data3.PlayerType = dunjiaSpear.PlayerType
	list = append(list, data3)
	data4 := &message.FreeSoldierData{}
	dunjiaFashi := &role.Role.DunJiaFaShi
	dunjiaFashi.PlayerType = 10
	data4.BodyId = dunjiaFashi.BodyId
	data4.WeapId = dunjiaFashi.WeapId
	data4.TouKuiId = dunjiaFashi.TouKuiId
	data4.PlayerType = dunjiaFashi.PlayerType
	list = append(list, data4)
	return &list
}
func (role *OnlineRole)GetBuJiaSoldierData()*[]*message.FreeSoldierData{
	list := make([]*message.FreeSoldierData,0)
	data := &message.FreeSoldierData{}
	dunjiaArrow := &role.Role.BuJiaArrow
	dunjiaArrow.PlayerType = 2
	data.BodyId = dunjiaArrow.BodyId
	data.WeapId = dunjiaArrow.WeapId
	data.TouKuiId = dunjiaArrow.TouKuiId
	data.PlayerType = dunjiaArrow.PlayerType
	list = append(list, data)
	data2 := &message.FreeSoldierData{}
	dunjiaDaoDun := &role.Role.BuJiaDaoDun
	dunjiaDaoDun.PlayerType = 1
	data2.BodyId = dunjiaDaoDun.BodyId
	data2.WeapId = dunjiaDaoDun.WeapId
	data2.TouKuiId = dunjiaDaoDun.TouKuiId
	data2.PlayerType = dunjiaDaoDun.PlayerType
	list = append(list, data2)
	data3 := &message.FreeSoldierData{}
	dunjiaSpear := &role.Role.BuJiaSpear
	dunjiaSpear.PlayerType = 7
	data3.BodyId = dunjiaSpear.BodyId
	data3.WeapId = dunjiaSpear.WeapId
	data3.TouKuiId = dunjiaSpear.TouKuiId
	data3.PlayerType = dunjiaSpear.PlayerType
	list = append(list, data3)
	data4 := &message.FreeSoldierData{}
	dunjiaFashi := &role.Role.BuJiaFaShi
	dunjiaFashi.PlayerType = 10
	data4.BodyId = dunjiaFashi.BodyId
	data4.WeapId = dunjiaFashi.WeapId
	data4.TouKuiId = dunjiaFashi.TouKuiId
	data4.PlayerType = dunjiaFashi.PlayerType
	list = append(list, data4)
	return &list
}
func (role *OnlineRole) ChangeFreeSoldierData(soldierType int,data message.FreeSoldierData) bool {
	if soldierType > 3 {
		return false
	}
	var soldierData *bean.FreeSoldierData
	roleData := &role.Role.FreeSoldierData[soldierType]
	roleData.CarrierType = data.CarrierType
	roleData.WeapId = data.WeapId
	roleData.TouKuiId = data.TouKuiId
	roleData.BodyId = data.BodyId
	carryType := data.CarrierType
	//然后
	if soldierType == 0{
		if carryType == 0{
			soldierData = &role.Role.BuJiaArrow
		}else{
			soldierData = &role.Role.DunJiaArrow
		}
	}else if soldierType == 1{
		if carryType == 0{
			soldierData = &role.Role.BuJiaDaoDun
		}else{
			soldierData = &role.Role.DunJiaDaoDun
		}
	}else if soldierType == 2{
		if carryType == 0{
			soldierData = &role.Role.BuJiaSpear
		}else{
			soldierData = &role.Role.DunJiaSpear
		}
	}else if soldierType == 3{
		if carryType == 0{
			soldierData = &role.Role.BuJiaFaShi
		}else{
			soldierData = &role.Role.DunJiaFaShi
		}
	}
	soldierData.TouKuiId = data.TouKuiId
	soldierData.WeapId = data.WeapId
	soldierData.BodyId = data.BodyId
	return true
}
///index=>士兵类型   equipIndex=>ItemType=>body
func (role *OnlineRole) ChangeFreeSoldierEquipId(soldierType int, equipIndex int, equipId int32,carryType int32) {
	if soldierType > 3 {
		return
	}
	if equipId == 0{
		return
	}
	var data *bean.FreeSoldierData
	if soldierType == 0{
		if carryType == 0{
			data = &role.Role.BuJiaArrow
		}else{
			data = &role.Role.DunJiaArrow
		}
	}else if soldierType == 1{
		if carryType == 0{
			data = &role.Role.BuJiaDaoDun
		}else{
			data = &role.Role.DunJiaDaoDun
		}
	}else if soldierType == 2{
		if carryType == 0{
			data = &role.Role.BuJiaSpear
		}else{
			data = &role.Role.DunJiaSpear
		}
	}else if soldierType == 3{
		if carryType == 0{
			data = &role.Role.BuJiaFaShi
		}else{
			data = &role.Role.DunJiaFaShi
		}
	}
	roleData := &role.Role.FreeSoldierData[soldierType]
	if equipIndex == 0 {
		roleData.BodyId = equipId
		roleData.CarrierType = data.CarrierType
		data.BodyId = equipId
	} else if equipIndex == 1 {
		roleData.WeapId = equipId
		roleData.CarrierType = data.CarrierType
		data.WeapId = equipId
	} else if equipIndex == 5 {
		roleData.TouKuiId = equipId
		roleData.CarrierType = data.CarrierType
		data.TouKuiId = equipId
	}
}
func (role *OnlineRole) GetFreeSoldierEquipId(index int, equipIndex int) int32 {
	if index > 3 {
		return -1
	}
	roleData := role.Role.FreeSoldierData[index]
	if equipIndex == 0 {
		return roleData.BodyId
	} else if equipIndex == 1 {
		return roleData.WeapId
	} else if equipIndex == 5 {
		return roleData.TouKuiId
	}
	return -1
}

///取得官方装备id
func (role *OnlineRole) GetFreeSoldierGuangFanEquipId(index int, equipIndex int) (int32, bool) {
	if index > 3 {
		return -1, false
	}
	roleData := role.Role.FreeSoldierData[index]
	id := int32(0)
	isGuangFan := false
	//布甲
	if equipIndex == 0 {
		//Body
		id = 0
		if roleData.BodyId == 0 {
			isGuangFan = true
		}
	} else if equipIndex == 1 {
		//Weap
		//根据职业
		if index == 0 {
			//Arrow
			id = 5001
			if roleData.WeapId == 5001 {
				isGuangFan = true
			}
		} else if index == 1 {
			//刀盾
			id = 4001
			if roleData.WeapId == 4001 {
				isGuangFan = true
			}
		} else if index == 2 {
			//枪兵
			id = 8001
			if roleData.WeapId == 8001 {
				isGuangFan = true
			}
		} else if index == 3 {
			//法师
			id = 7001
			if roleData.WeapId == 7001 {
				isGuangFan = true
			}
		}
	} else if equipIndex == 5 {
		//TouKui
		if roleData.CarrierType == 0 {
			id = 1
			if roleData.TouKuiId == 1 {
				isGuangFan = true
			}
		} else {
			id = 2001
			if roleData.TouKuiId == 2001 {
				isGuangFan = true
			}
		}
	}
	return id, isGuangFan
}
//func (role *OnlineRole) GetEmail(emailId int32) (*bean.Email, int32) {
//	for _, v := range role.Role.Emails {
//		if v.EmailId == emailId {
//			return &v, v.EmailIndex
//		}
//	}
//	return nil, -1
//}
func (role *OnlineRole) SetEmail(emailIndex int, get bool) bool {
	email := &role.Role.Emails[emailIndex]
	if email != nil {
		email.Get = get
		return true
	}
	return false
}
func (role *OnlineRole)GetEmailByIndex(index int)bean.Email{
	email := role.Role.Emails[index]
	return email
}
func (role *OnlineRole)GetCurMaxEmailNum()int{
	return len(role.Role.Emails)
}
func (role *OnlineRole) DeleteEmail(index int) bool {
	email := role.Role.Emails
	tempEmail := email[index]
	if tempEmail.Valid == false {
		return false
	}
	tempEmail.Clear()
	role.Role.Emails[index] = tempEmail
	return true
}
func (role *OnlineRole)DeleteEmailReally(){
	role.Role.Emails = role.Role.Emails[:role.GetMaxEmailNum()]
}
func (role *OnlineRole) AddEmail(email bean.Email) int32{
	len := len(role.Role.Emails)
	if len < 10{
		email := *new(bean.Email)
		email.Valid = true
		for i:= len;i<10;i++{
			role.Role.Emails = append(role.Role.Emails, email)
		}
	}
	for i:=0;i<10;i++{
		temp := role.Role.Emails[i]
		if temp.Valid {
			continue
		}
		var index int32
		index = int32(i)
		temp.EmailIndex = index
		temp.Title = email.Title
		temp.Content = email.Content
		temp.AwardList = email.AwardList
		temp.Get = email.Get
		temp.EmailTime = email.EmailTime
		role.Role.Emails[i] = temp
		return index
	}
	//已经满了，自动删除邮件
	for i:=0;i<10;i++{
		temp := role.Role.Emails[i]
		if temp.Get{
			//已经领取过的
			var index int32
			index = int32(i)
			temp.EmailIndex = index
			temp.Title = email.Title
			temp.Content = email.Content
			temp.AwardList = email.AwardList
			temp.Get = email.Get
			temp.EmailTime = email.EmailTime
			role.Role.Emails[i] = temp
			return index
		}
	}
	//如果都没有领取过，就放弃该邮件
	return -1
}
func (role *OnlineRole)AddEmptyEmail(index int){
	empty := new(bean.Email)
	empty.Valid = false
	len := len(role.Role.Emails)
	if index >= len{
		for i:=len-1;i<len;i++{
			role.Role.Emails = append(role.Role.Emails, *empty)
		}
	}
}
func (role *OnlineRole) GetMatchPlayer() face.IMatchPlayer {
	if role.MatchPlayer == nil {
		//创建
		p := new(matchManager.MatchPlayer)
		p.OnlineRole = role
		p.IsOwner = false
		p.IsPunish = false
		role.MatchPlayer = p
		return role.MatchPlayer
	} else {
		return role.MatchPlayer
	}
}
func (role *OnlineRole) GetAchieveMsgData() *message.AchievementData {
	data := new(message.AchievementData)
	achieve := role.Role.AchieveRecord
	//data.AllGameNum = achieve.AllGameNum
	data.AllGameTime = achieve.AllGameTime
	data.ConditionType = achieve.ConditionType
	data.ConditionValue = achieve.ConditionValue
	//data.BattleGameFailedNum = achieve.BattleGameFailedNum
	//data.BeDemageNum = achieve.BeDemageNum
	//data.BattleGameWinNum = achieve.BattleGameWinNum
	//data.BreadNum = achieve.BreadNum
	//data.CerealNum = achieve.CerealNum
	//data.DemageNum = achieve.DemageNum
	//data.FailedNum = achieve.FailedNum
	//data.HighestRankLevel = achieve.HighestRankLevel
	//data.KillBuildNum = achieve.KillBuildNum
	//data.KillDockNum = achieve.KillDockNum
	//data.KillFarmerNum = achieve.KillFarmerNum
	//data.KillHeroNum = achieve.KillHeroNum
	//data.KillLadderNum = achieve.KillLadderNum
	//data.KillMangoneNum = achieve.KillMangoneNum
	//data.KillSoldierNum = achieve.KillSoldierNum
	//data.KillStoreCarNum = achieve.KillStoreCarNum
	//data.MeatNum = achieve.MeatNum
	//data.MineNum = achieve.MineNum
	//data.PaiWeiFailedNum = achieve.PaiWeiFailedNum
	//data.PaiWeiWinNum = achieve.PaiWeiWinNum
	//data.RoomFailedNum = achieve.RoomFailedNum
	//data.RoomWinNum = achieve.RoomWinNum
	//data.SimulateFailedNum = achieve.SimulateFailedNum
	//data.SimulateWinNum = achieve.SimulateWinNum
	//data.TreeNum = achieve.TreeNum
	//data.WineNum = achieve.WineNum
	//data.WinNum = achieve.WinNum
	return data
}
func (role *OnlineRole) SetAchieveRecord(data *message.AchievementData) {
	achieve := &role.Role.AchieveRecord
	//achieve.AllGameNum = data.AllGameNum
	achieve.AllGameTime = data.AllGameTime
	achieve.ConditionType = data.ConditionType
	achieve.ConditionValue = data.ConditionValue
	//achieve.BattleGameFailedNum = data.BattleGameFailedNum
	//achieve.BeDemageNum = data.BeDemageNum
	//achieve.BattleGameWinNum = data.BattleGameWinNum
	//achieve.BreadNum = data.BreadNum
	//achieve.CerealNum = data.CerealNum
	//achieve.DemageNum = data.DemageNum
	//achieve.FailedNum = data.FailedNum
	//achieve.HighestRankLevel = data.HighestRankLevel
	//achieve.KillBuildNum = data.KillBuildNum
	//achieve.KillDockNum = data.KillDockNum
	//achieve.KillFarmerNum = data.KillFarmerNum
	//achieve.KillHeroNum = data.KillHeroNum
	//achieve.KillLadderNum = data.KillLadderNum
	//achieve.KillMangoneNum = data.KillMangoneNum
	//achieve.KillSoldierNum = data.KillSoldierNum
	//achieve.KillStoreCarNum = data.KillStoreCarNum
	//achieve.MeatNum = data.MeatNum
	//achieve.MineNum = data.MineNum
	//achieve.PaiWeiFailedNum = data.PaiWeiFailedNum
	//achieve.PaiWeiWinNum = data.PaiWeiWinNum
	//achieve.RoomFailedNum = data.RoomFailedNum
	//achieve.RoomWinNum = data.RoomWinNum
	//achieve.SimulateFailedNum = data.SimulateFailedNum
	//achieve.SimulateWinNum = data.SimulateWinNum
	//achieve.TreeNum = data.TreeNum
	//achieve.WineNum = data.WineNum
	//achieve.WinNum = data.WinNum
}
func (role *OnlineRole) BuyCard(cardType bean.CardType) bool {
	card := &role.Role.Card
	now := time.Now()
	if card == nil {
		tempCard := bean.Card{
			CardType: cardType,
			BuyTime:  now,
			LastGetTime:now,
		}
		role.Role.Card = append(role.Role.Card, tempCard)
	} else {
		for _, v := range role.Role.Card {
			if &v != nil {
				if v.CardType == cardType {
					//说明已经存在，购买失败
					return false
				}
			}
		}
		tempCard1 := bean.Card{
			CardType: cardType,
			LastGetTime:now,
			BuyTime:now,
		}
		role.Role.Card = append(role.Role.Card, tempCard1)
	}
	return true
}
func (role *OnlineRole) GetDayTaskIdList() []int32 {
	return role.Role.DayGetTask
}
func (role *OnlineRole) GetCardIds()*[]int32{
	len := len(role.Role.Card)
	if len > 0{
		cards := make([]int32,0)
		for _,v := range role.Role.Card{
			if v.CardType != bean.NoneCard{
				//cards[k] = int32(v.CardType)
				cards = append(cards, int32(v.CardType))
			}
		}
		return &cards
	}else{
		return nil
	}
}
func (role *OnlineRole)GetCards()*[]bean.Card{
	return &role.Role.Card
}
func (role *OnlineRole)SetCard(card bean.Card,index int){
	role.Role.Card[index] = card
}
func (role *OnlineRole)DeleteCard(index int){
	card := role.Role.Card[index]
	card.Clear()
	role.Role.Card[index] = card
}
func (role *OnlineRole) GetIsGM() bool {
	return role.Role.IsGM
}
func (role *OnlineRole) BuyDiam(diamType bean.DiamType) {
	value := tool.GetDiamValueByDiamType(int32(diamType))
	role.AddDiam(value)
}
func (role *OnlineRole)AddTran(tran bean.Tran){
	role.Role.Trans = append(role.Role.Trans, tran)
	role.Role.Trans = append(role.Role.Trans, tran)
}
func (role *OnlineRole)RemoveTran(tranId string)(bool,bean.TranType,int32){
	trans := role.Role.Trans
	if len(trans) == 0{
		return false,0,0
	}
	tranType := bean.CardTranType
	var tranValue int32
	for k,v := range trans{
		if v.TranId == tranId{
			tranType = v.TranType
			tranValue = v.TranValue
			removeT := append(trans[:k],trans[k+1:]...)
			role.Role.Trans = removeT
			break
		}
	}
	return true,tranType,tranValue
}
//处理未完成的订单
func (role *OnlineRole)DoNoCompleteTran(){
	if len(role.Role.Trans) > 0{
		for _,v := range role.Role.Trans{
			if v.TranType == bean.DiamTranType{
				addDiam := tool.GetDiamValueByDiamType(v.TranValue)
				role.AddDiam(addDiam)
			}else {
				if role.BuyCard(bean.CardType(v.TranValue)) == false{
					log4g.Infof("增加功能卡失败:[%s]",v.TranId)
				}
			}
		}
	}
	role.Role.Trans = make([]bean.Tran,0)
}
func (role *OnlineRole)AddCardDayGetDiam(add int32){
	role.Role.CardDiamValue += add
}
func (role *OnlineRole)GetCardDayGetDiam()int32{
	return  role.Role.CardDiamValue
}
func (role *OnlineRole)SetCardDayGetDiam(value int32){
	role.Role.CardDiamValue = value
}