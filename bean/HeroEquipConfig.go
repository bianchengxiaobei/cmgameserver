package bean

type EItemType int
const (
	Body EItemType = iota
	Weap
	Shoe
	Consume//消耗品
	Material//材料
	TouKui
	Box
)
type EEquipApplyType int
const (
	HeroType EEquipApplyType = iota
	Soldier
	RoleType
	None
)
type HeroQualityMapEquipConfig struct {
	Data map[EItemQuality][]ServerEquipData
}
type HeroItemIdMapEquipConfig struct {
	Data map[int32]ServerEquipData
}
type ServerEquipData struct {
	ItemId int32
	ItemType EItemType
	ApplyType EEquipApplyType
	SellGold int32
	PlayerType int32
	GuanFang bool
}