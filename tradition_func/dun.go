package tradition_func

import "sort"

type Dun struct {
	Type              NormalType
	pokers            []*Poker
	mapScoreListPoker map[int][]*Poker      //按点分组
	mapHuaListPoker   map[PokerHua][]*Poker //按花分组

	isShunZi    bool  //顺子
	listTieZhi  []int //铁支 int为铁支的分值
	listSanTiao []int //三条 int为三条的分值
	listDui     []int //对子 int为对子的分值
	listDanPai  []int //单牌 int为单牌的分值
}

func NewDun(pokers []*Poker) *Dun {
	SortPoker(pokers)
	dun := &Dun{
		pokers:            pokers,
		mapHuaListPoker:   make(map[PokerHua][]*Poker, 0),
		mapScoreListPoker: make(map[int][]*Poker, 0),
		isShunZi:          false,
	}
	dun.Statistic()
	dun.Type = dun.CalNormalType()
	return dun
}

//Statistic 将牌作统计
func (this *Dun) Statistic() {
	this.statisticsHuaOrScore()
	this.statistics1234()
}

//statisticsHuaOrScore 分组:按花分组或按分值分组
func (this *Dun) statisticsHuaOrScore() {
	for _, poker := range this.pokers {
		if this.mapScoreListPoker[poker.Score] == nil {
			this.mapScoreListPoker[poker.Score] = []*Poker{}
		}
		this.mapScoreListPoker[poker.Score] = append(this.mapScoreListPoker[poker.Score], poker)
		if this.mapHuaListPoker[poker.Hua] == nil {
			this.mapHuaListPoker[poker.Hua] = []*Poker{}
		}
		this.mapHuaListPoker[poker.Hua] = append(this.mapHuaListPoker[poker.Hua], poker)
	}
}

//statistics1234 统计单牌，对子，三条，铁支
func (this *Dun) statistics1234() {

	//统计:单牌，对子，三条，铁支
	for score, pokers := range this.mapScoreListPoker {
		count := len(pokers)
		if count == 1 {
			this.listDanPai = append(this.listDanPai, score)
		} else if count == 2 {
			this.listDui = append(this.listDui, score)
		} else if count == 3 {
			this.listSanTiao = append(this.listSanTiao, score)
		} else if count == 4 {
			this.listTieZhi = append(this.listTieZhi, score)
		}
	}
	sort.Ints(this.listDanPai)
	sort.Ints(this.listDui)
	sort.Ints(this.listSanTiao)
	sort.Ints(this.listTieZhi)
}

func (this *Dun) CalNormalType() NormalType {
	SortPoker(this.pokers)
	isShunZi := func() bool {
		lastScore := 0
		for _, poker := range this.pokers {
			if lastScore == 0 {
				lastScore = poker.Score
			} else if poker.Score == lastScore+1 || (lastScore == 5 && poker.Score == 14) {
				lastScore = poker.Score
			} else if poker.Score != lastScore+1 {
				return false
			}
		}
		return true
	}

	this.isShunZi = isShunZi()

	if this.isShunZi && len(this.mapHuaListPoker) == 1 {
		return TONG_HUA_SHUN
	}

	if len(this.listTieZhi) > 0 {
		return TIE_ZHI
	}

	if len(this.listSanTiao) > 0 && len(this.listDui) > 0 {
		return HU_LU
	}

	if this.isShunZi == false && len(this.mapHuaListPoker) == 1 {
		return TONG_HUA
	}

	if this.isShunZi == true && len(this.mapHuaListPoker) > 1 {
		return SHUN_ZI
	}

	if len(this.listSanTiao) > 0 {
		return SAN_TIAO
	}

	if len(this.listDui) == 2 {
		return LIANG_DUI
	}

	if len(this.listDui) == 1 {
		return DUI_ZI
	}

	return WU_LONG
}

func (this *Dun) Compare(other *Dun) CompareResult {
	if this.Type > other.Type {
		return Better
	}
	if this.Type < other.Type {
		return Worse
	}
	switch this.Type {
	case WU_LONG:
		return ComparePokerScore(this.pokers, other.pokers)
	case TONG_HUA_SHUN:
		return ComparePokerScore(this.pokers, other.pokers)
	case TONG_HUA:
		return ComparePokerScore(this.pokers, other.pokers)
	case SHUN_ZI:
		return ComparePokerScore(this.pokers, other.pokers)
	case DUI_ZI:
		if this.listDui[0] > other.listDui[0] {
			return Better
		} else if this.listDui[0] < other.listDui[0] {
			return Worse
		}
		thisLen := len(this.listDanPai)
		otherLen := len(other.listDanPai)
		if thisLen == 1 && otherLen == 1 {
			if this.listDanPai[thisLen-1] > other.listDanPai[otherLen-1] {
				return Better
			} else if this.listDanPai[thisLen-1] < other.listDanPai[otherLen-1] {
				return Worse
			}
			return Same
		}
		if thisLen == 3 && otherLen == 1 {
			if this.listDanPai[thisLen-1] > other.listDanPai[otherLen-1] {
				return Better
			} else if this.listDanPai[thisLen-1] < other.listDanPai[otherLen-1] {
				return Worse
			}
			return Better
		}
		if thisLen == 1 && otherLen == 3 {
			if this.listDanPai[thisLen-1] > other.listDanPai[otherLen-1] {
				return Better
			} else if this.listDanPai[thisLen-1] < other.listDanPai[otherLen-1] {
				return Worse
			}
			return Worse
		}
		if this.listDanPai[thisLen-1] > other.listDanPai[otherLen-1] {
			return Better
		} else if this.listDanPai[thisLen-1] < other.listDanPai[otherLen-1] {
			return Worse
		}
		if this.listDanPai[thisLen-2] > other.listDanPai[otherLen-2] {
			return Better
		} else if this.listDanPai[thisLen-2] < other.listDanPai[otherLen-2] {
			return Worse
		}

		if this.listDanPai[thisLen-3] > other.listDanPai[otherLen-3] {
			return Better
		} else if this.listDanPai[thisLen-3] < other.listDanPai[otherLen-3] {
			return Worse
		}
		return Same
	case LIANG_DUI:
		if this.listDui[1] > other.listDui[1] {
			return Better
		} else if this.listDui[1] < other.listDui[1] {
			return Worse
		}
		if this.listDui[0] > other.listDui[0] {
			return Better
		} else if this.listDui[0] < other.listDui[0] {
			return Worse
		}
		if this.listDanPai[0] > other.listDanPai[0] {
			return Better
		} else if this.listDanPai[0] < other.listDanPai[0] {
			return Worse
		}
		return Same
	case SAN_TIAO:
		if this.listSanTiao[0] > other.listSanTiao[0] {
			return Better
		} else if this.listSanTiao[0] < other.listSanTiao[0] {
			return Worse
		}
		thisLen := len(this.listDanPai)
		otherLen := len(other.listDanPai)
		if thisLen == 0 && otherLen == 0 {
			return Same
		}
		if thisLen == 2 && otherLen == 0 {
			return Better
		}
		if thisLen == 0 && otherLen == 2 {
			return Worse
		}
		if this.listDanPai[thisLen-1] > other.listDanPai[otherLen-1] {
			return Better
		} else if this.listDanPai[thisLen-1] < other.listDanPai[otherLen-1] {
			return Worse
		}
		if this.listDanPai[thisLen-2] > other.listDanPai[otherLen-2] {
			return Better
		} else if this.listDanPai[thisLen-2] < other.listDanPai[otherLen-2] {
			return Worse
		}
		return Same

	case TIE_ZHI:
		if this.listTieZhi[0] > other.listTieZhi[0] {
			return Better
		} else if this.listTieZhi[0] < other.listTieZhi[0] {
			return Worse
		}
		if this.listDanPai[0] > other.listDanPai[0] {
			return Better
		} else if this.listDanPai[0] < other.listDanPai[0] {
			return Worse
		}
		return Same
	case HU_LU:
		if this.listSanTiao[0] > other.listSanTiao[0] {
			return Better
		} else if this.listSanTiao[0] < other.listSanTiao[0] {
			return Worse
		}
		if this.listDui[0] > other.listDui[0] {
			return Better
		} else if this.listDui[0] < other.listDui[0] {
			return Worse
		}
		return Same
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
