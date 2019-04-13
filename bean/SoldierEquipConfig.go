package bean

type SoldierQualityEquipConfig struct {
	Data map[EItemQuality][]ServerEquipData
}
type SoldierItemIdEquipConfig struct {
	Data map[int32]ServerEquipData
}