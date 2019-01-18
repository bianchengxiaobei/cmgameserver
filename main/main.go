package main

import (
	"github.com/bianchengxiaobei/cmgo/log4g"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/gameserver"
)
var server *gameserver.GameServer

func main() {
	log4g.LoadConfig("/logConfig.txt")

	server := gameserver.NewGameServer()
	server.Init("gameBaseConfig.txt", "gameSessionConfig.txt","AchieveConfig.txt","TaskConfig.txt")

	server.Run()
	network.WaitSignal()
	server.Close()
}
