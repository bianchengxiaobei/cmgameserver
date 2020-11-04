package main

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/gameserver"
	//_ "net/http/pprof"
)
var server *gameserver.GameServer

func main() {
	log4g.LoadConfig("/logConfig.txt")

	server := gameserver.NewGameServer()
	server.Init("gameBaseConfig.txt", "gameSessionConfig.txt",
		"AchieveConfig.txt","TaskConfig.txt","ServerBoxItemConfig.txt",
		"ServerHeroEquipConfig.txt","ServerSoldierEquipConfig.txt",
		"ServerEmailConfig.txt","ServerMaterialConfig.txt","ServerHeroConfig.txt",
			"ShopEquipConfig.txt","ShopServerCardConfig.txt","ServerSkillConfig.txt",
				"GiftConfig.txt","LevelUpgradeConfig.txt","ChatConfig.txt","CardInfoConfig.txt")

	server.Run()
	//http.ListenAndServe("0.0.0.0:6464",nil)
	network.WaitSignal()
	server.Close()
}
