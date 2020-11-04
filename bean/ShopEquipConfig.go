package bean

type ServerShopEquipConfig struct {
	ShopEquip map[int32]ShopEquipServerData
}
type ShopEquipServerData struct {
	EquipId int32
	Price int32
}