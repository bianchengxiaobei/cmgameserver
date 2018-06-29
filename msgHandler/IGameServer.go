package msgHandler

import (
	"github.com/bianchengxiaobei/cmgo/db"
	"cmgameserver/roleManager"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/golang/protobuf/proto"
)

type IGameServer interface {
	GetId() int32
	GetDBManager() *db.MongoBDManager
	GetRoleManager() *roleManager.RoleManager
	WriteInnerMsg(session network.SocketSessionInterface,roleId int64,msgId int,msg proto.Message)
}
