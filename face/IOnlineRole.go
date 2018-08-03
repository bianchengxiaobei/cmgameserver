package face

import (
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/bean"
)

type IOnlineRole interface {
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
	GetNickName()string
	SetNickName(nickName string)
	GetRoomId()int32
	SetRoomId(roomId int32)
	GetLevel()int32
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
	IsConnected()bool
	SetConnected(conn bool)
	IsLoadFinished()bool
	SetLoadFinished(finished bool)
	GetBattleId() int32
	SetBattleId(id int32)
}