package face

type IMatchRoom interface {
	AddOneMatchTeam(team IMatchTeam)
	RemoveOneMatchTeam(team IMatchTeam)
}
