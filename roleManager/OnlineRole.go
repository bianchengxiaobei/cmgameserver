package roleManager

import (
	"cmgameserver/bean"
	"github.com/bianchengxiaobei/cmgo/network"
)


type OnlineRole struct {
	Role     bean.Role
	GateId   int32
	UserName string
	BattleInfo	BattleInfo
	Connected	bool
	GateSession	network.SocketSessionInterface
}
type BattleInfo struct {
	RoomId 		int32
	BattleId	int32
	IsInBattling	bool
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
//是否掉线或者连接
func (role *OnlineRole)IsConnected()bool{
	return role.Connected
}
func (role *OnlineRole)SetConnected(conn bool){
	role.Connected = conn
}
func (role *OnlineRole)IsLoadFinished()bool{
	return role.BattleInfo.IsLoadFinished
}
func (role *OnlineRole)SetLoadFinished(finished bool){
	role.BattleInfo.IsLoadFinished = finished
}
func (role *OnlineRole)GetBattleId() int32{
	return role.BattleInfo.BattleId
}
func (role *OnlineRole)SetBattleId(id int32){
	role.BattleInfo.BattleId = id
}