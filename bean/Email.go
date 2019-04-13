package bean

import "cmgameserver/message"

type Email struct{
	EmailId int32
	EmailTime int64
	Get   bool
}
func (this *Email)ToMessageData() *message.Email{
	data := new(message.Email)
	data.EmailId = this.EmailId
	data.Time = this.EmailTime
	data.BGet = this.Get
	return data
}