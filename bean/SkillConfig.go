package bean

type ServerSkillConfig struct {
	Datas  map[int32]ServerSkillConfigData
}
type ServerSkillConfigData struct {
	SkillId int32
	LearnPoint int32
}