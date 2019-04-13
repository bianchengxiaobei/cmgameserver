package bean

type Item struct {
	ItemId int32
	ItemSeed int32
	ItemNum int32
	ItemTime int64
}

func (item *Item)Clear(){
	item.ItemId = 0
	item.ItemSeed = 0
	item.ItemNum = 0
	item.ItemTime = 0
}