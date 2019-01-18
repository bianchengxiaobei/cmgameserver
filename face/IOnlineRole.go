package face

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/bean"
	"github.com/bianchengxiaobei/cmgo/db"
	"time"
	"cmgameserver/message"
)

type IOnlineRole interface {
	GetRole()*bean.Role
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
	SetGateId(int32)
	GetNickName() string
	SetNickName(nickName string)
	GetRoomId() int32
	SetRoomId(roomId int32)
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
	GetPingTime()time.Time
	SetPingTime(time time.Time)
	IsInBattling() bool
	SetInBattling(value bool)
	IsInRooming() bool
	SetInRooming(value bool)
	GetBattleId() int32
	SetBattleId(id int32)
	GetMaxBagNum() int32
	SetMaxBagNum(num int32)
	AddGold(gold int32)
	AddExp(exp int32)
	UpdateNextExp()
	AddHeroExp(heroId int32, exp int32) (int32,bool)
	//HasHero(heroId int32) bool
	GetItem(index int32) *bean.Item
	GetHero(heroId int32)*bean.Hero
	UpdateDB(manager *db.MongoBDManager)
	QuitBattle()
	BuyHero(heroId int32) bool
	WinLevel(level int32)
	AddGetTaskAward(taskId int32) bool
	AddGetAchieveAward(achieveId int32) bool
	GetSignAward()bool
	GetTaskSeed()int32
	GetSoldierData(index int) *message.FreeSoldierData
	ChangeFreeSoldierData(index int,data *message.FreeSoldierData)bool
	ChangeFreeSoldierEquipId(index int, equipIndex int, equipId int32)
	GetFreeSoldierEquipId(index int, equipIndex int)int32
}