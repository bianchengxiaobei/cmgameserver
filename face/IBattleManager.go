package face

type IBattleManager interface {
	CreateBattle(roleIds [4]IOnlineRole) IBattle
	GetBattleInFree() IBattle
	GetBattle(battleId int32)IBattle
}