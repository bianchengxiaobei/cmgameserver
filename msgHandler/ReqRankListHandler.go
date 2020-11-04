package msgHandler

import (
	"cmgameserver/face"
	"github.com/bianchengxiaobei/cmgo/network"
	"cmgameserver/message"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/bean"
	"time"
)

type ReqRankListHandler struct {
	GameServer face.IGameServer
}

func (handler *ReqRankListHandler) Action(session network.SocketSessionInterface, msg interface{}) {
	if innerMsg, ok := msg.(network.InnerWriteMessage); ok {
		if protoMsg, ok := innerMsg.MsgData.(*message.C2M_ReqRankList); ok {
			role := handler.GameServer.GetRoleManager().GetOnlineRole(innerMsg.RoleId)
			if role == nil{
				return
			}
			//查询段位榜单的列表
			updateTime := handler.GameServer.GetRankListTime()
			nowTime := time.Now().Unix()
			deltaTime := nowTime - updateTime
			//每隔一个小时
			if updateTime == 0 || deltaTime > 3600{
				//更新排行榜
				//段位榜
				dbSession := handler.GameServer.GetDBManager().Get()
				if dbSession == nil{
					return
				}
				handler.GameServer.SetRankListTime(nowTime)
				allRole := make([]bean.Role,30)
				c := dbSession.DB("sanguozhizhan").C("Role")
				err := c.Find(nil).Sort("-rankscore").Limit(30).All(&allRole)
				if err != nil {
					log4g.Errorf("更新RankLevel出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
					return
				}
				handler.UpdateRankLevelRankList(allRole)
				err = c.Find(nil).Sort("-herocount").Limit(30).All(&allRole)
				if err != nil {
					log4g.Errorf("更新HeroCount出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
					return
				}
				handler.UpdateRoleHeroCountList(allRole)
				err = c.Find(nil).Sort("-level").Limit(30).All(&allRole)
				if err != nil {
					log4g.Errorf("更新Level出错[%s],RoleId:%d", err.Error(), innerMsg.RoleId)
					return
				}
				handler.UpdateRoleLevelRankList(allRole)
			}
			returnMsg := new(message.M2C_RankListResult)
			if returnMsg.RankListItem == nil{
				returnMsg.RankListItem = make([]*message.RankListItem,0)
			}
			if protoMsg.RankType == message.ERankListType_RankLevelIndex{
				returnMsg.RankType = message.ERankListType_RankLevelIndex
				items := handler.GameServer.GetRankLevelRankList()
				for _,v := range items{
					if &v != nil && v.NickName != ""{
						if v.RoleId == innerMsg.RoleId{
							returnMsg.MySortOrder = v.RankNum
						}
						messageItem := v.ToMessageData()
						returnMsg.RankListItem = append(returnMsg.RankListItem, messageItem)
					}
				}
			}else if protoMsg.RankType == message.ERankListType_HeroIndex{
				returnMsg.RankType = message.ERankListType_HeroIndex
				items := handler.GameServer.GetRoleHeroCountRankList()
				for _,v := range items{
					if &v != nil && v.NickName != ""{
						if v.RoleId == innerMsg.RoleId{
							returnMsg.MySortOrder = v.RankNum
						}
						messageItem := v.ToMessageData()
						returnMsg.RankListItem = append(returnMsg.RankListItem, messageItem)
					}
				}
			}else if protoMsg.RankType == message.ERankListType_RoleLevelIndex{
				returnMsg.RankType = message.ERankListType_RoleLevelIndex
				items := handler.GameServer.GetRoleLevelRankList()
				for _,v := range items{
					if &v != nil && v.NickName != ""{
						if v.RoleId == innerMsg.RoleId{
							returnMsg.MySortOrder = v.RankNum
						}
						messageItem := v.ToMessageData()
						returnMsg.RankListItem = append(returnMsg.RankListItem, messageItem)
					}
				}
			}
			handler.GameServer.WriteInnerMsg(role.GetGateSession(), role.GetRoleId(), 5059, returnMsg)
		}
	}
}
func (handler *ReqRankListHandler)UpdateRankLevelRankList(allRoles []bean.Role){
	rankList := handler.GameServer.GetRankLevelRankList()
	for k,v := range allRoles {
		if &v != nil{
			item := &rankList[k]
			if item != nil{
				item.RankNum = int32(k + 1)//排行从1开始
				item.NickName = v.NickName
				item.Value = v.RankScore
				item.RankLevel = v.RankScore
				item.RoleId = v.RoleId
			}else {
				item = new(bean.RankListItem)
				item.RankNum = int32(k + 1)
				item.NickName = v.NickName
				item.Value = v.RankScore
				item.RankLevel = v.RankScore
				item.RoleId = v.RoleId
				rankList[k] = *item
			}
		}
	}
}
func (handler *ReqRankListHandler) UpdateRoleHeroCountList(allRoles []bean.Role){
	rankList := handler.GameServer.GetRoleHeroCountRankList()
	for k,v := range allRoles {
		if &v != nil{
			item := &rankList[k]
			if item != nil{
				item.RankNum = int32(k + 1)
				item.NickName = v.NickName
				item.Value = v.HeroCount
				item.RoleId = v.RoleId
				item.RankLevel = v.RankScore
			}else {
				item = new(bean.RankListItem)
				item.RankNum = int32(k + 1)
				item.NickName = v.NickName
				item.Value = v.HeroCount
				item.RankLevel = v.RankScore
				item.RoleId = v.RoleId
				rankList[k] = *item
			}
		}
	}
}
func (handler *ReqRankListHandler)UpdateRoleLevelRankList(allRoles []bean.Role){
	rankList := handler.GameServer.GetRoleLevelRankList()
	for k,v := range allRoles {
		if &v != nil{
			item := &rankList[k]
			if item != nil{
				item.RankNum = int32(k + 1)
				item.NickName = v.NickName
				item.Value = v.Level
				item.RankLevel = v.RankScore
				item.RoleId = v.RoleId
			}else {
				item = new(bean.RankListItem)
				item.RankNum = int32(k + 1)
				item.NickName = v.NickName
				item.Value = v.Level
				item.RankLevel = v.RankScore
				item.RoleId = v.RoleId
				rankList[k] = *item
			}
		}
	}
}