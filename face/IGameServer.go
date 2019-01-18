package face

import (
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/golang/protobuf/proto"
)

type IGameServer interface {
	GetId() int32
	GetDBManager() *db.MongoBDManager
	GetRoleManager() IRoleManager
	GetRoomManager() IRoomManager
	GetBattleManager() IBattleManager
	WriteInnerMsg(session network.SocketSessionInterface,roleId int64,msgId int,msg proto.Message)
	GetGameVersion()string
}
