package bean

import (
	"time"
	"github.com/bianchengxiaobei/cmgo/log4g"
	"cmgameserver/tool"
)

type CardType int32
//月卡类型
const(
	NoneCard  CardType = -1
	YueCard CardType = 0//月卡
	NianCard CardType  = 1//年卡
	JiCard CardType = 2//季卡=>3个月
	ZhongShenCard CardType = 3//终身卡
)

type DiamType int32

const(
	Six        DiamType = 0 //6颗钻石 6
	Fourteen   DiamType = 1 //14钻石  12
	ThirtyFive DiamType = 2 //35颗钻石 30
	OneBag     DiamType = 3 //50钻石  42
	TwoBag     DiamType = 4 //100钻石  78
	OneBox     DiamType = 5 //一箱150  108
	TwoBox     DiamType = 6 //两箱300   204
)


type Card struct {
	CardType CardType
	BuyTime time.Time
	LastGetTime time.Time
}
func (this *Card)Clear(){
	this.BuyTime = time.Now()
	this.CardType = NoneCard
}
//是否过期了
func (this *Card)IsTimeOut(now time.Time)bool{
	if this.CardType == YueCard{
		nextTime := this.BuyTime.AddDate(0,1,0).Unix()
		nowTime := now.Unix()
		if nowTime > nextTime{
			//说明是过期了
			return true
		}
	}else if this.CardType == JiCard{
		if this.BuyTime.AddDate(0,3,0).Unix() < now.Unix(){
			//说明是过期了
			return true
		}
	}else if this.CardType == NianCard{
		if this.BuyTime.AddDate(1,0,0).Unix() < now.Unix(){
			//说明是过期了
			return true
		}
	}
	return false
}
func (this *Card)GetDayDiam(now time.Time,cardConfig CardInfoConfig)int32{
	if this.CardType == NoneCard{
		return 0
	}
	if tool.IsSameDay(now,this.LastGetTime){
		//说明领取过了
		return 0
	}
	var allValue int32
	data := cardConfig.Cards[int32(this.CardType)]
	nextDay := this.BuyTime.AddDate(0,1,0).Day()
	nowDay := now.Day()
	buyDay := this.BuyTime.Day()
	if nextDay == nowDay || buyDay == nowDay{
		allValue += data.MonthGet
	}
	if tool.IsSameDay(this.BuyTime,now){
		this.LastGetTime = now
	}
	leftDay := nowDay - this.LastGetTime.Day()
	if leftDay <= 0{
		//log4g.Infof("月卡天数[%d]",leftDay)
		leftDay = 1
	}
	allValue += data.DayGet * int32(leftDay)
	this.LastGetTime = now
	log4g.Infof("获得月卡[%d]奖励[%d]",this.CardType,allValue)
	return allValue
}