package bean

type Hero struct {
	RoleId	int64
	HeroId	int32
	Level 	int32
	Exp     int32
	ItemIds [3]int32
	LearnSkill [10]int32
	SkillPoint int32
	NextLevelExp int32
	Skill1 int32
	Skill2 int32
}

func (hero *Hero)AddExp(exp int32) int32{
	hero.Exp += exp
	hero.NextLevelExp = hero.Level * hero.Level * 25 + 25 * hero.Level
	//英雄升级所用经验
	if hero.Exp >= hero.NextLevelExp{
		hero.Level++
		hero.SkillPoint++
		hero.NextLevelExp = hero.Level * hero.Level * 25 + 25 * hero.Level
	}
	return hero.Exp
}
func (hero *Hero)HeroLearnSkill(skillId int32)bool{
	if hero.SkillPoint > 0{
		return  false
	}
	for k,v := range hero.LearnSkill{
		if v > 0{
			if v == skillId{
				return false
			}
		}else{
			//这里扣除技能点
			hero.LearnSkill[k] = skillId
			return true
		}
	}
	return false
}
func (hero *Hero)HasSkill(skillId int32)bool{
	for _,v := range hero.LearnSkill{
		if v > 0{
			if v == skillId{
				return true
			}
		}
	}
	return false
}