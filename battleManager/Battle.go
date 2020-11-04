package battleManager

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"time"
	"sync"
	"github.com/bianchengxiaobei/cmgo/log4g"
)

type Battle struct {
	BattleId   int32
	RoomId     int32
	PauseRoleId int64
	PlayerCount int//战斗人数
	done       chan int//如果是1表示重新开始，0表示结束,如果是2表示暂停,3表示暂停结束
	timer      *time.Ticker
	state      face.BattleState
	FrameCount int32
	BattleType face.BattleType
	msg        *message.M2C_BattleFrame
	Players    [4]face.IOnlineRole
	GameServer face.IGameServer
	frameCache map[int32][]*message.Command
	sync.RWMutex
}

func (this *Battle)GetBattleId()int32{
	return this.BattleId
}

func (battle *Battle) Start() {
	go func() {
	Loop:
		for {
			select {
			case ok := <-battle.done:
				//进入结算，然后进入空闲状态
				if ok == 0 {
					battle.Reset()
					break Loop
				}else if ok == 2{
					break Loop
				}
 			case <-battle.timer.C:
				battle.BattleLoop()
			}
		}
		//阻塞，等待重新利用线程
		ok := <-battle.done
		if ok == 1 {
			goto Loop
		}else if ok == 3{
			goto Loop
		}
	}()
}
func (battle *Battle) Finish() {
	log4g.Infof("战斗结束[%d]",battle.BattleId)
	rmsg := new(message.M2C_BattleFinished)
	battleState := battle.GetBattleState()
	for _,v := range battle.Players{
		if v != nil{
			if v.IsInBattling(){
				v.SetInBattling(false)
				log4g.Info("22222")
			}
			if v.GetBattleId() > 0{
				v.SetBattleId(0)
				log4g.Info("3333")
			}
			if v.IsConnected(){
				battle.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5108,rmsg)
			}
		}
	}
	if battleState == face.GameOver{
		return
	}
	if battleState == face.Pause{
		battle.Reset()
		battle.SetBattleState(face.GameOver)
	}else{
		battle.SetBattleState(face.GameOver)
		battle.done <- 0
	}
}
func (battle *Battle) ReStart(roleIds *[4]face.IOnlineRole,battleType face.BattleType) {
	battle.Players = *roleIds
	battle.PlayerCount = 0
	for _,v := range battle.Players {
		if v != nil{
			battle.PlayerCount++//人数增加1
			battle.RoomId = v.GetRoomId()
			v.SetBattleId(battle.BattleId)
			v.SetInBattling(true)
		}
	}
	log4g.Infof("重新开始战斗{%d}",battle.BattleId)
	battle.BattleType = battleType
	battle.SetBattleState(face.InBattling)
	battle.done <- 1
}
func (battle *Battle)PauseBattle(){
	battle.SetBattleState(face.Pause)
	battle.done <- 2
}
//从暂停中恢复游戏
func (battle *Battle)RestartFormPause(msg *message.C2M2C_ReStartPauseBattle){
	battle.SetBattleState(face.InBattling)
	battle.done <- 3//发送3或者1
	//发送给所有在该房间战斗的玩家
	for _, v := range battle.Players {
		if v != nil{
			v.SetAgreePause(false)
			if v.IsConnected(){
				battle.GameServer.WriteInnerMsg(v.GetGateSession(), v.GetRoleId(), 5100, msg)
			}
		}
	}
}
func (battle *Battle) GetBattleState() face.BattleState {
	battle.RLock()
	defer battle.RUnlock()
	return battle.state
}
func (battle *Battle) SetBattleState(state face.BattleState) {
	battle.Lock()
	defer battle.Unlock()
	battle.state = state
}
func (battle *Battle)GetBattleMember() [4]face.IOnlineRole{
	return battle.Players
}
func (battle *Battle) Reset() {
	battle.FrameCount = 0
	battle.SetBattleState(face.Free)
	for _,v := range battle.Players{
		if v != nil && v.IsConnected(){
			v.SetLoadFinished(false)//战斗结束设置战斗加载没完成
		}
	}
	//清空角色列表
	battle.Players[0] = nil
	battle.Players[1] = nil
	battle.Players[2] = nil
	battle.Players[3] = nil
	//清除缓存帧
	for k,_:=range battle.frameCache{
		delete(battle.frameCache,k)
	}
}
func (battle *Battle) CalculateAward() {

}
func (battle *Battle)GetBattleType()face.BattleType{
	return battle.BattleType
}
func (battle *Battle)GetPauseRoleId() int64{
	return battle.PauseRoleId
}
func (battle *Battle)SetPauseRoleId(roleId int64){
	battle.PauseRoleId = roleId
}
//设置战斗成员是否同意暂停(0-不能暂停，1-所有玩家都暂停,2-表示有玩家还没同意)
func (battle *Battle)SetAgreePause(bAgree bool,roleId int64)int32{
	if bAgree == false{
		for _,v := range battle.Players{
			if v != nil{
				v.SetAgreePause(false)
			}
		}
		return 0
	}
	for _,v := range battle.Players{
		if v != nil{
			if v.GetRoleId() == roleId{
				v.SetAgreePause(bAgree)
			}
		}
	}
	//看是否所有人都同意了，如果所有人同意，执行暂停
	for _,v := range battle.Players{
		if v != nil{
			//有一个人不同意，就跳出
			if v.GetAgreePause() == false{
				return 2
			}
		}
	}
	//开始暂停
	battle.PauseBattle()
	return 1
}
func (battle *Battle) AddFrameCommand(playerId int32, cmdType int32, param string) {
	battle.Lock()
	defer battle.Unlock()
	if battle.state == face.GameOver{
		log4g.Info("111")
		return
	}
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
	battle.RLock()
	cmd := battle.frameCache[battle.FrameCount]
	battle.FrameCount++
	battle.RUnlock()
	if cmd != nil && len(cmd) > 0{
		battle.msg.Cmd = cmd
	}
	hasPlayerNoLeave := false
	for _, v := range battle.Players {
		if v != nil && v.IsConnected(){
			hasPlayerNoLeave = true
			battle.GameServer.WriteInnerMsg(v.GetGateSession(), v.GetRoleId(), 5012, battle.msg)
		}
	}
	if hasPlayerNoLeave == false || battle.PlayerCount == 1{
		//全部离线或者只留下一个玩家
		battle.AllPlayerLeave()
	}
}
//检测所有玩家离线,删除该战斗
func (battle *Battle)AllPlayerLeave()  {
	//删除房间
	if battle.BattleType == face.FreeRoomBattleType{
		if battle.GameServer.GetRoomManager().DeleteRoom(battle.RoomId){
			battle.Finish()
		}else{
			log4g.Infof("删除房间失败[%d]",battle.RoomId)
			battle.Finish()
		}
	}else if battle.BattleType == face.PaiWeiBattleType{
		room := battle.GameServer.GetMatchManager().GetMatchBattleRoom(battle.RoomId)
		if room != nil{
			room.Clear()
			battle.Finish()
		}
	}
}
func (battle *Battle)CheckAllLeave(){
	hasPlayerNoLeave := false
	for _, v := range battle.Players {
		if v != nil && v.IsConnected(){
			hasPlayerNoLeave = true
		}
	}
	if hasPlayerNoLeave == false{
		battle.AllPlayerLeave()
	}
}
func (battle *Battle)RemovePlayer(roleId int64)bool{
	battle.Lock()
	defer battle.Unlock()
	for k,v := range battle.Players {
		if v != nil {
			if v.GetRoleId() == roleId {
				battle.PlayerCount--
				battle.Players[k] = nil
				//log4g.Infof("移除战斗成员成功![%d]",roleId)
			}
		}
	}
	//看是否战斗还剩下一个人，如果一个人就代表战斗结束
	if battle.PlayerCount == 1{
		//代表游戏结束
		return true
	}else{
		return false
	}
}
func (battle *Battle)GetSaveFrames() map[int32][]*message.Command{
	return battle.frameCache
}
func (battle *Battle)GetFrameCount()int32{
	battle.RLock()
	defer battle.RUnlock()
	return battle.FrameCount
}
//func (battle *Battle)SurrenderPlayer(roleId int64,groupId int32){
//	for _, v := range battle.Players {
//		if v != nil && v.IsConnected(){
//			if v.GetRoleId() != roleId{
//				battle.GameServer.WriteInnerMsg(v.GetGateSession(),v.GetRoleId(),5101,msg)
//			}
//		}
//	}
//}