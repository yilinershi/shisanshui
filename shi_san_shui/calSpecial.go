package shi_san_shui

func (this *Tree) CalSpecial() (bool, SpecialType) {
	if this.IsZhiZunQingLong() {
		return true, ZhiZunQingLong
	}
	if this.IsYiTiaoLong() {
		return true, YiTiaoLong
	}
	if this.IsShiErHuangZu() {
		return true, ShiErHuangZu
	}
	if this.IsSanTongHuaShun() {
		return true, SanTongHuaShun
	}
	if this.IsSanFenTianXia() {
		return true, SanFenTianXia
	}
	if this.IsQuanDa() {
		return true, QuanDa
	}
	if this.IsQuanXiao() {
		return true, QuanXiao
	}
	if this.IsChouYiSe() {
		return true, ChouYiSe
	}
	if this.IsSiTaoSanTiao() {
		return true, SiTaoSanTiao
	}
	if this.IsWuDuiSanTiao() {
		return true, WuDuiSanTiao
	}
	if this.IsLiuDuiBan() {
		return true, LiuDuiBan
	}
	if this.IsSanTongHua() {
		return true, SanTongHua
	}
	if this.IsSanSunZi() {
		return true, SanSunZi
	}
	return false, None
}

//IsZhiZunQingLong 是否是"至尊青龙"
func (this *Tree) IsZhiZunQingLong() bool {
	return len(this._mapScoreListPoker) == 13 && len(this._mapHuaListPoker) == 1
}

//IsYiTiaoLong 是否是"一条龙"
func (this *Tree) IsYiTiaoLong() bool {
	return len(this._mapScoreListPoker) == 13 && len(this._mapHuaListPoker) > 1
}

//IsShiErHuangZu 是否是"十二皇族",13张牌中12张牌分值大于等于10
func (this *Tree) IsShiErHuangZu() bool {
	//分少于10的牌的数量
	lessScoreTCount := 0
	for _, poker := range this.pokers {
		if poker.Score < 10 {
			lessScoreTCount++
			if lessScoreTCount >= 2 {
				return false
			}
		}
	}

	return true
}

//IsSanTongHuaShun 是否是"三同花顺"
func (this *Tree) IsSanTongHuaShun() bool {
	//三同花顺，只可能是3个同样的花，也可能是2个花色，但每个花下都是连续的，一定不是一个花色，因为一个花色就是"至尊青龙"
	if len(this._mapHuaListPoker) == 1 || len(this._mapHuaListPoker) == 4 {
		return false
	}

	//同一个花下的所有牌，点数是连续的
	for _, pokers := range this._mapHuaListPoker {
		if len(pokers) == 3 || len(pokers) == 5 || len(pokers) == 8 || len(pokers) == 10 {
			SortPoker(pokers)
			var last *Poker = nil
			for _, poker := range pokers {
				if last == nil {
					last = poker
				} else {
					if poker.Score == last.Score+1 {
						last = poker
					} else {
						return false
					}
				}
			}
		} else {
			return false
		}
	}

	return true
}

//IsSanFenTianXia 是否是"三分天下"
func (this *Tree) IsSanFenTianXia() bool {
	//有3个铁支就是三分天下
	if len(this._listTieZhi) != 3 {
		return false
	}
	return true
}

//IsQuanDa 是否是"全大"
func (this *Tree) IsQuanDa() bool {
	for score, _ := range this._mapScoreListPoker {
		if score < 8 {
			return false
		}
	}
	return true
}

//IsQuanXiao 是否是"全小"
func (this *Tree) IsQuanXiao() bool {
	for score, _ := range this._mapScoreListPoker {
		if score > 8 {
			return false
		}
	}
	return true
}

//IsChouYiSe 即只有一个花色，且花色必需是全红或全黑
func (this *Tree) IsChouYiSe() bool {
	if len(this._mapHuaListPoker) != 2 {
		return false
	}
	var lastHua PokerHua = NoneHua
	for hua, _ := range this._mapHuaListPoker {
		if lastHua == NoneHua {
			lastHua = hua
		} else {
			if hua-lastHua != 0 && hua-lastHua != 2 && hua-lastHua != -2 {
				return false
			}
		}
	}
	return true
}

//IsSiTaoSanTiao 是否四套三条,即AAA,BBB,CCC,DDD,E
func (this *Tree) IsSiTaoSanTiao() bool {
	if len(this._listSanTiao) == 4 {
		return true
	}

	return false
}

//IsWuDuiSanTiao 是否五对三条,即AA、BB、CC、DD、EE、FFF
func (this *Tree) IsWuDuiSanTiao() bool {
	if len(this._listDui) == 5 && len(this._listSanTiao) == 1 {
		return true
	}

	return false
}

//IsLiuDuiBan 是否六对半,即AA、BB、CC、DD、EE、FF、G
func (this *Tree) IsLiuDuiBan() bool {
	if len(this._listDui) == 6 {
		return true
	}

	return false
}

//IsSanTongHua 是否三同花
func (this *Tree) IsSanTongHua() bool {
	//只可能是2种或3种花色，一种花色是同花顺
	if len(this._mapHuaListPoker) == 4 || len(this._mapHuaListPoker) == 1 {
		return false
	}

	//如果是2种花色，其中一花为8张或10，另一个花为5张或3张；如果是3种花色，同一个花下只能是3张或5张
	for _, pokers := range this._mapHuaListPoker {
		count := len(pokers)
		if count != 3 && count != 5 && count != 8 && count != 10 {
			return false
		}
	}

	return true
}

//IsSanSunZi 是否三顺子
func (this *Tree) IsSanSunZi() bool {
	this.splitSunZi()
	for _, node1 := range this.Nodes {
		if node1.normalType == SHUN_ZI {
			middle := NewTree(node1.rest)
			middle.splitSunZi()
			for _, node2 := range middle.Nodes {
				if node2.normalType == SHUN_ZI {
					right := node2.rest
					SortPoker(right)
					if right[0].Score+1 == right[1].Score && right[1].Score+1 == right[2].Score {
						//case1:分值连续
						return true
					} else if right[0].Point == Poker2 && right[1].Point == Poker3 && right[2].Point == PokerA {
						//case2:A,2,3 算
						return true
					}
				}
			}
		}
	}
	return false
}
