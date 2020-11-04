package matchManager

import (
	"cmgameserver/face"
	"sync/atomic"
	"sync"
	"cmgameserver/message"
)

//匹配队伍（双排，单排）
var teamId int32
type MatchTeam struct {
	MatchPlayers   map[int]face.IMatchPlayer//队伍里面的玩家
	TeamId         int32//队伍id
	PlayerSize     int
	RoomMatchIndex int
	BInMatching    bool
	sync.RWMutex
}


func CreateMatchTeam(teamOwner face.IMatchPlayer)face.IMatchTeam{
	teamOwner.SetBOwner(true)//设置第一个加入的为房主
	team := MatchTeam{
		TeamId:0,
		MatchPlayers:make(map[int]face.IMatchPlayer,0),
		RoomMatchIndex:0,
		BInMatching:false,
	}
	gTeamId := atomic.AddInt32(&teamId,1)
	team.TeamId = gTeamId
	team.AddMatchPlayer(teamOwner)
	//log4g.Infof("创建匹配队伍[%d]",team.TeamId)
	return &team
}
//这里得加锁
func (this *MatchTeam)AddMatchPlayer(player face.IMatchPlayer)bool{
	this.Lock()
	defer this.Unlock()
	//从0开始排
	for i:= 0;i<10;i++{
		temp := this.MatchPlayers[i]
		if temp == nil{
			player.SetPosIndex(int32(i))
			//this.SetBInMatchingNoLock(true)//正在匹配
			player.SetMatchTeamId(this.TeamId)
			this.MatchPlayers[i] = player
			this.PlayerSize += 1
			return true
		}
	}
	return false
}
//移除一个队员
func (this *MatchTeam)RemoveMatchPlayer(player face.IMatchPlayer,server face.IGameServer){
	this.Lock()
	defer this.Unlock()
	tempPosIndex := int(player.GetPosIndex())
	tempPlayer := this.MatchPlayers[tempPosIndex]
	if tempPlayer != nil && tempPlayer == player{
		returnMsg := new(message.M2C_RemoveMatchPlayerFromMatchTeam)
		returnMsg.RoleId = tempPlayer.GetRoleId()
		//发送给所有玩家，该玩家退出队伍的消息
		for _,v := range this.MatchPlayers{
			//发送移除该队伍消息
			role := v.GetOnlineRole()
			if role != nil{
				server.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5065,returnMsg)
			}
		}
		tempPlayer.Clear()
		delete(this.MatchPlayers, tempPosIndex)
		this.PlayerSize--
	}
}

func (this *MatchTeam)GetMatchTeamId()int32{
	return  this.TeamId
}
func (this *MatchTeam)Clear(server face.IGameServer){
	this.TeamId = 0
	if len(this.MatchPlayers) > 0{
		returnMsg := new(message.M2C_RemoveMatchTeam)
		for _,v := range this.MatchPlayers{
			//发送移除该队伍消息
			role := v.GetOnlineRole()
			if role != nil{
				server.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5064,returnMsg)
			}
			v.SetPosIndex(0)
			v.SetBOwner(false)
			v.SetMatchTeamId(0)
		}
	}
	this.SetBInMatchingNoLock(false)
	//清空
	this.MatchPlayers = make(map[int]face.IMatchPlayer,0)
}
func (this *MatchTeam)ClearNoSend(){
	this.TeamId = 0
	if len(this.MatchPlayers) > 0{
		for _,v := range this.MatchPlayers{
			v.SetBOwner(false)
			v.SetMatchTeamId(0)
		}
	}
	//清空
	this.MatchPlayers = make(map[int]face.IMatchPlayer,0)
}
///取得队伍队长的段位分数
func (this *MatchTeam)GetTeamOwnerScore()int32{
	if len(this.MatchPlayers) > 0{
		owner := this.MatchPlayers[0]
		if owner != nil{
			return owner.GetOnlineRole().GetRankScore()
		}
	}
	return 0
}
func (this *MatchTeam)GetPlayerSize()int{
	this.RLock()
	defer this.RUnlock()
	return this.PlayerSize
}
func (this *MatchTeam)GetAllMatchPlayer()[]face.IMatchPlayer{
	if this.PlayerSize > 0{
		players := make([]face.IMatchPlayer,this.PlayerSize)
		for i := 0;i<this.PlayerSize;i++{
			players = append(players, this.MatchPlayers[i])
		}
		return players
	}else{
		return nil
	}
}
func (this *MatchTeam)GetMatchRoomIndex()int{
	return this.RoomMatchIndex
}
func (this *MatchTeam)SetMatchRoomIndex(index int){
	this.RoomMatchIndex = index
}
func (this *MatchTeam)GetBInMatching()bool{
	this.RLock()
	defer this.RUnlock()
	return this.BInMatching
}
func (this *MatchTeam)SetBInMatchingWithLock(value bool){
	this.Lock()
	defer this.Unlock()
	this.BInMatching = value
	for _,v := range this.MatchPlayers{
		if v != nil{
			v.SetBInMatching(value)
		}
	}
	//println("1:",this.BInMatching)
}
func (this *MatchTeam)SetBInMatchingNoLock(value bool){
	this.BInMatching = value
	for _,v := range this.MatchPlayers{
		if v != nil{
			v.SetBInMatching(value)
		}
	}
	//println("2:",this.BInMatching)
}