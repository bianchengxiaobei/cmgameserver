package battleManager

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"time"
	"sync"
)

type Battle struct {
	BattleId     int32
	done         chan bool
	timer        *time.Ticker
	state        face.BattleState
	FrameCount   int32
	CommandIndex int
	msg          *message.M2C_BattleFrame
	players      [4]face.IOnlineRole
	GameServer   face.IGameServer
	CacheCommand [4]message.Command
	sync.RWMutex
}

func (battle *Battle) Start() {
	go func() {
	Loop:
		for {
			select {
			case ok := <-battle.done:
				//进入结算，然后进入空闲状态
				if ok == false {
					break Loop
				}
			case <-battle.timer.C:
				battle.BattleLoop()
			}
		}
		//结算
		battle.CalculateAward()
		battle.Reset()
		ok := <-battle.done
		if ok {
			goto Loop
		}
	}()
}
func (battle *Battle) Finish() {
	battle.state = face.GameOver
	battle.done <- false
}
func (battle *Battle) ReStart() {
	battle.done <- true
}
func (battle *Battle) GetBattleState() face.BattleState {
	return battle.state
}
func (battle *Battle) SetBattleState(state face.BattleState) {
	battle.state = state
}
func (battle *Battle) Reset() {
	battle.FrameCount = 0
	battle.state = face.Free
}
func (battle *Battle) CalculateAward() {

}
func (battle *Battle) AddFrameCommand(playerId int32, cmdType int32, param string) {
	battle.Lock()
	defer battle.Unlock()
	var cmd *message.Command
	//如果已经存在，直接修改
	cmd = &battle.CacheCommand[battle.CommandIndex]
	cmd.PlayerId = playerId
	cmd.CommandType = cmdType
	cmd.Param = param
	battle.CommandIndex++
	battle.msg.Cmd = append(battle.msg.Cmd, cmd)
}
func (battle *Battle) BattleLoop() {
	//发送当前帧给房间内所有玩家
	battle.msg.FrameCount = battle.FrameCount
	for _, v := range battle.players {
		if v != nil {
			battle.GameServer.WriteInnerMsg(v.GetGateSession(), v.GetRoleId(), 5012, battle.msg)
		}
	}
	battle.Lock()
	battle.msg.Cmd = battle.msg.Cmd[battle.CommandIndex:]
	battle.CommandIndex = 0
	battle.Unlock()
	battle.FrameCount++
}
