package face

type IMatchTeam interface {
	GetMatchTeamId()int32
	RemoveMatchPlayer(player IMatchPlayer,server IGameServer)
	Clear(server IGameServer)
	ClearNoSend()
	GetTeamOwnerScore()int32
	GetPlayerSize()int
	GetAllMatchPlayer()[]IMatchPlayer
	GetMatchRoomIndex()int
	SetMatchRoomIndex(index int)
	GetBInMatching()bool
	SetBInMatchingWithLock(value bool)
	SetBInMatchingNoLock(value bool)
}
