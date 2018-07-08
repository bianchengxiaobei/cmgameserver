package roomManager

import (
	"cmgameserver/face"
	"sync"
	"container/list"
)

type RoomManager struct {
	sync.RWMutex
	rooms map[int32]face.IRoom
	recycleRoomIds	*list.List
	GameServer  face.IGameServer
} 
var genRoomId = int32(0)

func NewRoomManager(server face.IGameServer) *RoomManager{
	return &RoomManager{
		rooms:make(map[int32]face.IRoom),
		recycleRoomIds:list.New(),
		GameServer:server,
	}
}
//创建房间
func (roomManager *RoomManager)CreateRoom(roleId int64) face.IRoom{
	var roomId int32
	roomManager.Lock()
	defer roomManager.Unlock()
	if roomManager.recycleRoomIds.Len() > 0{
		roomId = roomManager.recycleRoomIds.Back().Value.(int32)
	}else {
		genRoomId++
		roomId = genRoomId
	}
	role := roomManager.GameServer.GetRoleManager().GetOnlineRole(roleId)
	if role != nil{
		role.SetRoomId(roomId)
	}
	room := &Room{
		RoomOwnerId:roleId,
		RoomId:roomId,
		GroupIdPool:[6]int32{0,1,2,3,4,5},
		RoomMembers:make(map[int64]*RoomMember),
	}
	room.RoomOwnerGroupId = room.getRandomId()
	room.roomRoleIds[room.roomIndex] = roleId
	room.roomIndex++
	roomManager.rooms[room.RoomId] = room
	return room
}
//删除房间
func (roomManager *RoomManager)DeleteRoom(roomId int32) bool{
	roomManager.Lock()
	defer roomManager.Unlock()
	room := roomManager.rooms[roomId]
	if room == nil{
		return false
	}else{
		//遍历房间内所有玩家，删除战斗信息
		roleIds := room.GetRoomRoleIds()
		for _,v:= range roleIds{
			if v > 0{
				role := roomManager.GameServer.GetRoleManager().GetOnlineRole(v)
				if role != nil{
					role.SetRoomId(0)
				}
			}
		}
		delete(roomManager.rooms, roomId)
		roomManager.recycleRoomIds.PushBack(roomId)
	}
	return true
}
func (roomManager *RoomManager)GetRoomByRoomId(roomId int32)face.IRoom{
	roomManager.RLock()
	defer  roomManager.RUnlock()
	return roomManager.rooms[roomId]
}
func (roomManager *RoomManager)GetAllRoom()map[int32]face.IRoom{
	return roomManager.rooms
}
//检查是否房间内所有玩家加载完成
func (roomManager *RoomManager)CheckAllRoomMemberLoadFinished(roomId int32) bool{
	room := roomManager.rooms[roomId]
	roleIds := room.GetRoomRoleIds()
	for _,v:= range roleIds{
		if v > 0{
			role := roomManager.GameServer.GetRoleManager().GetOnlineRole(v)
			if role != nil{
				if role.IsLoadFinished() == false && role.IsConnected(){
					return false
				}
			}
		}
	}
	return true
}