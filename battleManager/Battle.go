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
	msg          *message.M2C_BattleFrame
	players      [4]face.IOnlineRole
	GameServer   face.IGameServer
	frameCache	map[int32][]*message.Command
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
	battle.RLock()
	defer battle.RUnlock()
	var cmd *message.Command
	//如果已经存在，直接修改
	cmd = new(message.Command)
	cmd.PlayerId = playerId
	cmd.CommandType = cmdType
	cmd.Param = param
	battle.frameCache[battle.FrameCount] = append(battle.frameCache[battle.FrameCount],cmd)
	//log4g.Infof("f1111:%d",battle.FrameCount)
}
func (battle *Battle) BattleLoop() {
	var cmdLen int
	cmdLen = len(battle.msg.Cmd)
	if cmdLen > 0{
		battle.msg.Cmd = battle.msg.Cmd[cmdLen:]
	}
	//发送当前帧给房间内所有玩家
	battle.msg.FrameCount = battle.FrameCount
	cmd := battle.frameCache[battle.FrameCount]
	if cmd != nil && len(cmd) > 0{
		//log4g.Info("fewfee")
		battle.msg.Cmd = cmd
	}
	for _, v := range battle.players {
		if v != nil && v.IsConnected(){
			battle.GameServer.WriteInnerMsg(v.GetGateSession(), v.GetRoleId(), 5012, battle.msg)
		}
	}
	battle.Lock()
	battle.FrameCount++
	battle.Unlock()
}
