package roleManager

import (
	"sync"
	"cmgameserver/bean"
	"github.com/bianchengxiaobei/cmgo/db"
	"gopkg.in/mgo.v2/bson"
)
type OnlineRole struct {
	Role bean.Role
	GateId	int32
	UserName string
}
type RoleManager struct {
	lock 		sync.RWMutex
	onlineRoles	map[int64]*OnlineRole
	dbManager	*db.MongoBDManager
}

func NewRoleManager(db *db.MongoBDManager)  *RoleManager{
	return &RoleManager{
		onlineRoles:make(map[int64]*OnlineRole),
		dbManager:db,
	}
}
func (manager *RoleManager)GetOnlineRole(roleId int64) *OnlineRole{
	manager.lock.RLock()
	defer manager.lock.RUnlock()
	return manager.onlineRoles[roleId]
}
func (manager *RoleManager)AddOnlineRole(role *OnlineRole){
	manager.lock.Lock()
	defer manager.lock.Unlock()
	manager.onlineRoles[role.Role.RoleId] = role
}
func (manager *RoleManager)NewOnlineRole(roleId int64) *OnlineRole{
	var err error
	if manager.dbManager != nil{
		dbSession := manager.dbManager.Get()
		if dbSession != nil{
			role := bean.Role{}
			c := dbSession.DB("sanguozhizhan").C("Role")
			err = c.Find(bson.M{"roleid":roleId}).One(&role)
			if err != nil{
				return nil
			}
			onlineRole := OnlineRole{
				Role:role,
			}
			return &onlineRole
		}
	}
	return nil
}
