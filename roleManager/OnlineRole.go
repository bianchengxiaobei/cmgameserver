package roleManager

import (
	"cmgameserver/bean"
	"github.com/bianchengxiaobei/cmgo/network"
	"sync"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/db"
	"gopkg.in/mgo.v2/bson"
	"time"
	"cmgameserver/message"
)


type OnlineRole struct {
	Role     bean.Role
	Heros	map[int32]bean.Hero
	GateId   int32
	UserName string
	BattleInfo	BattleInfo
	Connected	bool
	NextLevelExp int32
	GateSession	network.SocketSessionInterface
	PingTime	time.Time
	//ArrowerSoldierData message.FreeSoldierData
	//DaoDunSoldierData message.FreeSoldierData
	//SpearSoldierData message.FreeSoldierData
	//FashiSoldierData message.FreeSoldierData
	sync.RWMutex
}
type BattleInfo struct {
	RoomId 		int32
	BattleId	int32
	IsInBattling	bool
	IsInRooming		bool
	IsLoadFinished  bool
}
func (role *OnlineRole)GetRole()*bean.Role{
	return &role.Role
}
func (role *OnlineRole)SetGateSession(session network.SocketSessionInterface){
	role.GateSession = session
}
func (role *OnlineRole)GetGateSession() network.SocketSessionInterface{
	return role.GateSession
}
func (role *OnlineRole)GetRoleId() int64{
	return role.Role.RoleId
}
func (role *OnlineRole)GetServerId() int32{
	return role.Role.ServerId
}
func (role *OnlineRole)GetUserId() int64{
	return role.Role.UserId
}
func (role *OnlineRole)GetUseName() string{
	return role.UserName
}
func (role *OnlineRole)GetGateId() int32{
	return role.GateId
}
func (role *OnlineRole)SetRoleId(roleId int64) {
	role.Role.RoleId = roleId
}
func (role *OnlineRole)SetServerId(serverId int32) {
	role.Role.ServerId = serverId
}
func (role *OnlineRole)SetUserId(userId int64){
	role.Role.UserId = userId
}
func (role *OnlineRole)SetUseName(name string) {
	role.UserName = name
}
func (role *OnlineRole)SetGateId(gateId int32) {
	role.GateId = gateId
}
func (role *OnlineRole)GetNickName()string{
	return role.Role.NickName
}
func (role *OnlineRole)SetNickName(nickName string){
	role.Role.NickName = nickName
}
func (role *OnlineRole)GetRoomId()int32{
	return role.BattleInfo.RoomId
}
func (role *OnlineRole)SetRoomId(roomId int32){
	role.BattleInfo.RoomId = roomId
}
func (role *OnlineRole)GetLevel() int32{
	return role.Role.Level
}
func (role *OnlineRole)SetLevel(level int32){
	role.Role.Level = level
}
func (role *OnlineRole)GetDiam() int32{
	return role.Role.Diam
}
func (role *OnlineRole)SetDiam(diam int32){
	role.Role.Diam = diam
}
func (role *OnlineRole)GetGold() int32{
	return role.Role.Gold
}
func (role *OnlineRole)SetGold(gold int32){
	role.Role.Gold = gold
}
func (role *OnlineRole)GetExp() int32{
	return role.Role.Exp
}
func (role *OnlineRole)SetExp(exp int32){
	role.Role.Exp = exp
}
func (role *OnlineRole)GetAvatarId() int32{
	return role.Role.AvatarId
}
func (role *OnlineRole)SetAvatarId(avatarId int32){
	role.Role.AvatarId = avatarId
}
func (role *OnlineRole)GetAllHero() map[int32]bean.Hero{
	return role.Heros
}
//是否掉线或者连接
func (role *OnlineRole)IsConnected()bool{
	role.RLock()
	defer role.RUnlock()
	return role.Connected
}
func (role *OnlineRole)SetConnected(conn bool){
	role.Lock()
	defer role.Unlock()
	role.Connected = conn
}
func (role *OnlineRole)IsInBattling()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInBattling
}
func (role *OnlineRole)SetInBattling(value bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInBattling = value
}
func (role *OnlineRole)IsInRooming()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInRooming
}
func (role *OnlineRole)SetInRooming(value bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInRooming = value
}
func (role *OnlineRole)IsLoadFinished()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsLoadFinished
}
func (role *OnlineRole)SetLoadFinished(finished bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsLoadFinished = finished
}
func (role *OnlineRole)GetBattleId() int32{
	return role.BattleInfo.BattleId
}
func (role *OnlineRole)SetBattleId(id int32){
	role.BattleInfo.BattleId = id
}
func (role *OnlineRole)GetMaxBagNum() int32{
	return role.Role.MaxBagNum
}
func (role *OnlineRole)SetMaxBagNum(num int32){
	role.Role.MaxBagNum = num
}
func (role *OnlineRole)GetPingTime()time.Time{
return  role.PingTime
}
func (role *OnlineRole)SetPingTime(time time.Time){
	role.PingTime = time
}
func (role *OnlineRole)AddGold(gold int32){
	role.Role.Gold += gold
}
func (role *OnlineRole)AddExp(exp int32){
	role.Role.Exp += exp
	if role.Role.Exp >= role.NextLevelExp{
		role.Role.Level++
		role.UpdateNextExp()
	}
}
func (role *OnlineRole)UpdateNextExp(){
	role.NextLevelExp = 75 * role.Role.Level * role.Role.Level - 25 * role.Role.Level
}
func (role *OnlineRole)AddHeroExp(heroId int32,exp int32) (int32,bool) {
	if key,ok := role.Heros[heroId];ok{
		return key.AddExp(exp),ok
	}else{
		log4g.Infof("Role[%d]不包含该英雄[%d]",role.GetRoleId(),heroId)
		return -1,false
	}
}
//func (role *OnlineRole)HasHero(heroId int32) bool{
//	if _,ok := role.Heros[heroId];ok{
//		return true
//	}else{
//		return false
//	}
//}
func (role *OnlineRole)GetItem(index int32) *bean.Item {
	if index >= role.GetMaxBagNum() {
		return nil
	}else{
		return &role.Role.Items[index]
	}
}
func (role *OnlineRole)GetHero(heroId int32)*bean.Hero{
	if hero,ok := role.Heros[heroId];ok{
		return &hero
	}else{
		return nil
	}
}
func (role *OnlineRole)UpdateDB(manager *db.MongoBDManager)  {
	dbSession := manager.Get()
	if dbSession != nil {
		//更新角色
		c := dbSession.DB("sanguozhizhan").C("Role")
		err := c.Update(bson.M{"roleid": role.Role.RoleId}, role.Role)
		if err != nil {
			log4g.Errorf("更新Role出错[%s],RoleId:%d", err.Error(), role.Role.RoleId)
			return
		}
	}
}
func (role *OnlineRole)QuitBattle(){
	role.BattleInfo.RoomId = 0
	role.BattleInfo.BattleId = 0
	role.SetLoadFinished(false)
	role.SetInBattling(false)
	role.SetInRooming(false)
}
func (role *OnlineRole)HasHero(heroId int32) bool{
	hero := role.Heros[heroId]
	if &hero == nil{
		return  false
	}
	return true
}
func (role *OnlineRole)BuyHero(heroId int32) bool{
	hasHero := role.HasHero(heroId)
	if hasHero == false{
		return false
	}
	if role.GetGold() < 8000{
		return false
	}
	role.AddGold(-8000)
	hero := bean.Hero{}
	hero.RoleId = role.Role.RoleId
	hero.HeroId = heroId
	hero.Level = 1
	hero.Exp = 0
	role.Heros[heroId] = hero
	return true
}
func (role *OnlineRole)WinLevel(level int32){
	if level == 0{
		return
	}
	has := false
	for _,v := range role.Role.WinLevel{
		if level == v{
			has = true
			break
		}
	}
	if has == false{
		role.Role.WinLevel = append(role.Role.WinLevel, level)
	}
}
func (role *OnlineRole)AddGetTaskAward(taskId int32) bool{
	for k,v := range role.Role.DayGetTask{
		if v == 0{
			role.Role.DayGetTask[k] = taskId
			return true
		}else if v == taskId{
			return false
		}
	}
	role.Role.DayGetTask = append(role.Role.DayGetTask, taskId)
	return true
}
func (role *OnlineRole)AddGetAchieveAward(achieveId int32) bool{
	for _,v := range role.Role.Achievement{
		if v == achieveId{
			return false
		}
	}
	role.Role.Achievement = append(role.Role.Achievement, achieveId)
	return true
}
//获取签到奖励
func (role *OnlineRole)GetSignAward()bool{
	if role.Role.GetSign == false{
		role.AddGold(100)
		role.Role.GetSign = true
		return true
	}else{
		return false
	}
}
func (role *OnlineRole)GetTaskSeed()int32{
	return role.Role.TaskSeed
}
func (role *OnlineRole)GetSoldierData(index int) *message.FreeSoldierData{
	data := &message.FreeSoldierData{}
	arrowerData := role.Role.FreeSoldierData[index]
	data.CarrierType = arrowerData.CarrierType
	data.PlayerType = arrowerData.PlayerType
	data.BodyId = arrowerData.BodyId
	data.TouKuiId = arrowerData.TouKuiId
	data.WeapId = arrowerData.TouKuiId
	return data
}
func (role *OnlineRole)ChangeFreeSoldierData(index int,data *message.FreeSoldierData)bool{
	if index > 3{
		return false
	}
	roleData := role.Role.FreeSoldierData[index]
	roleData.CarrierType = data.CarrierType
	roleData.WeapId = data.WeapId
	roleData.TouKuiId = data.TouKuiId
	roleData.BodyId = data.BodyId
	return true
}
func (role *OnlineRole)ChangeFreeSoldierEquipId(index int, equipIndex int, equipId int32){
	if index > 3{
		return
	}
	roleData := role.Role.FreeSoldierData[index]
	if equipIndex == 0{
		roleData.BodyId = equipId
	}else if equipIndex == 1{
		roleData.WeapId = equipId
	}else if equipIndex == 5{
		roleData.TouKuiId = equipId
	}
}
func (role *OnlineRole)GetFreeSoldierEquipId(index int, equipIndex int)int32{
	if index > 3{
		return -1
	}
	roleData := role.Role.FreeSoldierData[index]
	if equipIndex == 0{
		return roleData.BodyId
	}else if equipIndex == 1{
		return roleData.WeapId
	}else if equipIndex == 5{
		return roleData.TouKuiId
	}
}