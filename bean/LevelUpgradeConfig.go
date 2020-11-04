package bean

type LevelUpgradeConfig struct{
	Datas map[int32]LevelUpgradeData
}
type LevelUpgradeData struct {
	Level int32
	Awards []int32
	Gold int32
	Diam int32
	HeroId int32
}