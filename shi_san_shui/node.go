package shi_san_shui

import (
	"fmt"
)

type Node struct {
	normalType NormalType
	pokers     []*Poker
	rest       []*Poker
	dui        *Dui
	liangDui   *LiangDui
	sanTiao    *SanTiao
	tieZhi     *TieZhi
	huLu       *HuLu
}

func NewNode() *Node {
	return &Node{}
}

func (this *Node) String() string {
	pokerDesc := ""
	for _, poker := range this.pokers {
		pokerDesc += poker.Desc
	}
	restDesc := ""
	for _, poker := range this.rest {
		restDesc += poker.Desc
	}
	desc := fmt.Sprintf("结点：{类型=【%s】,左节点=【%s】,右节点=【%s】}", this.normalType, pokerDesc, restDesc)
	return desc
}

//CompareInter 同一个牌型下，不同节点间的比较
func (this *Node) CompareInter(other *Node) CompareResult {
	if this.normalType > other.normalType {
		return Better
	}

	if this.normalType < other.normalType {
		return Worse
	}

	switch this.normalType {
	case WU_LONG, TONG_HUA_SHUN, TONG_HUA, SHUN_ZI:
		return ComparePokerScore(this.pokers, other.pokers)
	case DUI_ZI:
		return this.dui.CompareInter(other.dui)
	case LIANG_DUI:
		return this.liangDui.CompareInter(other.liangDui)
	case SAN_TIAO:
		return this.sanTiao.CompareInter(other.sanTiao)
	case TIE_ZHI:
		return this.tieZhi.CompareInter(other.tieZhi)
	case HU_LU:
		return this.huLu.CompareInter(other.huLu)
	default:
		return Same
	}
}

//CompareExternal 不同牌型下，同位置节点间之前的比较
func (this *Node) CompareExternal(other *Node) CompareResult {
	if this.normalType > other.normalType {
		return Better
	}
	if this.normalType < other.normalType {
		return Worse
	}
	switch this.normalType {
	case WU_LONG, TONG_HUA_SHUN, TONG_HUA, SHUN_ZI:
		return ComparePokerScore(this.pokers, other.pokers)
	case DUI_ZI:
		return this.dui.CompareExternal(other.dui)
	case LIANG_DUI:
		return this.liangDui.CompareExternal(other.liangDui)
	case SAN_TIAO:
		return this.sanTiao.CompareExternal(other.sanTiao)
	case TIE_ZHI:
		return this.tieZhi.CompareExternal(other.tieZhi)
	case HU_LU:
		return this.huLu.CompareExternal(other.huLu)
	default:
		return Same
	}
}

func ComparePokerScore(pokers1, pokers2 []*Poker) CompareResult {
	SortPoker(pokers1)
	SortPoker(pokers2)
	count1 := len(pokers1)
	count2 := len(pokers2)
	min := 0
	if count1 > count2 {
		min = count2
	} else {
		min = count1
	}

	for i := 1; i <= min; i++ {
		poker1 := pokers1[count1-i]
		poker2 := pokers2[count2-i]
		if poker1.Score > poker2.Score {
			return Better
		} else if poker1.Score < poker2.Score {
			return Worse
		}
	}

	if count1 > count2 {
		return Better
	}
	if count1 < count2 {
		return Worse
	}
	//平局也是true
	return Same
}

type Dui struct {
	DuiScore  int //对的牌的
	Dan1Score int //最小的单
	Dan2Score int //次小的单
	Dan3Score int //最大的单
}

//CompareInter 对的内部比较
func (this *Dui) CompareInter(other *Dui) CompareResult {
	if this.DuiScore > other.DuiScore {
		return Better
	} else if this.DuiScore < other.DuiScore {
		return Worse
	}

	if this.Dan3Score < other.Dan3Score {
		return Better
	} else if this.Dan3Score > other.Dan3Score {
		return Worse
	}

	if this.Dan2Score < other.Dan2Score {
		return Better
	} else if this.Dan2Score > other.Dan2Score {
		return Worse
	}

	if this.Dan1Score < other.Dan1Score {
		return Better
	} else if this.Dan1Score > other.Dan1Score {
		return Worse
	}
	return Same
}

//CompareExternal 对的外部比较，所有牌越大越好
func (this *Dui) CompareExternal(other *Dui) CompareResult {
	if this.DuiScore > other.DuiScore {
		return Better
	} else if this.DuiScore < other.DuiScore {
		return Worse
	}

	if this.Dan3Score > other.Dan3Score {
		return Better
	} else if this.Dan3Score < other.Dan3Score {
		return Worse
	}

	if this.Dan2Score > other.Dan2Score {
		return Better
	} else if this.Dan2Score < other.Dan2Score {
		return Worse
	}

	if this.Dan1Score > other.Dan1Score {
		return Better
	} else if this.Dan1Score < other.Dan1Score {
		return Worse
	}
	return Same
}

type LiangDui struct {
	Dui1Score int
	Dui2Score int
	DanScore  int //最小的单
}

//CompareInter 两对的内部比较，对是越大越好，但单牌越小越好
func (this *LiangDui) CompareInter(other *LiangDui) CompareResult {
	if this.Dui1Score > other.Dui1Score {
		return Better
	} else if this.Dui1Score < other.Dui1Score {
		return Worse
	}
	if this.Dui2Score > other.Dui2Score {
		return Better
	} else if this.Dui2Score < other.Dui2Score {
		return Worse
	}
	if this.DanScore < other.DanScore {
		return Better
	} else if this.DanScore > other.DanScore {
		return Worse
	}
	return Same
}

//CompareExternal 两对的外部比较，所有牌是越大越好
func (this *LiangDui) CompareExternal(other *LiangDui) CompareResult {
	if this.Dui1Score > other.Dui1Score {
		return Better
	} else if this.Dui1Score < other.Dui1Score {
		return Worse
	}
	if this.Dui2Score > other.Dui2Score {
		return Better
	} else if this.Dui2Score < other.Dui2Score {
		return Worse
	}
	if this.DanScore > other.DanScore {
		return Better
	} else if this.DanScore < other.DanScore {
		return Worse
	}
	return Same
}

type SanTiao struct {
	SanTiaoScore int //三条的点
	Dan1Score    int //大一点的单张
	Dan2Score    int //小一点的单张
}

//CompareInter 三条的内部比较，三条越大越好，带的两张单牌越小越好
func (this *SanTiao) CompareInter(other *SanTiao) CompareResult {
	if this.SanTiaoScore > other.SanTiaoScore {
		return Better
	} else if this.SanTiaoScore < other.SanTiaoScore {
		return Worse
	}
	if this.Dan1Score < other.Dan1Score {
		return Better
	} else if this.Dan1Score > other.Dan1Score {
		return Worse
	}
	if this.Dan2Score < other.Dan2Score {
		return Better
	} else if this.Dan2Score > other.Dan2Score {
		return Worse
	}
	return Same
}

//CompareExternal 三条外部比较，三条只比三条的牌，
func (this *SanTiao) CompareExternal(other *SanTiao) CompareResult {
	if this.SanTiaoScore > other.SanTiaoScore {
		return Better
	} else if this.SanTiaoScore < other.SanTiaoScore {
		return Worse
	}
	return Same
}

type TieZhi struct {
	TieZhiScore int
	DanScore    int
}

//CompareInter 铁支的只有内部比较，当铁支的点相同时，另一个单牌越小越好
func (this *TieZhi) CompareInter(other *TieZhi) CompareResult {
	if this.TieZhiScore > other.TieZhiScore {
		return Better
	} else if this.TieZhiScore < other.TieZhiScore {
		return Worse
	}
	if this.DanScore < other.DanScore {
		return Better
	} else if this.DanScore > other.DanScore {
		return Worse
	}
	return Same
}

//CompareExternal 铁支的外部比较,外部比较不用比较单牌，因为别人不可能相同的铁支
func (this *TieZhi) CompareExternal(other *TieZhi) CompareResult {
	if this.TieZhiScore > other.TieZhiScore {
		return Better
	} else if this.TieZhiScore < other.TieZhiScore {
		return Worse
	}
	return Same
}

type HuLu struct {
	HuLuScore int
	DuiScore  int
}

//CompareInter 葫芦也只有内部比较，当葫芦的点相同时，另一个对越小越好
func (this *HuLu) CompareInter(other *HuLu) CompareResult {
	if this.HuLuScore > other.HuLuScore {
		return Better
	} else if this.HuLuScore < other.HuLuScore {
		return Worse
	}
	if this.DuiScore < other.DuiScore {
		return Better
	} else if this.DuiScore > other.DuiScore {
		return Worse
	}
	return Same
}

//CompareExternal 葫芦也的外部比较，因为两个不同牌型，三条肯定不一样
func (this *HuLu) CompareExternal(other *HuLu) CompareResult {
	if this.HuLuScore > other.HuLuScore {
		return Better
	} else if this.HuLuScore < other.HuLuScore {
		return Worse
	}
	return Same
}
