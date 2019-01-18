package roleManager

import (
	"sync"
	"cmgameserver/bean"
	"gopkg.in/mgo.v2/bson"
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/message"
	"time"
)
type RoleManager struct {
	lock        sync.RWMutex
	onlineRoles map[int64]face.IOnlineRole//[roleId]role
	onlineSessionRoles map[network.SocketSessionInterface]map[int64]face.IOnlineRole//[session][roleId]Role
	GameServer  face.IGameServer
}

func NewRoleManager(server face.IGameServer)  *RoleManager{
	return &RoleManager{
		onlineRoles: make(map[int64]face.IOnlineRole),
		onlineSessionRoles:make(map[network.SocketSessionInterface]map[int64]face.IOnlineRole),
		GameServer:  server,
	}
}
func (manager *RoleManager)GetOnlineRole(roleId int64) face.IOnlineRole{
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.onlineRoles[roleId]
}
func (manager *RoleManager)AddOnlineRole(role face.IOnlineRole){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.onlineRoles[role.GetRoleId()] = role
	session := role.GetGateSession()
	roleMap := manager.onlineSessionRoles[session]
	if roleMap == nil{
		roleMap = make(map[int64]face.IOnlineRole)
		manager.onlineSessionRoles[session] = roleMap
	}
	roleMap[role.GetRoleId()] = role
}
func (manager *RoleManager)RemoveOnlineRole(role face.IOnlineRole){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	delete(manager.onlineRoles,role.GetRoleId())
	delete(manager.onlineSessionRoles[role.GetGateSession()],role.GetRoleId())
}
func (manager *RoleManager)NewOnlineRole(roleId int64) face.IOnlineRole{
	var err error
	if manager.GameServer.GetDBManager() != nil{
		dbSession := manager.GameServer.GetDBManager().Get()
		if dbSession != nil{
			role := bean.Role{}
			var heros []bean.Hero
			c := dbSession.DB("sanguozhizhan").C("Role")
			err = c.Find(bson.M{"roleid":roleId}).One(&role)
			if err != nil{
				return nil
			}
			c = dbSession.DB("sanguozhizhan").C("Hero")
			err = c.Find(bson.M{"roleid":roleId}).All(&heros)
			if err != nil{
				return nil
			}
			if len(role.Items) == 0{
				//说明游戏还没有背包
				for i:=0;i<int(role.MaxBagNum);i++{
					item := bean.Item{}
					item.ItemId = 0
					role.Items = append(role.Items, item)
				}
			}
			if len(heros) == 0{
				//说明游戏还没有英雄，免费送，刚开始
				hero := bean.Hero{}
				hero.RoleId = roleId
				hero.HeroId = 1
				hero.Level = 1
				hero.ItemIds[2] = 15000
				err = c.Insert(&hero)
				if err != nil{
					return nil
				}
				heros = append(heros, hero)
			}
			now := time.Now()
			if role.LoginTime != now{
				//清空每日任务
				for k,_ := range role.DayGetTask{
					role.DayGetTask[k] = 0
				}
				role.LoginTime = now
				role.TaskSeed = int32(now.Nanosecond())
				role.GetSign = false
			}
			onlineRole := OnlineRole{
				Role:role,
				Heros:make(map[int32]bean.Hero),
				BattleInfo:BattleInfo{},
				Connected:true,
				PingTime:time.Now(),
			}
			if len(heros) > 0{
				for _,v := range heros{
					onlineRole.Heros[v.HeroId] = v
				}
			}
			return &onlineRole
		}
	}
	return nil
}
func (manager *RoleManager)GetAllOnlineRole(gateSession network.SocketSessionInterface) map[int64]face.IOnlineRole{
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.onlineSessionRoles[gateSession]
}
func (manager *RoleManager)RoleQuit(roleId int64){
	defer func() {
		if err := recover();err != nil{
			log4g.Info(err.(error).Error())
			return
		}
	}()
	//如果role正在战斗中，判断战斗的进程，比如房间中，或者是正在战斗
	role := manager.GetOnlineRole(roleId)
	if role != nil{
		role.SetConnected(false)
		//更新数据库
		role.UpdateDB(manager.GameServer.GetDBManager())
		if role.IsInBattling(){
			//退出战斗
			role.SetLoadFinished(false)
		}else{
			if role.IsInRooming(){
				//退出该房间，发送消息
				roomManger := manager.GameServer.GetRoomManager()
				roomId := role.GetRoomId()
				room := roomManger.GetRoomByRoomId(roomId)
				if room != nil{
					if room.IsRoomOwner(roleId){
						//如果是房主，删除房间，并且通知所有人
						if roomManger.DeleteRoom(roomId){
							log4g.Infof("删除房间[%d]成功!",roomId)
						}else{
							log4g.Infof("删除房间[%d]失败!",roomId)
						}
					}else{
						//移除成员，并且通知房间所有人
						if roomManger.RemoveOneMemberByRoom(room,role) == false{
							log4g.Infof("移除房间[%d]成员[%d]失败!",roomId,roleId)
						}
					}
				}
			}
			//移除缓存，如果是战斗中，就不移除，等过游戏结束如果玩家还么有连接上来再移除
			manager.RemoveOnlineRole(role)
		}
		//发送给网关移除角色
		message := &message.M2G_RoleQuitGate{
			RoleId: roleId,
		}
		manager.GameServer.WriteInnerMsg(role.GetGateSession(),roleId,10005,message)
	}else{
		log4g.Infof("不存在玩家[%d],退出游戏服!",roleId)
	}
}
