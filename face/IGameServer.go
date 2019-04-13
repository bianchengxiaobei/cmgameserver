package face

import (
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/golang/protobuf/proto"
	"cmgameserver/bean"
)

type IGameServer interface {
	GetId() int32
	GetDBManager() *db.MongoBDManager
	GetRoleManager() IRoleManager
	GetRoomManager() IRoomManager
	GetBattleManager() IBattleManager
	WriteInnerMsg(session network.SocketSessionInterface,roleId int64,msgId int,msg proto.Message)
	GetGameVersion()string
	GetBoxItemInfoConfig() bean.BoxItemInfoConfig
	GetHeroQualityEquipConfig() bean.HeroQualityMapEquipConfig
	GetSoldierQualityEquipConfig() bean.SoldierQualityEquipConfig
	GetHeroItemIdEquipConfig() bean.HeroItemIdMapEquipConfig
	GetSoldierItemIdEquipConfig() bean.SoldierItemIdEquipConfig
	GetEmailConfig() bean.ServerEmailConfig
	GetMaterialConfig() bean.ServerMaterialConfig
	GetRankListTime() int64
	SetRankListTime(time int64)
	GetRankLevelRankList () *[30]bean.RankListItem
	GetRoleHeroCountRankList() *[30]bean.RankListItem
	GetRoleLevelRankList() *[30]bean.RankListItem
}
