package roleManager

import (
	"cmgameserver/bean"
	"github.com/bianchengxiaobei/cmgo/network"
	"sync"
)


type OnlineRole struct {
	Role     bean.Role
	Heros	map[int32]bean.Hero
	GateId   int32
	UserName string
	BattleInfo	BattleInfo
	Connected	bool
	GateSession	network.SocketSessionInterface
	sync.RWMutex
}
type BattleInfo struct {
	RoomId 		int32
	BattleId	int32
	IsInBattling	bool
	IsInRooming		bool
	IsLoadFinished  bool
}
func (role *OnlineRole)SetGateSession(session network.SocketSessionInterface){
	role.GateSession = session
}
func (role *OnlineRole)GetGateSession() network.SocketSessionInterface{
	return role.GateSession
}
func (role *OnlineRole)GetRoleId() int64{
	return role.Role.RoleId
}
func (role *OnlineRole)GetServerId() int32{
	return role.Role.ServerId
}
func (role *OnlineRole)GetUserId() int64{
	return role.Role.UserId
}
func (role *OnlineRole)GetUseName() string{
	return role.UserName
}
func (role *OnlineRole)GetGateId() int32{
	return role.GateId
}
func (role *OnlineRole)SetRoleId(roleId int64) {
	role.Role.RoleId = roleId
}
func (role *OnlineRole)SetServerId(serverId int32) {
	role.Role.ServerId = serverId
}
func (role *OnlineRole)SetUserId(userId int64){
	role.Role.UserId = userId
}
func (role *OnlineRole)SetUseName(name string) {
	role.UserName = name
}
func (role *OnlineRole)SetGateId(gateId int32) {
	role.GateId = gateId
}
func (role *OnlineRole)GetNickName()string{
	return role.Role.NickName
}
func (role *OnlineRole)SetNickName(nickName string){
	role.Role.NickName = nickName
}
func (role *OnlineRole)GetRoomId()int32{
	return role.BattleInfo.RoomId
}
func (role *OnlineRole)SetRoomId(roomId int32){
	role.BattleInfo.RoomId = roomId
}
func (role *OnlineRole)GetLevel() int32{
	return role.Role.Level
}
func (role *OnlineRole)SetLevel(level int32){
	role.Role.Level = level
}
func (role *OnlineRole)GetDiam() int32{
	return role.Role.Diam
}
func (role *OnlineRole)SetDiam(diam int32){
	role.Role.Diam = diam
}
func (role *OnlineRole)GetGold() int32{
	return role.Role.Gold
}
func (role *OnlineRole)SetGold(gold int32){
	role.Role.Gold = gold
}
func (role *OnlineRole)GetExp() int32{
	return role.Role.Exp
}
func (role *OnlineRole)SetExp(exp int32){
	role.Role.Exp = exp
}
func (role *OnlineRole)GetAvatarId() int32{
	return role.Role.AvatarId
}
func (role *OnlineRole)SetAvatarId(avatarId int32){
	role.Role.AvatarId = avatarId
}
func (role *OnlineRole)GetAllHero() map[int32]bean.Hero{
	return role.Heros
}
//是否掉线或者连接
func (role *OnlineRole)IsConnected()bool{
	role.RLock()
	defer role.RUnlock()
	return role.Connected
}
func (role *OnlineRole)SetConnected(conn bool){
	role.Lock()
	defer role.Unlock()
	role.Connected = conn
}
func (role *OnlineRole)IsInBattling()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInBattling
}
func (role *OnlineRole)SetInBattling(value bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInBattling = value
}
func (role *OnlineRole)IsInRooming()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsInRooming
}
func (role *OnlineRole)SetInRooming(value bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsInRooming = value
}
func (role *OnlineRole)IsLoadFinished()bool{
	role.RLock()
	defer role.RUnlock()
	return role.BattleInfo.IsLoadFinished
}
func (role *OnlineRole)SetLoadFinished(finished bool){
	role.Lock()
	defer role.Unlock()
	role.BattleInfo.IsLoadFinished = finished
}
func (role *OnlineRole)GetBattleId() int32{
	return role.BattleInfo.BattleId
}
func (role *OnlineRole)SetBattleId(id int32){
	role.BattleInfo.BattleId = id
}