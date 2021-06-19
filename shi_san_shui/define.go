package shi_san_shui

import "sort"

type NormalType int

const (
	WU_LONG       NormalType = 0
	DUI_ZI        NormalType = 1
	LIANG_DUI     NormalType = 2
	SAN_TIAO      NormalType = 3
	SHUN_ZI       NormalType = 4
	TONG_HUA      NormalType = 5
	HU_LU         NormalType = 6
	TIE_ZHI       NormalType = 7
	TONG_HUA_SHUN NormalType = 8
)

func (this NormalType) String() string {
	switch this {
	case WU_LONG:
		return "乌龙"
	case DUI_ZI:
		return "对子"
	case LIANG_DUI:
		return "两对"
	case SAN_TIAO:
		return "三条"
	case SHUN_ZI:
		return "顺子"
	case TONG_HUA:
		return "同花"
	case HU_LU:
		return "葫芦"
	case TIE_ZHI:
		return "铁支"
	case TONG_HUA_SHUN:
		return "同花顺"
	default:
		return "乌龙"
	}
}

type SpecialType int

const (
	None SpecialType = iota
	ZhiZunQingLong
	YiTiaoLong
	ShiErHuangZu
	SanTongHuaShun
	SanFenTianXia
	QuanDa
	QuanXiao
	ChouYiSe
	SiTaoSanTiao
	WuDuiSanTiao
	LiuDuiBan
	SanSunZi
	SanTongHua
)

func (this SpecialType) String() string {
	switch this {
	case None:
		return "不是特殊牌型"
	case ZhiZunQingLong:
		return "至尊青龙"
	case YiTiaoLong:
		return "一条龙"
	case ShiErHuangZu:
		return "十二皇族"
	case SanTongHuaShun:
		return "三同花顺"
	case SanFenTianXia:
		return "三分天下"
	case QuanDa:
		return "全大"
	case QuanXiao:
		return "全小"
	case ChouYiSe:
		return "凑一色"
	case SiTaoSanTiao:
		return "四套三条"
	case WuDuiSanTiao:
		return "五对三条"
	case LiuDuiBan:
		return "六对半"
	case SanSunZi:
		return "三顺子"
	case SanTongHua:
		return "三同花"
	default:
		return ""
	}
}

type PokerPoint int

const (
	PokerNone PokerPoint = iota
	PokerA
	Poker2
	Poker3
	Poker4
	Poker5
	Poker6
	Poker7
	Poker8
	Poker9
	PokerT
	PokerJ
	PokerQ
	PokerK
)

func (this PokerPoint) String() string {
	switch this {
	case PokerA:
		return "A"
	case Poker2:
		return "2"
	case Poker3:
		return "3"
	case Poker4:
		return "4"
	case Poker5:
		return "5"
	case Poker6:
		return "6"
	case Poker7:
		return "7"
	case Poker8:
		return "8"
	case Poker9:
		return "9"
	case PokerT:
		return "T"
	case PokerJ:
		return "J"
	case PokerQ:
		return "Q"
	case PokerK:
		return "K"
	default:
		return ""
	}
}

func (this PokerPoint) Score() int {
	switch this {
	case PokerA:
		return 14
	case Poker2:
		return 2
	case Poker3:
		return 3
	case Poker4:
		return 4
	case Poker5:
		return 5
	case Poker6:
		return 6
	case Poker7:
		return 7
	case Poker8:
		return 8
	case Poker9:
		return 9
	case PokerT:
		return 10
	case PokerJ:
		return 11
	case PokerQ:
		return 12
	case PokerK:
		return 13
	default:
		return 0
	}
}

type PokerHua int

const (
	NoneHua PokerHua = iota
	Hua1
	Hua2
	Hua3
	Hua4
)

func (this PokerHua) String() string {
	switch this {
	case Hua1:
		return "♦"
	case Hua2:
		return "♣"
	case Hua3:
		return "♥"
	case Hua4:
		return "♠"
	default:
		return ""
	}
}

func SortPoker(pokers []*Poker) {
	sort.Slice(pokers, func(i, j int) bool {
		if pokers[i].Score == pokers[j].Score {
			return pokers[i].Hua < pokers[j].Hua
		}
		return pokers[i].Score < pokers[j].Score
	})
}


type DescHua string

const (
	NoneHuaDesc DescHua = ""
	Hua1Desc    DescHua = "♦"
	Hua2Desc    DescHua = "♣"
	Hua3Desc    DescHua = "♥"
	Hua4Desc    DescHua = "♠"
)

func (this DescHua) ToPokerHua() PokerHua {
	switch this {
	case Hua1Desc:
		return Hua1
	case Hua2Desc:
		return Hua2
	case Hua3Desc:
		return Hua3
	case Hua4Desc:
		return Hua4
	default:
		return NoneHua
	}
}

type DescPoint string

const (
	DescPointNone DescPoint = ""
	DescPointA    DescPoint = "A"
	DescPoint2    DescPoint = "2"
	DescPoint3    DescPoint = "3"
	DescPoint4    DescPoint = "4"
	DescPoint5    DescPoint = "5"
	DescPoint6    DescPoint = "6"
	DescPoint7    DescPoint = "7"
	DescPoint8    DescPoint = "8"
	DescPoint9    DescPoint = "9"
	DescPointT    DescPoint = "T"
	DescPointJ    DescPoint = "J"
	DescPointQ    DescPoint = "Q"
	DescPointK    DescPoint = "K"
)

func (this DescPoint) ToPokerPoint() PokerPoint {
	switch this {
	case DescPointA:
		return PokerA
	case DescPoint2:
		return Poker2
	case DescPoint3:
		return Poker3
	case DescPoint4:
		return Poker4
	case DescPoint5:
		return Poker5
	case DescPoint6:
		return Poker6
	case DescPoint7:
		return Poker7
	case DescPoint8:
		return Poker8
	case DescPoint9:
		return Poker9
	case DescPointT:
		return PokerT
	case DescPointJ:
		return PokerJ
	case DescPointQ:
		return PokerQ
	case DescPointK:
		return PokerK
	default:
		return PokerNone
	}
}

type CompareResult int

const (
	Same   CompareResult = 0
	Better CompareResult = 1
	Worse  CompareResult = 2
)
