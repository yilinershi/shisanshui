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

//CompareInter true：一样大或是自己比别人大，false:自己经别人小
func (this *Node) CompareInter(other *Node) CompareResult {

	if this.normalType > other.normalType{
		return Better
	}else if this.normalType < other.normalType{
		return Worse
	}else {
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
}

//CompareInter true：一样大或是自己比别人大，false:自己经别人小
func (this *Node) CompareExternal(other *Node) CompareResult {

	if this.normalType > other.normalType{
		return Better
	}else if this.normalType < other.normalType{
		return Worse
	}else {
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
		if poker1.Score != poker2.Score {
			if poker1.Score > poker2.Score {
				return Better
			} else {
				return Worse
			}
		}
	}

	if count1 != count2 {
		if count1 > count2 {
			return Better
		} else {
			return Worse
		}
	}

	//平局也是true
	return Same
}

type Dui struct {
	duiScore int //对的牌的
	dan1     int //最小的单
	dan2     int //次小的单
	dan3     int //最大的单
}

//CompareInter 对的内部比较，所有牌越大越好
func (this *Dui) CompareInter(other *Dui) CompareResult {
	if this.duiScore > other.duiScore {
		return Better
	} else if this.duiScore < other.duiScore {
		return Worse
	}

	if this.dan3 < other.dan3 {
		return Better
	} else if this.dan3 > other.dan3 {
		return Worse
	}

	if this.dan2 < other.dan2 {
		return Better
	} else if this.dan2 > other.dan2 {
		return Worse
	}

	if this.dan1 < other.dan1 {
		return Better
	} else if this.dan1 > other.dan1 {
		return Worse
	}
	return Same
}

//CompareExternal 对的外部比较，所有牌越大越好
func (this *Dui) CompareExternal(other *Dui) CompareResult {
	if this.duiScore > other.duiScore {
		return Better
	} else if this.duiScore < other.duiScore {
		return Worse
	}

	if this.dan3 > other.dan3 {
		return Better
	} else if this.dan3 < other.dan3 {
		return Worse
	}

	if this.dan2 > other.dan2 {
		return Better
	} else if this.dan2 < other.dan2 {
		return Worse
	}

	if this.dan1 > other.dan1 {
		return Better
	} else if this.dan1 < other.dan1 {
		return Worse
	}
	return Same
}

type LiangDui struct {
	dui1Score int
	dui2Score int
	dan       int //最小的单
}

//CompareInter 两对的内部比较，对是越大越好，但单牌越小越好
func (this *LiangDui) CompareInter(other *LiangDui) CompareResult {
	if this.dui1Score > other.dui1Score {
		return Better
	} else if this.dui1Score < other.dui1Score {
		return Worse
	}
	if this.dui2Score > other.dui2Score {
		return Better
	} else if this.dui2Score < other.dui2Score {
		return Worse
	}
	if this.dan < other.dan {
		return Better
	} else if this.dan > other.dan {
		return Worse
	}
	return Same
}

//CompareExternal 两对的外部比较，所有牌是越大越好
func (this *LiangDui) CompareExternal(other *LiangDui) CompareResult {
	if this.dui1Score > other.dui1Score {
		return Better
	} else if this.dui1Score < other.dui1Score {
		return Worse
	}
	if this.dui2Score > other.dui2Score {
		return Better
	} else if this.dui2Score < other.dui2Score {
		return Worse
	}
	if this.dan > other.dan {
		return Better
	} else if this.dan < other.dan {
		return Worse
	}
	return Same
}

type SanTiao struct {
	sanTiaoScore int //三条的点
	bigDanPai    int //大一点的单张
	smallDanPai  int //小一点的单张
}

//CompareInter 三条的内部比较，三条越大越好，带的两张单牌越小越好
func (this *SanTiao) CompareInter(other *SanTiao) CompareResult {
	if this.sanTiaoScore > other.sanTiaoScore {
		return Better
	} else if this.sanTiaoScore < other.sanTiaoScore {
		return Worse
	}
	if this.bigDanPai < other.bigDanPai {
		return Better
	} else if this.bigDanPai > other.bigDanPai {
		return Worse
	}
	if this.smallDanPai < other.smallDanPai {
		return Better
	} else if this.smallDanPai > other.smallDanPai {
		return Worse
	}
	return Same
}

//CompareExternal 三条外部比较，三条只比三条的牌，
func (this *SanTiao) CompareExternal(other *SanTiao) CompareResult {
	if this.sanTiaoScore > other.sanTiaoScore {
		return Better
	} else if this.sanTiaoScore < other.sanTiaoScore {
		return Worse
	}
	return Same
}

type TieZhi struct {
	tieZhiScore int
	danPaiScore int
}

//CompareInter 铁支的只有内部比较，当铁支的点相同时，另一个单牌越小越好
func (this *TieZhi) CompareInter(other *TieZhi) CompareResult {
	if this.tieZhiScore > other.tieZhiScore {
		return Better
	} else if this.tieZhiScore < other.tieZhiScore {
		return Worse
	}
	if this.danPaiScore < other.danPaiScore {
		return Better
	} else if this.danPaiScore > other.danPaiScore {
		return Worse
	}
	return Same
}

//CompareExternal 铁支的外部比较,外部比较不用比较单牌，因为别人不可能相同的铁支
func (this *TieZhi) CompareExternal(other *TieZhi) CompareResult {
	if this.tieZhiScore > other.tieZhiScore {
		return Better
	} else if this.tieZhiScore < other.tieZhiScore {
		return Worse
	}
	return Same
}

type HuLu struct {
	HuLuScore int
	duiScore  int
}

//CompareInter 葫芦也只有内部比较，当葫芦的点相同时，另一个对越小越好
func (this *HuLu) CompareInter(other *HuLu) CompareResult {
	if this.HuLuScore > other.HuLuScore {
		return Better
	} else if this.HuLuScore < other.HuLuScore {
		return Worse
	}
	if this.duiScore < other.duiScore {
		return Better
	} else if this.duiScore > other.duiScore {
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
