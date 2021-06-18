package shi_san_shui

import "sort"

type Tree struct {
	pokers              []*Poker              //用来分节点的牌
	_mapScoreListPoker  map[int][]*Poker      //按点分组
	_mapHuaListPoker    map[PokerHua][]*Poker //按花分组
	_listShunZi         [][2]int              //顺子
	_isHaveSpecialSunZi bool                  //是否有特殊的顺子，如A2345
	_listTongHua        []PokerHua            //同花
	_mapTongHuaSun      map[PokerHua][2]int   //同花顺
	_listTieZhi         []int                 //铁支 int为铁支的分值
	_listSanTiao        []int                 //三条 int为三条的分值
	_listDui            []int                 //对子 int为对子的分值
	_listDanPai         []int                 //单牌 int为单牌的分值
	Nodes               []*Node               //二叉树的节点
}

func NewTree(pokers []*Poker) *Tree {
	SortPoker(pokers)
	tree := &Tree{
		pokers:              pokers,
		_mapHuaListPoker:    make(map[PokerHua][]*Poker, 0),
		_mapScoreListPoker:  make(map[int][]*Poker, 0),
		_mapTongHuaSun:      make(map[PokerHua][2]int, 0),
		Nodes:               make([]*Node, 0),
		_isHaveSpecialSunZi: false,
	}
	tree.Statistic()
	return tree
}

//Statistic 将牌作统计
func (this *Tree) Statistic() {
	this.statisticsHuaOrScore()
	this.statisticsSunZi()
	this.statistics1234()
	this.statisticsSpecialSunZi()
}

//statisticsHuaOrScore 分组:按花分组或按分值分组
func (this *Tree) statisticsHuaOrScore() {
	for _, poker := range this.pokers {
		if this._mapScoreListPoker[poker.Score] == nil {
			this._mapScoreListPoker[poker.Score] = []*Poker{}
		}
		this._mapScoreListPoker[poker.Score] = append(this._mapScoreListPoker[poker.Score], poker)
		if this._mapHuaListPoker[poker.Hua] == nil {
			this._mapHuaListPoker[poker.Hua] = []*Poker{}
		}
		this._mapHuaListPoker[poker.Hua] = append(this._mapHuaListPoker[poker.Hua], poker)
	}
}

//statisticsSunZi 统计顺子，从顺a到顺b
func (this *Tree) statisticsSunZi() {
	_tempStart := 0
	_tempEnd := 0
	_tempCount := 0
	for score := 2; score <= 14; score++ {
		if _, ok := this._mapScoreListPoker[score]; ok {
			if _tempStart == 0 {
				_tempStart = score //从哪开始连接，比如从2分开始连续
				_tempCount = 1     //共连续的长度
			} else {
				_tempCount = _tempCount + 1
			}
			if _tempCount >= 5 {
				_tempEnd = score
				this._listShunZi = append(this._listShunZi, [2]int{_tempEnd - 4, _tempEnd})
			}
		} else {
			_tempStart = 0
			_tempCount = 1
		}


	}
}

func (this *Tree) statisticsSpecialSunZi() {
	specialSunZiSocre := [5]int{14, 2, 3, 4, 5}
	for _, score := range specialSunZiSocre {
		if _, ok := this._mapScoreListPoker[score]; !ok {
			this._isHaveSpecialSunZi = false
			return
		}
	}
	this._isHaveSpecialSunZi = true
}

//statistics1234 统计单牌，对子，三条，铁支
func (this *Tree) statistics1234() {
	var isDanPai = func(poker *Poker) bool {
		//是否在顺子里
		for _, shunZi := range this._listShunZi {
			if poker.Score >= shunZi[0] && poker.Score <= shunZi[1] {
				return false
			}
		}

		//该牌是否在同花里
		for _, pokers := range this._mapHuaListPoker {
			if len(pokers)>=5{
				for _, p := range pokers {
					if p==poker{
						return false
					}
				}
			}
		}

		return true
	}

	//统计:单牌，对子，三条，铁支
	for score, pokers := range this._mapScoreListPoker {
		count := len(pokers)
		if count == 1 {
			if isDanPai(pokers[0]) {
				this._listDanPai = append(this._listDanPai, score)
			}
		} else if count == 2 {
			this._listDui = append(this._listDui, score)
		} else if count == 3 {
			this._listSanTiao = append(this._listSanTiao, score)
		} else if count == 4 {
			this._listTieZhi = append(this._listTieZhi, score)
		}
	}
	sort.Ints(this._listDanPai)
	sort.Ints(this._listDui)
	sort.Ints(this._listSanTiao)
	sort.Ints(this._listTieZhi)
}
