package bean

import "cmgameserver/message"

type Email struct{
	EmailIndex int32
	EmailTime int64
	Title string
	Content string
	AwardList []int32
	Get   bool
	Valid bool//是否是有效的
}
func (this *Email)ToMessageData() *message.Email{
	data := new(message.Email)
	data.Title = this.Title
	data.Content = this.Content
	data.AwardList = this.AwardList
	data.Time = this.EmailTime
	data.EmailIndex = this.EmailIndex
	data.BGet = this.Get
	return data
}
func (this *Email)Clear(){
	this.Title = ""
	this.Content = ""
	this.EmailIndex = 0
	this.EmailTime = 0
	this.Get = false
	this.Valid = false
}