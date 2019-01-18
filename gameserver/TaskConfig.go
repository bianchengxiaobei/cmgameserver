package gameserver

type TaskConfig struct {
	Tasks map[int32]TaskConfigData
}
type AwardType int
const (
	Gold AwardType = iota
	Exp
	Diam
	Item
)
type TaskConfigData struct {
	TaskId	int32
	Award map[AwardType]int
}
type AchieveConfigData struct {
	AchieveId	int32
	Award map[AwardType]int
}
type AchieveConfig struct {
	Achieves map[int32]AchieveConfigData
} 