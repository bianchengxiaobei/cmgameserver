package bean


type ShopHeroCardData struct {
	CardId int32
	Price int32
}
type ShopHeroCardConfig struct {
	Cards map[int32]ShopHeroCardData
}
