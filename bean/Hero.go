package bean

type Hero struct {
	RoleId	int64
	HeroId	int32
	Level 	int32
	//Exp     int32
	ItemIds [3]Item
	LearnSkill [10]int32
	SkillPoint int32
	//NextLevelExp int32
	Skill1 int32
	Skill2 int32
}

func (hero *Hero)UpgradeLevel() (level int32,skillPoint int32){
	hero.Level++
	hero.SkillPoint++
	return hero.Level,hero.SkillPoint
}
func (hero *Hero)HeroLearnSkill(skillId int32,needSkillPoint int32)bool{
	if hero.SkillPoint < needSkillPoint{
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
			hero.SkillPoint -= needSkillPoint
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