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

//Compare true：一样大或是自己比别人大，false:自己经别人小
func (this *Node) Compare(other *Node) bool {
	if this.normalType == other.normalType {
		switch this.normalType {
		case WU_LONG, TONG_HUA_SHUN, TONG_HUA, SHUN_ZI:
			return ComparePokerScore(this.pokers, other.pokers)
		case DUI_ZI:
			return this.dui.Compare(other.dui)
		case LIANG_DUI:
			return this.liangDui.Compare(other.liangDui)
		case SAN_TIAO:
			return this.sanTiao.Compare(other.sanTiao)
		case TIE_ZHI:
			return this.tieZhi.Compare(other.tieZhi)
		case HU_LU:
			return this.huLu.Compare(other.huLu)
		}
	}
	return this.normalType > other.normalType
}

func ComparePokerScore(pokers1, pokers2 []*Poker) bool {
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
			return poker1.Score > poker2.Score
		}
	}

	if count1 != count2 {
		return count1 > count2
	}

	//平局也是true
	return true
}

type Dui struct {
	duiScore int
	dan1     int //最小的单
	dan2     int //次小的单
	dan3     int //最大的单
}

func (this *Dui) Compare(other *Dui) bool {
	if this.duiScore != other.duiScore {
		return this.duiScore > other.duiScore
	}
	if this.dan3 != other.dan3 {
		return this.dan3 > other.dan3
	}
	if this.dan2 != other.dan2 {
		return this.dan2 > other.dan2
	}
	return this.dan1 > other.dan1
}

type LiangDui struct {
	dui1Score int
	dui2Score int
	dan       int //最小的单
}

func (this *LiangDui) Compare(other *LiangDui) bool {
	if this.dui2Score != other.dui2Score {
		return this.dui2Score > other.dui2Score
	}
	if this.dui1Score != other.dui1Score {
		return this.dui1Score > other.dui1Score
	}
	return this.dan > other.dan
}

type SanTiao struct {
	sanTiaoScore int //三条的点
	bigDanPai    int //大一点的单张
	smallDanPai  int //小一点的单张
}

//Compare 优秀的三条是，三条越大越好，带的两张单牌越小越好
func (this *SanTiao) Compare(other *SanTiao) bool {
	if this.sanTiaoScore == other.sanTiaoScore {
		if this.bigDanPai == other.bigDanPai {
			return this.smallDanPai < other.smallDanPai
		}
		return this.bigDanPai < other.bigDanPai
	}

	return this.sanTiaoScore > other.sanTiaoScore
}

type TieZhi struct {
	tieZhiScore int
	danPaiScore int
}

//Compare 铁支的只有内部比较，当铁支的点相同时，另一个单牌越小越好
func (this *TieZhi) Compare(other *TieZhi) bool {
	if this.tieZhiScore == other.tieZhiScore {
		return this.danPaiScore < other.danPaiScore
	}
	return this.tieZhiScore > other.tieZhiScore
}

type HuLu struct {
	HuLuScore int
	duiScore int
}

//Compare 葫芦也只有内部比较，当葫芦的点相同时，另一个对越小越好
func (this *HuLu) Compare(other *HuLu) bool {
	if this.HuLuScore == other.HuLuScore {
		return this.duiScore < other.duiScore
	}
	return this.HuLuScore > other.HuLuScore
}
