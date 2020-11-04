package bean

type ServerEmailData struct {
	EmailId   	int32
	Awards		[]int32
	OtherAwards map[int32]int32
}

type ServerEmailConfig struct {
	Data map[int32]ServerEmailData
} 