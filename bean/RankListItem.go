package bean

import "cmgameserver/message"

type RankListItem struct {
	RankNum int32
	RankLevel int32
	NickName string
	Value int32
	RoleId int64
} 

func (this RankListItem)ToMessageData() *message.RankListItem{
	item := new(message.RankListItem)
	item.RankLevel = this.RankLevel
	item.Value = this.Value
	item.NickName = this.NickName
	item.RankNum = this.RankNum
	return item
}