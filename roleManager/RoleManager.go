package roleManager

import (
	"sync"
	"cmgameserver/bean"
	"gopkg.in/mgo.v2/bson"
	"cmgameserver/face"
)
type RoleManager struct {
	lock        sync.RWMutex
	onlineRoles map[int64]face.IOnlineRole
	GameServer  face.IGameServer
}

func NewRoleManager(server face.IGameServer)  *RoleManager{
	return &RoleManager{
		onlineRoles: make(map[int64]face.IOnlineRole),
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
			if len(heros) == 0{
				//说明游戏还没有英雄，免费送，刚开始
				hero := bean.Hero{}
				hero.RoleId = roleId
				hero.HeroId = 1
				hero.Level = 1
				err = c.Insert(&hero)
				if err != nil{
					return nil
				}
				heros = append(heros, hero)
			}
			onlineRole := OnlineRole{
				Role:role,
				Heros:make(map[int32]bean.Hero),
				BattleInfo:BattleInfo{},
				Connected:true,
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
func (manager *RoleManager)GetAllOnlineRole()map[int64]face.IOnlineRole{
	return manager.onlineRoles
}
