package bean

type BoxItemInfoData struct {
	BoxChange [4]int
	Quality EItemQuality
}
type BoxItemInfoConfig struct {
	BoxList map[int32]BoxItemInfoData
}
type EItemQuality int

const (
	White EItemQuality = iota
	Green
	Blue
	Orange
)
