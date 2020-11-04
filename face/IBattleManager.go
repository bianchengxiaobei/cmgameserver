package face

type IBattleManager interface {
	CreateBattle(roleIds *[4]IOnlineRole,battleType BattleType) IBattle
	GetBattleInFree() IBattle
	GetBattle(battleId int32)IBattle
}