package face

import (
	"cmgameserver/bean"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/network"
	"time"
)


type IOnlineRole interface {
	GetRole() *bean.Role
	SetGateSession(session network.SocketSessionInterface)
	GetGateSession() network.SocketSessionInterface
	GetRoleId() int64
	GetServerId() int32
	GetUserId() int64
	GetUseName() string
	GetGateId() int32
	SetRoleId(int64)
	SetServerId(int32)
	SetUserId(int64)
	SetUseName(string)
	GetSaiJiId()int32
	SetSaiJiId(id int32)
	SetGateId(int32)
	GetNickName() string
	SetNickName(nickName string)
	GetRoomId() int32
	SetRoomId(roomId int32)
	SetBattleSeed(seed int32)
	GetBattleSeed()int32
	GetLevel() int32
	SetLevel(level int32)
	GetDiam() int32
	SetDiam(diam int32)
	GetGold() int32
	SetGold(gold int32)
	GetExp() int32
	SetExp(exp int32)
	GetAvatarId() int32
	SetAvatarId(avatarId int32)
	GetAllHero() map[int32]bean.Hero
	IsConnected() bool
	SetConnected(conn bool)
	IsLoadFinished() bool
	SetLoadFinished(finished bool)
	GetPingTime() time.Time
	SetPingTime(time time.Time)
	IsInBattling() bool
	SetInBattling(value bool)
	IsInRooming() bool
	SetInRooming(value bool)
	GetLoginTime()time.Time
	GetOnlineDiam() int32
	SetOnlineDiam(value int32)
	SetLoginTime(time time.Time)
	GetBattleId() int32
	SetBattleId(id int32)
	GetAgreePause() bool
	SetAgreePause(agree bool)
	GetSex() int32
	SetSex(sex int32)
	SetSign(sign string)
	GetSign() string
	GetMaxBagNum() int
	GetMaxEmailNum() int
	GetCurMaxEmailNum()int
	DeleteEmailReally()
	SetMaxBagNum(num int)
	SetMaxEmailNum(num int)
	AddGold(gold int32)
	AddExp(exp int32)(bool,int32)
	AddDiam(diam int32)
	GetMailAddresss() string
	SetMailAddresss(address string)
	GetQQ() int32
	SetQQ(qq int32)
	GetWeiXin() string
	SetWeiXin(weixin string)
	GetPhone() int32
	SetPhone(phone int32)
	UpdateNextExp()
	UpgradeHeroLevel(heroId int32) (level int32,skillPoint int32,value bool)
	//HasHero(heroId int32) bool
	GetItem(index int32) bean.Item
	SetItem(index int32, item bean.Item)
	DeleteItem(index int32)
	AddItem(item message.Item, hasNum bool) int32
	AddItemNoMsg(itemId int32, itemSeed int32, itemTime int64, hasNum bool) int32

	GetHero(heroId int32) *bean.Hero
	SetHero(hero bean.Hero)
	AddHero(hero bean.Hero)
	SetHeroCount(count int32)
	UpdateDB(manager *db.MongoBDManager)
	QuitBattle()
	BuyHero(heroId int32,buyType int32) bool
	WinLevel(level int32)bool
	AddGetTaskAward(taskId int32) bool
	AddGetAchieveAward(achieveId int32) bool
	GetSignAward() bool
	GetTaskSeed() int32
	GetSoldierData(index int) *message.FreeSoldierData
	GetDunJiaSoldierData()*[]*message.FreeSoldierData
	GetBuJiaSoldierData()*[]*message.FreeSoldierData
	ChangeFreeSoldierData(index int, data message.FreeSoldierData) bool
	ChangeFreeSoldierEquipId(index int, equipIndex int, equipId int32,carryType int32)
	GetFreeSoldierEquipId(index int, equipIndex int) int32
	GetFreeSoldierGuangFanEquipId(index int, equipIndex int) (int32, bool)
	//GetEmail(emailId int32) (*bean.Email, int32)
	GetEmailByIndex(index int)bean.Email
	SetEmail(emailIndex int,get bool)bool
	DeleteEmail(index int)bool
	AddEmail(email bean.Email)int32
	AddEmptyEmail(index int)
	GetMatchPlayer()IMatchPlayer
	GetRankScore()int32
	AddRankScore(value int32)int32
	GetAchieveMsgData()*message.AchievementData
	SetAchieveRecord(data *message.AchievementData)
	BuyCard(cardType bean.CardType)bool//购买月卡
	BuyDiam(diamType bean.DiamType)//购买钻石
	GetDayTaskIdList() []int32
	HasGetAchievementId()[]int32
	GetCardIds()*[]int32
	GetCards()*[]bean.Card
	SetCard(card bean.Card,index int)
	AddCardDayGetDiam(add int32)
	GetCardDayGetDiam()int32
	SetCardDayGetDiam(value int32)
	DeleteCard(index int)
	GetIsGM()bool//是否是游戏管理员
	AddTran(tran bean.Tran)
	RemoveTran(tranId string)(bool,bean.TranType,int32)
	DoNoCompleteTran()
	GetGuideFinised()bool
	SetGuideFinised(value bool)
	GetGiftCode(code string)bool
	GetInSimulateBattle()bool
	SetInSimulateBattle(bValue bool)
}
