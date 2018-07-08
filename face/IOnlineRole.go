package face

import "github.com/bianchengxiaobei/cmgo/network"

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
	IsConnected()bool
	SetConnected(conn bool)
	IsLoadFinished()bool
	SetLoadFinished(finished bool)
	GetBattleId() int32
	SetBattleId(id int32)
}