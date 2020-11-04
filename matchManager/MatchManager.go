package matchManager

import (
	"cmgameserver/face"
	"cmgameserver/message"
	"sync"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"sync/atomic"
	"github.com/golang/protobuf/proto"
	"time"
)


var MatchBattleRoomId int32

type MatchManager struct {
	GameServer face.IGameServer
	//可能排位有好几种游戏模式
	AllSingleMatchRoom map[int32]face.IMatchRoom
	AllDoubleMatchRoom map[int32]face.IMatchRoom
	AllMatchTeam       map[int32]face.IMatchTeam
	AllMatchBattleRoom map[int32]face.IMatchBattleRoom
	sync.RWMutex
}

func NewMatchManager(server face.IGameServer) *MatchManager{
	m := MatchManager{
		GameServer:         server,
		AllSingleMatchRoom: make(map[int32]face.IMatchRoom),
		AllMatchTeam:       make(map[int32]face.IMatchTeam),
		AllMatchBattleRoom:make(map[int32]face.IMatchBattleRoom),
	}
	return &m
}
func NewSingleMatchRoom(manager *MatchManager) MatchSingleRoom {
	room := MatchSingleRoom{
		AllBaiYinMatchTeams:make(map[int]face.IMatchTeam,4),
		AllBoJingMatchTeams:make(map[int]face.IMatchTeam,4),
		AllHuangJingMatchTeams:make(map[int]face.IMatchTeam,4),
		AllQingtongMatchTeams:make(map[int]face.IMatchTeam,4),
		AllZuanShiMatchTeams:make(map[int]face.IMatchTeam,4),
		QingTongNum:0,
		BaiYinNum:0,
		HuangJingNum:0,
		BoJingNum:0,
		ZuanShiNum:0,
		MatchManager:manager,
	}
	return room
}
func NewDoubleMatchRoom(manager *MatchManager) MatchDoubleRoom {
	room := MatchDoubleRoom{
		AllBaiYinMatchTeams:make(map[int]face.IMatchTeam,4),
		AllBoJingMatchTeams:make(map[int]face.IMatchTeam,4),
		AllHuangJingMatchTeams:make(map[int]face.IMatchTeam,4),
		AllQingtongMatchTeams:make(map[int]face.IMatchTeam,4),
		AllZuanShiMatchTeams:make(map[int]face.IMatchTeam,4),
		QingTongNum:0,
		BaiYinNum:0,
		HuangJingNum:0,
		BoJingNum:0,
		ZuanShiNum:0,
		MatchManager:manager,
	}
	return room
}
func NewMatchBattleRoom(id int32,matchManager *MatchManager)MatchBattleRoom{
	battleRoom := MatchBattleRoom{
		RoomId:id,
		AllGroupMatchPlayer:make([]face.IMatchPlayer,1),
		GroupIdPool:[6]int32{0,1,2,3,4,5},
		selectTimer:time.NewTimer(time.Second * 30),
		Stop:false,
		Seed:int32(time.Now().Unix()),
		MatchManger:matchManager,
	}
	return  battleRoom
}


//玩家点击排位按钮，创建排位队伍，不加入匹配，只是创建
func (this *MatchManager)ReqAutoCreateMatchTeam(role face.IOnlineRole){
	matchPlayer := role.GetMatchPlayer()
	if matchPlayer != nil{
		if matchPlayer.GetMatchTeamId() > 0{
			//说明已经存在就放回
			return
		}
		team := CreateMatchTeam(matchPlayer)
		if team != nil{
			this.AddMatchTeam(team)
			//然后发送给客户端，自己的匹配队伍信息
			returnMsg := new(message.M2C_MatchTeamInfo)
			msgPlayer := new(message.MatchTeamPlayer)
			msgPlayer.RoleId = role.GetRoleId()
			msgPlayer.NickName = role.GetNickName()
			msgPlayer.AvatarId = role.GetAvatarId()
			p := matchPlayer
			msgPlayer.BOwner = p.GetBOwner()
			msgPlayer.PosIndex = p.GetPosIndex()
			returnMsg.Players = append(returnMsg.Players, msgPlayer)
			//发送
			this.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5063, returnMsg)
		}
	}
}
//添加一个匹配队伍
func (this *MatchManager)AddMatchTeam(team face.IMatchTeam){
	this.Lock()
	defer this.Unlock()
	this.AllMatchTeam[team.GetMatchTeamId()] = team
}
//移除一个匹配队伍(发送消息)
func (this *MatchManager)RemoveMatchTeam(teamId int32){
	this.Lock()
	defer this.Unlock()
	team := this.AllMatchTeam[teamId]
	if team == nil{
		return
	}else{
		//发送消息给房间内的所有玩家
		team.Clear(this.GameServer)
		delete(this.AllMatchTeam, teamId)
		//log4g.Infof("解散匹配队伍[%d]",teamId)
	}
}
func (this *MatchManager)RemoveMatchTeamWhenEnterRoom(teamId int32){
	this.Lock()
	defer this.Unlock()
	team := this.AllMatchTeam[teamId]
	if team == nil{
		return
	}else{
		team.ClearNoSend()
		delete(this.AllMatchTeam, teamId)
	}
}
//从一个队伍中移除队员
func (this *MatchManager)RemovePlayerFromMatchTeam(matchPlayer face.IMatchPlayer){
	team := this.AllMatchTeam[matchPlayer.GetMatchTeamId()]
	if team == nil{
		return
	}
	if matchPlayer.GetBOwner() == true{
		if team.GetBInMatching() == false{
			//如果是队长的话，直接解散队伍
			this.RemoveMatchTeam(matchPlayer.GetMatchTeamId())
		}else{
			//如果正在匹配中，得先从房间内移除队伍
			if team.GetPlayerSize() == 1{
				room := this.AllSingleMatchRoom[matchPlayer.GetMatchRoomId()]
				if room != nil{
					room.RemoveOneMatchTeam(team)
				}
			}else{
				room := this.AllDoubleMatchRoom[matchPlayer.GetMatchRoomId()]
				if room != nil{
					room.RemoveOneMatchTeam(team)
				}
			}
			//通知客户端
			this.RemoveMatchTeam(matchPlayer.GetMatchTeamId())
		}
	}else{
		//普通队员
		team := this.AllMatchTeam[matchPlayer.GetMatchTeamId()]
		if team != nil{
			team.RemoveMatchPlayer(matchPlayer,this.GameServer)
		}
	}
}
//该匹配队伍开始匹配（根据玩家（房主）的排位分和等级来判断）
func (this *MatchManager)TeamStartMatching(teamId int32,mode face.PaiWeiGameMode){
	team := this.AllMatchTeam[teamId]
	if team == nil{
		return
	}
	modeId := int32(mode)
	if team.GetPlayerSize() == 1{
		room := this.AllSingleMatchRoom[modeId]
		if room == nil{
			//创建房间
			tempRoom := NewSingleMatchRoom(this)
			room = &tempRoom
			this.AllSingleMatchRoom[modeId] = room
		}
		room.AddOneMatchTeam(team)
	}else{
		room := this.AllDoubleMatchRoom[modeId]
		if room == nil{
			//创建房间
			tempRoom := NewDoubleMatchRoom(this)
			room = &tempRoom
			this.AllDoubleMatchRoom[modeId] = room
		}
		room.AddOneMatchTeam(team)
	}
}
func (this *MatchManager)CancelStartMatch(matchPlayer face.IMatchPlayer){
	team := this.AllMatchTeam[matchPlayer.GetMatchTeamId()]
	if team == nil{
		return
	}
	if team.GetBInMatching(){
		//如果正在匹配中，得先从房间内移除队伍
		if team.GetPlayerSize() == 1{
			room := this.AllSingleMatchRoom[matchPlayer.GetMatchRoomId()]
			if room != nil{
				room.RemoveOneMatchTeam(team)
			}
		}else{
			room := this.AllDoubleMatchRoom[matchPlayer.GetMatchRoomId()]
			if room != nil{
				room.RemoveOneMatchTeam(team)
			}
		}
	}else{
		log4g.Infof("不在匹配中")
	}
}
func (this *MatchManager)GetMatchBattleRoom(roomId int32)face.IMatchBattleRoom{
	room := this.AllMatchBattleRoom[roomId]
	if room != nil{
		return room
	}
	return nil
}
//移除战斗房间
func (this *MatchManager)RemoveMatchBattleRoom(roomId int32){
	room := this.AllMatchBattleRoom[roomId]
	if room != nil{
		delete(this.AllMatchBattleRoom, roomId)
		log4g.Infof("移除匹配房间[%d]",roomId)
	}
}
func (this *MatchManager)EnterMatchBattleRoom(teams map[int]face.IMatchTeam){
	id := atomic.AddInt32(&MatchBattleRoomId,1)
	room := NewMatchBattleRoom(id,this)
	this.AllMatchBattleRoom[id] = &room
	room.AddAllPlayer(teams)
	for _,v := range teams{
		if v != nil{
			this.RemoveMatchTeamWhenEnterRoom(v.GetMatchTeamId())
		}
	}
	//发送消息给所有在房间内的玩家
	players := room.AllGroupMatchPlayer
	returnMsg := new(message.M2C_EnterMatchBattleRoom)
	for _,v := range players{
		if v == nil{
			continue
		}
		role := v.GetOnlineRole()
		msgPlayer := new(message.MatchRoomPlayer)
		msgPlayer.RoleId = v.GetRoleId()
		msgPlayer.GroupId = v.GetGroupId()
		msgPlayer.PosIndex = v.GetPosIndex()
		msgPlayer.AvatarId = role.GetAvatarId()
		msgPlayer.NickName = role.GetNickName()
		msgPlayer.CityId = 1
		msgPlayer.Arrower = role.GetSoldierData(0)
		msgPlayer.Daodun = role.GetSoldierData(1)
		msgPlayer.Spear = role.GetSoldierData(2)
		msgPlayer.Fashi = role.GetSoldierData(3)
		returnMsg.Players = append(returnMsg.Players, msgPlayer)
	}
	for _,v := range players{
		if v == nil{
			continue
		}
		role := v.GetOnlineRole()
		this.GameServer.WriteInnerMsg(role.GetGateSession(),role.GetRoleId(),5067,returnMsg)
	}
}
func (this *MatchManager)OnePlayerQuitMatchBattleRoom(player face.IMatchPlayer){
	room := this.AllMatchBattleRoom[player.GetMatchRoomId()]
	if room != nil{
		room.ClearNotifyAllPlayer()
	}
}
func (this *MatchManager)SetMatchRoomPlayerPrepare(roomId int32,roleId int64,value bool)bool{
	room := this.AllMatchBattleRoom[roomId]
	if room != nil{
		return room.SetPlayerPrepare(roleId,value)
	}else{
		log4g.Infof("不存在房间[%d]",roomId)
		return false
	}
}
func (this *MatchManager)SendMsgToAllBattleRoomPlayer(msgId int,msg proto.Message,roomId int32){
	room := this.AllMatchBattleRoom[roomId]
	if room != nil{
		room.SendMstToAllPlayer(msgId,msg)
	}
}
func (this *MatchManager)CheckAllPlayerLoadFinished(roomId int32)(bool,*[4]face.IOnlineRole){
	room := this.AllMatchBattleRoom[roomId]
	if room != nil{
		return room.CheckAllPlayerLoadFinished()
	}else{
		log4g.Infof("不存在房间[%d]",roomId)
		return false,nil
	}
}