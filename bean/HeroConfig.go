package bean

type HeroConfig struct {
	Data map[int32]HeroData
}
type HeroData struct {
	HeroId int32
	GuanFangBodyId int32
	GuanFangWeapId int32
	GuanFangShoeId int32
}
