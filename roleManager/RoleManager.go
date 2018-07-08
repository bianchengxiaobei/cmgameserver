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
			c := dbSession.DB("sanguozhizhan").C("Role")
			err = c.Find(bson.M{"roleid":roleId}).One(&role)
			if err != nil{
				return nil
			}
			onlineRole := OnlineRole{
				Role:role,
				BattleInfo:BattleInfo{},
				Connected:true,
			}
			return &onlineRole
		}
	}
	return nil
}
