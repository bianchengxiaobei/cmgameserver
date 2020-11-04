package face

import "cmgameserver/message"

type BattleState int
const (
	Free BattleState = iota
	InBattling
	Pause
	GameOver
)
type IBattle interface {
	Start()
	GetBattleId()int32
	SetPauseRoleId(roleId int64)
	GetPauseRoleId() int64
	ReStart(roleIds *[4]IOnlineRole,battleType BattleType)
	GetBattleState() BattleState
	SetBattleState(state BattleState)
	AddFrameCommand(playerId int32, cmdType int32, param string)
	RemovePlayer(roleId int64)bool
	Finish()
	GetSaveFrames() map[int32][]*message.Command
	GetFrameCount()int32
	PauseBattle()
	RestartFormPause(msg *message.C2M2C_ReStartPauseBattle)
	//SurrenderPlayer(roleId int64,groupId int32)
	GetBattleMember() [4]IOnlineRole
	SetAgreePause(bAgree bool,roleId int64) int32
	CheckAllLeave()
	GetBattleType()BattleType
}