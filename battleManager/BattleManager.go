package battleManager

import (
	"cmgameserver/face"
	"time"
	"sync"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type BattleManager struct {
	GameServer face.IGameServer
	battles map[int32]face.IBattle
	sync.RWMutex
}
var genBattleId = int32(0)

func NewBattleManager(server face.IGameServer) *BattleManager{
	return &BattleManager{
		battles:make(map[int32]face.IBattle),
		GameServer:server,
	}
}
func (battleManager *BattleManager) CreateBattle(roles *[4]face.IOnlineRole,battleType face.BattleType) face.IBattle{
	genBattleId++
	//设置战斗Id
	num := 0
	for _,v := range roles {
		if v != nil{
			num++
			v.SetBattleId(genBattleId)
			v.SetInBattling(true)
		}
	}
	battle := &Battle{
		BattleId:   genBattleId,
		done:       make(chan int,1),
		timer:      time.NewTicker(50 * time.Millisecond),
		Players:    *roles,
		PlayerCount:num,
		msg:        new(message.M2C_BattleFrame),
		frameCache: make(map[int32][]*message.Command),
		GameServer: battleManager.GameServer,
		BattleType:battleType,
	}
	for _,v := range roles{
		if v != nil{
			battle.RoomId = v.GetRoomId()
			break
		}
	}
	battle.SetBattleState(face.InBattling)
	battle.msg.Cmd = make([]*message.Command,0)
	battleManager.battles[battle.BattleId] = battle
	log4g.Infof("新增战斗[%d]",battle.BattleId)
	return battle
}
//取得空闲的战斗
func (battleManager *BattleManager)GetBattleInFree() face.IBattle{
	battleManager.RLock()
	defer battleManager.RUnlock()
	for _,v := range battleManager.battles{
		if v.GetBattleState() == face.Free{
			return v
		}
	}
	return nil
}
func (battleManager *BattleManager)GetBattle(battleId int32)face.IBattle{
	return battleManager.battles[battleId]
}