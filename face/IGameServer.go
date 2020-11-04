package face

import (
	"github.com/bianchengxiaobei/cmgo/db"
	"github.com/bianchengxiaobei/cmgo/network"
	"github.com/golang/protobuf/proto"
	"cmgameserver/bean"
	"cmgameserver/message"
)
type IGameServer interface {
	GetId() int32
	GetDBManager() *db.MongoBDManager
	GetRoleManager() IRoleManager
	GetRoomManager() IRoomManager
	GetMatchManager() IMatchManager
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
	GetHeroConfig() bean.HeroConfig
	GetGiftConfig() bean.GiftConfig
	GetAchieveConfig() bean.AchieveConfig
	GetTaskConfig() bean.TaskConfig
	GetShopEquipConfig() bean.ServerShopEquipConfig
	GetCardInfoConfig() bean.CardInfoConfig
	GetShopHeroCardConfig() bean.ShopHeroCardConfig
	GetSkillConfig() bean.ServerSkillConfig
	GetLevelUpgradeConfig() bean.LevelUpgradeConfig
	GetRankListTime() int64
	SetRankListTime(time int64)
	GetRankLevelRankList () *[30]bean.RankListItem
	GetRoleHeroCountRankList() *[30]bean.RankListItem
	GetRoleLevelRankList() *[30]bean.RankListItem
	GetChatInfo() *[10]message.ChatInfo
	AddChatInfo(avatarId int32,name string,content string,time int64,level int32,
		rankScore int32,sex int32,roomId int32)
	SetNeedClose()
	GetNeedClose()bool
	SubAllOnlineDiam(value int32)(bool,int32)
}
