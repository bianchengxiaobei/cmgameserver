package bean

type ServerEmailData struct {
	EmailId   	int32
	Awards		[]int32
}

type ServerEmailConfig struct {
	Data map[int32]ServerEmailData
} 