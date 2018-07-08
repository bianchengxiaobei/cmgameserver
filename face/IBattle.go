package face

type BattleState int
const (
	Free BattleState = iota
	InBattling
	GameOver
)
type IBattle interface {
	Start()
	GetBattleState() BattleState
	SetBattleState(state BattleState)
	AddFrameCommand(playerId int32, cmdType int32, param string)
}