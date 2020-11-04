package matchManager

import (
	"cmgameserver/face"
	"sync"
)

type MatchPlayer struct {
	OnlineRole face.IOnlineRole//在线角色
	IsOwner   bool//是否是房主
	IsPunish  bool //是否强制退出，需要5分钟才能再次进入匹配
	PosIndex int32//所在队伍的位置坐标
	BInMatching bool//是否正在匹配（用于玩家退出游戏的时候判断）
	BInBattleRoom bool //是否在匹配房间内
	TeamId   int32//所在匹配队伍的id
	RoomId  int32//所在Room的id
	GroupId int32
	Prepare bool//是否准备好了
	sync.RWMutex//读写锁
}

func (this *MatchPlayer)SetBOwner(value bool){
	this.IsOwner = value
}
func (this *MatchPlayer)GetBOwner()bool{
	return this.IsOwner
}
func (this *MatchPlayer)GetRoleId()int64{
	return this.OnlineRole.GetRoleId()
}
func (this *MatchPlayer)GetOnlineRole()face.IOnlineRole{
	return this.OnlineRole
}

func (this *MatchPlayer)SetPosIndex(value int32){
	this.PosIndex = value
}
func (this *MatchPlayer)GetPosIndex()int32{
	return this.PosIndex
}
func (this *MatchPlayer)GetBInMatching()bool{
	this.RLock()
	defer this.RUnlock()
	return this.BInMatching
}
func (this *MatchPlayer)SetBInMatching(value bool){
	this.BInMatching = value
}
func (this *MatchPlayer)GetBInBattleRoom()bool{
	return this.BInBattleRoom
}
func (this *MatchPlayer)SetBInBattleRoom(value bool){
	this.Lock()
	defer this.Unlock()
	this.BInBattleRoom = value
}
func (this *MatchPlayer)GetMatchTeamId()int32{
	return this.TeamId
}
func (this *MatchPlayer)SetMatchTeamId(value int32){
	this.TeamId = value
}
func (this *MatchPlayer)Clear(){
	this.SetMatchTeamId(0)
	this.SetBOwner(false)
	this.SetPosIndex(0)
	this.SetBInMatching(false)
	this.SetPrepare(false)
	this.SetGroupId(0)
	this.SetMatchRoomId(0)
	this.SetBInBattleRoom(false)
	role := this.GetOnlineRole()
	role.SetInBattling(false)
	role.SetLoadFinished(false)
	role.SetBattleId(0)
	role.SetBattleSeed(0)
	role.SetRoomId(0)
}
func (this *MatchPlayer)GetPrepare()bool{
	this.RLock()
	defer this.RUnlock()
	return this.Prepare
}

func (this *MatchPlayer)SetPrepare(value bool){
	this.Lock()
	defer this.Unlock()
	this.Prepare = value
}
func (this *MatchPlayer)SetGroupId(value int32){
	this.GroupId = value
}
func (this *MatchPlayer)GetGroupId()int32{
	return this.GroupId
}
func (this *MatchPlayer)GetMatchRoomId()int32{
	return this.RoomId
}
func (this *MatchPlayer)SetMatchRoomId(value int32){
	this.RoomId = value
}