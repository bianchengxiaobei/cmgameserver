package bean


type CardInfoConfigData struct {
	Price int32
	DayGet int32
	MonthGet int32
}
type CardInfoConfig struct {
	Cards map[int32]CardInfoConfigData
}