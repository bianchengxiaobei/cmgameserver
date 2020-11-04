package bean

type GiftConfig struct {
	Datas map[string]GiftData
}

type GiftData struct {
	Code string//礼品码
	Awards		[]int32//礼品码奖励
	Gold int32//礼品金币
	Diam int32//礼品钻石
	Exp int32//礼品经验
	HeroId int32//礼品武将
}