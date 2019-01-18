package face

import "cmgameserver/message"

type BattleState int
const (
	Free BattleState = iota
	InBattling
	GameOver
)
type IBattle interface {
	Start()
	ReStart(roleIds [4]IOnlineRole)
	GetBattleState() BattleState
	SetBattleState(state BattleState)
	AddFrameCommand(playerId int32, cmdType int32, param string)
	RemovePlayer(roleId int64)
	Finish()
	GetSaveFrames() map[int32][]*message.Command
	GetFrameCount()int32
}