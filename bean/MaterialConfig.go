package bean

type EMaterialType int
const (
	EGold EMaterialType = iota
	EExp
)

type ServerMaterialConfig struct {
	Data map[int32]ServerMaterial
}
type ServerMaterial struct {
	ItemId int32
	Quality EItemQuality
	MatType EMaterialType
	MatValue int
	SellGold int
}