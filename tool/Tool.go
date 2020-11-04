package tool

import (
	"time"
	"strconv"
)
func GetHashCode(num int64)int{
	println(num)
	println(int32(num))
	key := strconv.FormatInt(num,10)
	var index int = 0
	index = int(key[0])
	for k := 0; k < len(key); k++ {
		//1103515245是个好数字，使通过hashCode散列出的0-15的数字的概率是相等的
		index *= 1103515245 + int(key[k])
	}
	index >>= 27
	index &= 16 - 1
	return index
}
func GetRankLevelFromRankSocre(score int32)int32{
	if score >= 0 && score < 200{
		//青铜
		return 0
	}else if score >= 200 && score < 500{
		//白银
		return 1
	}else if score >= 500 && score < 800{
		//黄金
		return 2
	}else if score >= 800 && score < 1200{
		//铂金
		return 3
	}else{
		//钻石
		return 4
	}
}
func GetAddRankScoreFormRankLevel(score int32,win bool)int32{
	var reslut int32
	reslut = 30
	if score >= 0 && score < 200{
		//青铜
		if win == false{
			reslut = 10
		}
	}else if score >= 200 && score < 500{
		//白银
		if win == false{
			reslut = -5
		}
	}else if score >= 500 && score < 800{
		//黄金
		if win == false{
			reslut = -15
		}
	}else if score >= 800 && score < 1200{
		//铂金
		if win == false{
			reslut = -20
		}
	}else{
		//钻石
		if win == false{
			reslut = -30
		}
	}
	return reslut
}

//const(
//	Six        DiamType = 0 //6颗钻石 6  48
//	Fourteen   DiamType = 1 //14钻石  12  112
//	ThirtyFive DiamType = 2 //35颗钻石 30  280
//	OneBag     DiamType = 3 //50钻石  42   400
//	TwoBag     DiamType = 4 //100钻石  78    800
//	OneBox     DiamType = 5 //一箱150  108   1200
//	TwoBox     DiamType = 6 //两箱300   204   2400
//)
func GetDiamValueByDiamType(diamType int32) int32 {
	switch diamType {
	case 0:
		return 48
	case 1:
		return 112
		break
	case 2:
		return  280
	case 3:
		return 400
	case 4:
		return 800
	case 5:
		return 1200
	case 6:
		return 2400
	}
	return 0
}
func GetBoxPriceByBoxId(boxId int32,buyType int32) int32{
	var value int32
	if buyType == 0{
		//金币
		if boxId == 200000{
			value = 500
		}else if boxId == 200001{
			value = 1000
		}else if boxId == 200002{
			value = 2000
		}
		return value
	}else{
		if boxId == 200000{
			value = 50
		}else if boxId == 200001{
			value = 100
		}else if boxId == 200002{
			value = 200
		}else if boxId == 200010{
			value = 500
		}
		return value
	}
}
func GetHeroUpgradeNeedGold(heroLevel int32)int32{
	return  50 * heroLevel * heroLevel - heroLevel * 25
}
func IsSameDay(now time.Time,lastLogin time.Time)bool{
	nowYear,nowMonth,nowDay := now.Date()
	lastYear,lastMonth,lastDay := lastLogin.Date()
	if nowYear == lastYear && nowMonth == lastMonth && nowDay == lastDay{
		return true
	}
	return false
}