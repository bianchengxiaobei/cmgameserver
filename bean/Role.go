
package bean

import "time"

type Role struct {
	RoleId		int64
	UserId		int64
	NickName	string
	MailAddress string//绑定邮箱地址
	QQ			int32//绑定QQ账号
	WeiXin		string//绑定微信账号
	Phone		int32//绑定手机账号
	GuideFinished  bool //是否新手关卡完成了
	ServerId	int32
	Level	int32
	AvatarId	int32
	Gold	int32
	RankScore int32
	Diam    int32
	Exp     int32
	Sex 	int32
	HeroCount int32
	SaiJiId  int32
	OnlineGetDiam int32//在线获得的钻石奖励
	Card   []Card
	CardDiamValue int32//累积功能卡获得的数量
	Sign string
	MaxBagNum int
	MaxEmailNum int
	Items  []Item
	Trans  []Tran
	Emails []Email
	GiftCode []string
	AchieveRecord Achievement
	WinLevel []int32
	DayGetTask	[]int32
	TaskSeed	int32
	Achievement []int32
	LoginTime time.Time
	GetSign			bool
	FreeSoldierData   [4]FreeSoldierData
	DunJiaDaoDun  FreeSoldierData
	BuJiaDaoDun   FreeSoldierData
	DunJiaArrow  FreeSoldierData
	BuJiaArrow   FreeSoldierData
	DunJiaSpear  FreeSoldierData
	BuJiaSpear   FreeSoldierData
	DunJiaFaShi  FreeSoldierData
	BuJiaFaShi   FreeSoldierData
	IsGM		bool //是否是游戏管理员
}
type FreeSoldierData struct {
	PlayerType		int32
	CarrierType     int32
	TouKuiId		int32
	BodyId			int32
	WeapId			int32
}