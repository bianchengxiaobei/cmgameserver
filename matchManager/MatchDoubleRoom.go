package matchManager

import (
	"cmgameserver/face"
	"sync"
)

//双排房间
type MatchDoubleRoom struct {
	AllQingtongMatchTeams map[int]face.IMatchTeam //青铜段位(0,200)
	QingTongNum int
	AllBaiYinMatchTeams map[int]face.IMatchTeam //白银段位(200,500)
	BaiYinNum int
	AllHuangJingMatchTeams map[int]face.IMatchTeam //黄金段位(500,800)
	HuangJingNum int
	AllBoJingMatchTeams map[int]face.IMatchTeam //铂金段位(800,1200)
	BoJingNum int
	AllZuanShiMatchTeams map[int]face.IMatchTeam //钻石段位(1200+)
	ZuanShiNum int
	MatchManager face.IMatchManager
	sync.Mutex
}

func (this *MatchDoubleRoom)AddOneMatchTeam(team face.IMatchTeam){
	score := team.GetTeamOwnerScore()
	if score >= 0 && score < 200{
		//青铜
		this.Lock()
		this.AllQingtongMatchTeams[this.QingTongNum] = team
		this.QingTongNum++
		//如果人数满了(默认4个人,测试2个人)，就进入战斗选择画面（选择英雄携带技能，选择阵营颜色等等）
		if this.QingTongNum == OneTeamMaxPlayerNum{
			this.EnterMatchBattleRoom(this.AllQingtongMatchTeams)
			this.AllQingtongMatchTeams[0] = nil
			this.AllQingtongMatchTeams[1] = nil
			this.QingTongNum = 0
		}
		this.Unlock()
	}else if score >= 200 && score < 500{
		//白银
		this.Lock()
		this.AllBaiYinMatchTeams[this.BaiYinNum] = team
		this.BaiYinNum++
		if this.BaiYinNum == OneTeamMaxPlayerNum{
			this.EnterMatchBattleRoom(this.AllBaiYinMatchTeams)
			this.AllBaiYinMatchTeams[0] = nil
			this.AllBaiYinMatchTeams[1] = nil
			this.BaiYinNum = 0
		}
		this.Unlock()
	}else if score >= 500 && score < 800{
		//黄金
		this.Lock()
		this.AllHuangJingMatchTeams[this.HuangJingNum] = team
		this.HuangJingNum++
		if this.HuangJingNum == OneTeamMaxPlayerNum{
			this.EnterMatchBattleRoom(this.AllHuangJingMatchTeams)
			this.AllHuangJingMatchTeams[0] = nil
			this.AllHuangJingMatchTeams[1] = nil
			this.HuangJingNum = 0
		}
		this.Unlock()
	}else if score >= 800 && score < 1200{
		//铂金
		this.Lock()
		this.AllBoJingMatchTeams[this.BoJingNum] = team
		this.BoJingNum++
		if this.BoJingNum == OneTeamMaxPlayerNum{
			this.EnterMatchBattleRoom(this.AllBoJingMatchTeams)
			this.AllBoJingMatchTeams[0] = nil
			this.AllBoJingMatchTeams[1] = nil
			this.BoJingNum = 0
		}
		this.Unlock()
	}else if score >= 1200{
		//钻石
		this.Lock()
		this.AllZuanShiMatchTeams[this.ZuanShiNum] = team
		this.ZuanShiNum++
		if this.ZuanShiNum == OneTeamMaxPlayerNum{
			this.EnterMatchBattleRoom(this.AllZuanShiMatchTeams)
			this.AllZuanShiMatchTeams[0] = nil
			this.AllZuanShiMatchTeams[1] = nil
			this.ZuanShiNum = 0
		}
		this.Unlock()
	}
	team.SetBInMatchingWithLock(true)
}
func (this *MatchDoubleRoom)EnterMatchBattleRoom(all map[int]face.IMatchTeam){
	this.MatchManager.EnterMatchBattleRoom(all)
}
func (this *MatchDoubleRoom)RemoveOneMatchTeam(team face.IMatchTeam){
	score := team.GetTeamOwnerScore()
	if score >= 0 && score < 200{
		//青铜
		this.Lock()
		index := team.GetMatchRoomIndex()
		this.AllQingtongMatchTeams[index] = nil
		this.QingTongNum--
		this.Unlock()
	}else if score >= 200 && score < 500{
		//白银
		this.Lock()
		index := team.GetMatchRoomIndex()
		this.AllBaiYinMatchTeams[index] = nil
		this.BaiYinNum--
		this.Unlock()
	}else if score >= 500 && score < 800{
		//黄金
		this.Lock()
		index := team.GetMatchRoomIndex()
		this.AllHuangJingMatchTeams[index] = nil
		this.HuangJingNum--
		this.Unlock()
	}else if score >= 800 && score < 1200{
		//铂金
		this.Lock()
		index := team.GetMatchRoomIndex()
		this.AllBoJingMatchTeams[index] = nil
		this.BoJingNum--
		this.Unlock()
	}else if score >= 1200{
		//钻石
		this.Lock()
		index := team.GetMatchRoomIndex()
		this.AllZuanShiMatchTeams[index] = nil
		this.ZuanShiNum--
		this.Unlock()
	}
	team.SetBInMatchingWithLock(false)
}