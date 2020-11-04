package matchManager

import (
	"cmgameserver/face"
	"math/rand"
	"github.com/golang/protobuf/proto"
	"time"
	"sync"
	"cmgameserver/message"
)

//type MatchBattleRoomStage int32
//
//const(
//	SelectStage MatchBattleRoomStage = 0
//	PrepareStage MatchBattleRoomStage = 1
//)
//战斗房间
//里面有玩家，不同颜色，
//有定时，时间到就开始游戏

type MatchBattleRoom struct {
	RoomId int32
	AllGroupMatchPlayer  []face.IMatchPlayer//不同队伍的玩家
	GroupIdPool [6]int32
	MatchManger *MatchManager
	selectTimer *time.Timer
	//prepareTimer *time.Timer
	loadTimer  *time.Timer//玩家加载1分钟，有人没加载完成，就直接进入游戏
	Stop bool
	Seed int32//战斗随机种子
	sync.RWMutex
}

func (this *MatchBattleRoom)Start(){
	go func() {
		<-this.selectTimer.C
		this.SetStop(true)
		//进入准备阶段
		this.EnterPrepare()
	}()
}
func (this *MatchBattleRoom)GetRoomId()int32{
	return this.RoomId
}
func (this *MatchBattleRoom)GetRoomSeed()int32{
	return  this.Seed
}
func (this *MatchBattleRoom)GetAllPlayer()[]face.IMatchPlayer{
	return this.AllGroupMatchPlayer
}
func (this *MatchBattleRoom)EnterPrepare(){
	//发送消息给所有玩家
	returnMsg := new(message.M2C_MatchBattleRoomEnterPrepare)
	this.SendMstToAllPlayer(5071,returnMsg)
	go func() {
		<-time.After(time.Second * 10)
		//进入战斗状态
		startBattleMsg := new(message.M2C_StartPaiWeiGame)
		startBattleMsg.Seed = this.Seed
		this.SendMstToAllPlayer(5070,startBattleMsg)
		if this.loadTimer == nil{
			this.loadTimer = time.NewTimer(time.Second * 60)
		}
		<-this.loadTimer.C
		//进入游戏
		returnMsg := new(message.M2C_StartBattle)
		this.SendMstToAllPlayer(5011,returnMsg)
		//this.Clear()
	}()
}
func (this *MatchBattleRoom)AddAllPlayer(all map[int]face.IMatchTeam){
	var posIndex int32
	for _,v := range all{
		if v != nil{
			allPlayer := v.GetAllMatchPlayer()
			this.AllGroupMatchPlayer = append(this.AllGroupMatchPlayer, allPlayer...)
		}
	}
	posIndex = 0
	//分配groupId和pos
	for _,v := range this.AllGroupMatchPlayer{
		if v != nil{
			id := this.getRandomId()
			v.SetGroupId(id)
			v.SetMatchRoomId(this.RoomId)
			v.SetBInBattleRoom(true)
			role := v.GetOnlineRole()
			role.SetRoomId(this.RoomId)
			role.SetBattleSeed(this.Seed)
			v.SetPosIndex(posIndex)
			posIndex += 1
		}
	}
}

//取得随机id
func (this *MatchBattleRoom)getRandomId() int32{
Random:
	randomIndex := rand.Int31n(6)
	if this.GroupIdPool[randomIndex] == -1{
		goto Random
	}else{
		id := this.GroupIdPool[randomIndex]
		this.GroupIdPool[randomIndex] = -1
		return id
	}
}
func (this *MatchBattleRoom)SetPlayerPrepare(roleId int64,value bool)bool{
	result := false
	for _,v := range this.AllGroupMatchPlayer{
		if v != nil{
			if v.GetRoleId() == roleId{
				v.SetPrepare(value)
				result = true
				break
			}
		}
	}
	if value && result{
		//看其他玩家是否准备完毕。准备完毕就开始进入准备阶段
		allPrepare := true
		for _,v := range this.AllGroupMatchPlayer{
			if v != nil{
				if v.GetPrepare() == false{
					allPrepare = false
					break
				}
			}
		}
		if allPrepare{
			if this.GetStop() == false{
				this.selectTimer.Stop()
				this.SetStop(true)
				this.EnterPrepare()
			}
		}
	}
	return result
}
func (this *MatchBattleRoom)SendMstToAllPlayer(msgId int,msg proto.Message){
	for _,v := range this.AllGroupMatchPlayer{
		if v != nil{
			role := v.GetOnlineRole()
			if role != nil{
				this.MatchManger.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),msgId,msg)
			}
		}
	}
}
func (this *MatchBattleRoom)SetStop(value bool){
	this.Lock()
	defer this.Unlock()
	this.Stop = value
}
func (this *MatchBattleRoom)GetStop()bool{
	this.RLock()
	defer this.RUnlock()
	return  this.Stop
}
func (this *MatchBattleRoom)Clear(){
	if len(this.AllGroupMatchPlayer) > 0{
		for _,v := range this.AllGroupMatchPlayer{
			if v != nil{
				v.Clear()
			}
		}
		this.AllGroupMatchPlayer = nil
	}
	this.selectTimer = nil
	this.MatchManger.RemoveMatchBattleRoom(this.RoomId)
	this.MatchManger = nil
}
func (this *MatchBattleRoom)ClearNotifyAllPlayer(){
	returnMsg := new(message.M2C_RemoveMatchBattleRoom)
	this.SendMstToAllPlayer(5074,returnMsg)
	this.Clear()
}
//检测是否所有玩家都已经加载完成
func (this *MatchBattleRoom)CheckAllPlayerLoadFinished()(bool,*[4]face.IOnlineRole){
	////清理这个房间
	//this.Clear()
	allLoadFinished := true
	for _,v := range this.AllGroupMatchPlayer{
		if v != nil{
			role := v.GetOnlineRole()
			if role != nil{
				if role.IsLoadFinished() == false{
					allLoadFinished = false
					break
				}
			}
		}
	}
	//说明所有人加载完成
	if allLoadFinished{
		var allPlayer [4]face.IOnlineRole
		index := 0
		for _,v := range this.AllGroupMatchPlayer{
			if v != nil{
				allPlayer[index] = v.GetOnlineRole()
				index++
			}
		}
		//进入游戏
		this.loadTimer.Stop()
		//通知所有进入游戏
		returnMsg := new(message.M2C_StartBattle)
		this.SendMstToAllPlayer(5011,returnMsg)
		//this.Clear()
		return true,&allPlayer
	}else{
		return false,nil
	}
}