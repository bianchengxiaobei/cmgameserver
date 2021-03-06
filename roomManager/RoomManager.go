package roomManager

import (
	"cmgameserver/face"
	"sync"
	"container/list"
	"cmgameserver/message"
	"time"
	"github.com/bianchengxiaobei/cmgo/log4g"
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
		role.SetInRooming(true)
	}
	room := &Room{
		RoomOwnerId:roleId,
		RoomOwnerAvatarId:role.GetAvatarId(),
		RoomOwnerName:role.GetNickName(),
		RoomOwnerLevel:role.GetLevel(),
		RoomId:roomId,
		GroupIdPool:[6]int32{0,1,2,3,4,5},
		RoomMembers:make(map[int64]*RoomMember),
	}
	room.Seed = int32(time.Now().UnixNano())
	room.RoomOwnerGroupId = room.getRandomId()
	room.roomRoleIds[room.roomIndex] = roleId
	room.roomIndex++
	room.RoomOwnerArrowerData = role.GetSoldierData(0)
	room.RoomOwnerDaodunData = role.GetSoldierData(1)
	room.RoomOwnerSpearData = role.GetSoldierData(2)
	room.RoomOwnerFashiData = role.GetSoldierData(3)
	roomManager.rooms[room.RoomId] = room
	log4g.Infof("创建游戏房间[%d]",roomId)
	return room
}
//创建默认的房间
func (roomManager *RoomManager)CreateDefaultRoom(roleId int64) face.IRoom{
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
		role.SetInRooming(true)
	}
	room := &Room{
		RoomOwnerId:roleId,
		RoomOwnerAvatarId:role.GetAvatarId(),
		RoomOwnerName:role.GetNickName(),
		RoomOwnerLevel:role.GetLevel(),
		RoomId:roomId,
		GroupIdPool:[6]int32{0,1,2,3,4,5},
		RoomMembers:make(map[int64]*RoomMember),
	}
	room.Seed = int32(time.Now().UnixNano())
	room.RoomOwnerGroupId = room.getRandomId()
	room.roomRoleIds[room.roomIndex] = roleId
	room.roomIndex++
	room.RoomOwnerArrowerData = role.GetSoldierData(0)
	room.RoomOwnerDaodunData = role.GetSoldierData(1)
	room.RoomOwnerSpearData = role.GetSoldierData(2)
	room.RoomOwnerFashiData = role.GetSoldierData(3)
	room.SetRoomMaxPlayerNum(4)
	/*
	public enum EBattleType
{
    SimulateBattle = 0,
    Battle,
    PaiWei,
    NetRoom,
    LocalRoom,
}
	 */
	room.SetGameType(3)
	room.SetIsWarFow(true)
	room.SetIsOutsideMonster(true)
	room.SetRoomOwnerCityId(1)
	room.SetMapId(0)//默认地图id为0
	roomManager.rooms[room.RoomId] = room
	log4g.Infof("创建大厅游戏房间[%d]",roomId)
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
		rMsg := &message.M2C_RoomDelete{}
		rMsg.RoomId = roomId
		for _,v:= range roleIds{
			if v > 0{
				role := roomManager.GameServer.GetRoleManager().GetOnlineRole(v)
				if role != nil{
					if role.GetRoomId() != roomId{
						continue
					}
					role.SetRoomId(0)
					role.SetInRooming(false)
					role.SetInBattling(false)
					role.SetLoadFinished(false)
					if role.IsConnected() && role.GetBattleId() == 0{
						//通知
						roomManager.GameServer.WriteInnerMsg(role.GetGateSession(),v,5015,rMsg)
					}
					if role.GetBattleId() > 0{
						role.SetBattleId(0)
					}
				}
			}
		}
		room.SetInBattle(false)
		delete(roomManager.rooms, roomId)
		roomManager.recycleRoomIds.PushBack(roomId)
		log4g.Infof("回收房间[%d]",roomId)
	}
	return true
}
func (roomManager *RoomManager)RemoveOneMemberByRoom(room face.IRoom,role face.IOnlineRole) bool{
	roomManager.Lock()
	defer roomManager.Unlock()
	if room != nil{
		if room.LeaveOneMember(role.GetRoleId()){
			role.SetInRooming(false)
			role.SetRoomId(0)
			role.SetInBattling(false)
			//通知（包括自己）
			rMsg := &message.M2C_RoleQuitRoom{}
			rMsg.RoleId = role.GetRoleId()
			rMsg.RoomId = room.GetRoomId()
			roleIds := room.GetRoomRoleIds()
			for _,v:= range roleIds{
				if v > 0{
					r := roomManager.GameServer.GetRoleManager().GetOnlineRole(v)
					if r != nil{
						roomManager.GameServer.WriteInnerMsg(r.GetGateSession(),r.GetRoleId(),5016,rMsg)
					}
				}
			}
			//（包括自己）
			roomManager.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5016,rMsg)
			return true
		}else{
			return false
		}
	}else{
		return false
	}
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
	if room != nil{
		roleIds := room.GetRoomRoleIds()
		for _,v:= range roleIds{
			if v > 0{
				role := roomManager.GameServer.GetRoleManager().GetOnlineRole(v)
				if role != nil && role.IsConnected(){
					bLoad := role.IsLoadFinished()
					if bLoad == false{
						return false
					}
				}
			}
		}
		return true
	}
	return false
}