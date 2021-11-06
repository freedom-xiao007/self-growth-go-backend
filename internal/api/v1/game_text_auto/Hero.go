package game_text_auto

import (
	"github.com/kamva/mgm/v3"
	"github.com/mozillazg/go-pinyin"
)

// Hero 角色
// spirit 灵魂力：
//    用于维持角色的基本活动，当为0时，视为死亡或者冻结状态
//    会随着时间的推移而消耗
//    每次活动会消耗灵魂
//    用于维持角色的基本活动，当为0时，视为死亡或者冻结状态
//
// spiritAttack 神识：角色的灵魂力，用来计算角色的灵魂攻击力
//
// spiritDefence 神防：角色的神识防御力
//
// bleed 血气：角色生命值
//
// strong 力量：角色的力量，用来计算角色的物理攻击力
//
// shooting 命中率：角色的技巧，用来计算角色的命中率、必杀率和大部分技能的触发率
//
// attackSpeed 攻速：角色的速度。If的追击阈值是5，也就是说，当一名角色的速度高于敌方5点及时，该角色可以在敌方攻击后再攻击一次。速度也是影响角色回避率的属性
//
// dodge 闪避：角色的幸运，主要影响必杀回避（运%），对命中与回避也有些许影响（运/2%）
//
// defence 防御：角色的物理防御
//
// moveSpeed 移动力：角色与一回合内在平地可以移动的格子数量，基础上限为10（算上鞋子）
//
// level 星级：每提升已星级 所需碎片翻倍；基础 10；最大1000；星级无上限
//
// chip 碎片：每提升已星级 所需碎片翻倍；基础 10；最大1000；星级无上限
type Hero struct {
	mgm.DefaultModel `bson:",inline"`
	NameZW             string `json:"name_zw"`
	NamePY             string `json:"name_py"`
	Description      string `json:"description"`
	Spirit int64 `json:"spirit"`
	SpiritAttack int64 `json:"spiritAttack"`
	SpiritDefence int64 `json:"spiritDefence"`
	Bleed int64 `json:"bleed"`
	Strong int64 `json:"strong"`
	Shooting int64 `json:"shooting"`
	AttackSpeed int64 `json:"attackSpeed"`
	Dodge int64 `json:"dodge"`
	Defence int64 `json:"defence"`
	MoveSpeed int64 `json:"moveSpeed"`
	Level int64 `json:"level"`
	Chip int64 `json:"chip"`
}

func NewHero(name, desc string) *Hero {
	return &Hero{
		NameZW:        name,
		NamePY:        HeroNameZW2PY(name),
		Description:   desc,
		Spirit:        86840,
		SpiritAttack:  1,
		SpiritDefence: 1,
		Bleed:         1000,
		Strong:        1,
		Shooting:      1,
		AttackSpeed:   1,
		Dodge:         1,
		Defence:       1,
		MoveSpeed:     1,
		Level: 1,
		Chip: 0,
	}
}

func HeroNameZW2PY(name string) string {
	namePy := ""
	pyStr := pinyin.NewArgs()
	for _, item := range pinyin.Pinyin(name, pyStr) {
		namePy += item[0]
	}
	return namePy
}