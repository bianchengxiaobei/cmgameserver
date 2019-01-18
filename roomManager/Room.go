package roomManager

import (
	"sync"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"math/rand"
	"cmgameserver/message"
)

type Room struct {
	sync.RWMutex
	RoomId		int32
	MapId		int32
	RoomOwnerGroupId	int32
	RoomName 	string
	RoomOwnerId	int64
	RoomOwnerAvatarId  int32
	RoomOwnerName	string
	GameType	int32
	RoomMaxPlayerNum	int32
	Locked				bool
	RoomPassword		string
	RoomMembers			map[int64]*RoomMember
	roomRoleIds			[4]int64
	roomIndex 			int32
	GroupIdPool			[6]int32
	IsWarFow			bool
	RoomOwnerArrowerData  *message.FreeSoldierData
	RoomOwnerDaodunData  *message.FreeSoldierData
	RoomOwnerSpearData  *message.FreeSoldierData
	RoomOwnerFashiData  *message.FreeSoldierData
	Seed			int32
}
/*
public enum GroupColor
{
    蓝色,
    绿色,
    黄色,
    红色,
    紫色,
    青色
}
*/
type GroupId int32
const(
	Blue GroupId = iota
	Green
	Yellow
	Red
	Magenta
	Cyan
)
func (room *Room) GetRoomName() string {
	return room.RoomName
}

func (room *Room) SetRoomName(name string) {
	room.RoomName = name
}

func (room *Room) GetRoomId() int32 {
	return room.RoomId
}

func (room *Room) GetRoomOwnerId() int64 {
	return room.RoomOwnerId
}

func (room *Room) GetRoomMaxPlayerNum() int32 {
	return room.RoomMaxPlayerNum
}

func (room *Room) SetRoomMaxPlayerNum(num int32) {
	room.RoomMaxPlayerNum = num
}

func (room *Room) SetGameType(gameType int32) {
	room.GameType = gameType
}

func (room *Room) GetGameType() int32 {
	return room.GameType
}
func (room *Room) SetMapId(mapId int32) {
	room.MapId = mapId
}

func (room *Room) GetMapId() int32 {
	return room.MapId
}
func (room *Room)GetRoomOwnerGroupId() int32{
	return room.RoomOwnerGroupId
}
func (room *Room)SetRoomOwnerGroupId(id int32){
	room.RoomOwnerGroupId = id
}
func (room *Room)GetRoomOwnerName()string{
	return room.RoomOwnerName
}
func (room *Room)SetIsWarFow(value bool){
	room.IsWarFow = value
}
func (room *Room)GetIsWarFow() bool{
	return room.IsWarFow
}
func (room *Room)GetArrowerData() *message.FreeSoldierData{
	return room.RoomOwnerArrowerData
}
func (room *Room)GetDaodunData() *message.FreeSoldierData{
	return room.RoomOwnerDaodunData
}
func (room *Room)GetSpearData() *message.FreeSoldierData{
	return room.RoomOwnerSpearData
}
func (room *Room)GetFashiData() *message.FreeSoldierData{
	return room.RoomOwnerFashiData
}
//加入一个成员
func (room *Room) JoinOneMember(roleId int64) (int32,bool){
	room.Lock()
	defer room.Unlock()
	//if room.GetCurPlayerNum() > 4{
	//	return -1,false
	//}
	if room.RoomMembers[roleId] == nil{
		member := &RoomMember{
			RoleId:roleId,
			Prepare:false,
		}
		room.RoomMembers[roleId] = member
		room.roomRoleIds[room.roomIndex] = roleId
		room.roomIndex++
		//随机分配一个GroupId
		member.GroupId = room.getRandomId()
		return member.GroupId,true
	}
	return -1,false
}
//退出一个成员
func (room *Room) LeaveOneMember(roleId int64) bool{
	room.Lock()
	defer room.Unlock()
	member := room.RoomMembers[roleId]
	if member == nil{
		return false
	}else{
		room.addRandomId(member.GroupId)
		delete(room.RoomMembers,roleId)
		room.RemoveRoleId(roleId)
		room.roomIndex--
		return true
	}
}
//是否是房主
func (room *Room)IsRoomOwner(roleId int64)bool{
	if room.RoomOwnerId == roleId{
		return true
	}else{
		return false
	}
}
func (room *Room)GetRoomRoleIds() [4]int64{
	room.RLock()
	defer room.RUnlock()
	return room.roomRoleIds
}
func (room *Room)RemoveRoleId(roleId int64){
	for k,v := range room.roomRoleIds{
		if v == roleId{
			room.roomRoleIds[k] = 0
			return
		}
	}
	log4g.Infof("房间不存role:[%d]",roleId)
}
//取得当前房间的人数
func (room *Room)GetCurPlayerNum() int32{
	room.RLock()
	defer room.RUnlock()
	oNum := len(room.RoomMembers)
	return 1 + int32(oNum)
}
//取得随机id
func (room *Room)getRandomId() int32{
	Random:
	randomIndex := rand.Int31n(5)
	if room.GroupIdPool[randomIndex] == -1{
		goto Random
	}else{
		id := room.GroupIdPool[randomIndex]
		room.GroupIdPool[randomIndex] = -1
		return id
	}
}
//退出房间是回收小组id
func (room *Room)addRandomId(id int32){
	room.GroupIdPool[id] = id
}
//检查是否房间内玩家已经准备
func (room *Room)CheckRoomReady()bool{
	room.RLock()
	defer room.RUnlock()
	for _,v:=range room.RoomMembers{
		if v.Prepare == false{
			return false
		}
	}
	return true
}
//设置房间内玩家是否准备
func (room *Room)SetRoomMemberReady(ready bool,roleId int64) bool{
	member := room.RoomMembers[roleId]
	if member != nil{
		member.Prepare = ready
		return true
	}else{
		return false
	}
}
func (room *Room)GetSeed()int32{
	return room.Seed
}
func (room *Room)SetSeed(seed int32) {
	room.Seed = seed
}
func (room *Room)GetRoomOwnerAvatarId()int32  {
	return room.RoomOwnerAvatarId
}
func (room *Room)GetRoomMemberGroupId(roleId int64) int32{
	member := room.RoomMembers[roleId]
	if member == nil{
		log4g.Infof("member == nil,Id[%d]",roleId)
		return -1
	}
	return member.GroupId
}
func (room *Room)GetRoomMemberReady(roleId int64) bool{
	member := room.RoomMembers[roleId]
	if member == nil{
		log4g.Infof("member == nil,Id[%d]",roleId)
		return false
	}
	return member.Prepare
}