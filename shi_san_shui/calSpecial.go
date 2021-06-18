package shi_san_shui

func CalSpecial(tree *Tree) (bool, SpecialType) {
	if IsZhiZunQingLong(tree) {
		return true, ZhiZunQingLong
	}
	if IsYiTiaoLong(tree) {
		return true, YiTiaoLong
	}
	if IsShiErHuangZu(tree) {
		return true, ShiErHuangZu
	}
	if IsSanTongHuaShun(tree) {
		return true, SanTongHuaShun
	}
	if IsSanFenTianXia(tree) {
		return true, SanFenTianXia
	}
	if IsQuanDa(tree) {
		return true, QuanDa
	}
	if IsQuanXiao(tree) {
		return true, QuanXiao
	}
	if IsCouYiSe(tree) {
		return true, ChouYiSe
	}
	if IsSiTaoSanTiao(tree) {
		return true, SiTaoSanTiao
	}
	if IsWuDuiSanTiao(tree) {
		return true, WuDuiSanTiao
	}
	if IsLiuDuiBan(tree) {
		return true, LiuDuiBan
	}
	if IsSanTongHua(tree) {
		return true, SanTongHua
	}
	if IsSanSunZi(tree) {
		return true, SanSunZi
	}
	return false, None
}

//IsZhiZunQingLong 是否是"至尊青龙"
func IsZhiZunQingLong(tree *Tree) bool {
	return len(tree._mapScoreListPoker) == 13 && len(tree._mapHuaListPoker) == 1
}

//IsYiTiaoLong 是否是"一条龙"
func IsYiTiaoLong(tree *Tree) bool {
	return len(tree._mapScoreListPoker) == 13 && len(tree._mapHuaListPoker) > 1
}

//IsShiErHuangZu 是否是"十二皇族",13张牌中12张牌分值大于等于10
func IsShiErHuangZu(tree *Tree) bool {
	//分少于10的牌的数量
	lessScoreTCount := 0
	for _, poker := range tree.pokers {
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
func IsSanTongHuaShun(tree *Tree) bool {
	//三同花顺，只可能是3个同样的花，也可能是2个花色，但每个花下都是连续的，一定不是一个花色，因为一个花色就是"至尊青龙"
	if len(tree._mapHuaListPoker) == 1 || len(tree._mapHuaListPoker) == 4 {
		return false
	}

	//同一个花下的所有牌，点数是连续的
	for _, pokers := range tree._mapHuaListPoker {
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
func IsSanFenTianXia(tree *Tree) bool {
	//有3个铁支就是三分天下
	if len(tree._listTieZhi) != 3 {
		return false
	}
	return true
}

//IsQuanDa 是否是"全大"
func IsQuanDa(tree *Tree) bool {
	for score, _ := range tree._mapScoreListPoker {
		if score < 8 {
			return false
		}
	}
	return true
}

//IsQuanXiao 是否是"全小"
func IsQuanXiao(tree *Tree) bool {
	for score, _ := range tree._mapScoreListPoker {
		if score > 8 {
			return false
		}
	}
	return true
}

//IsCouYiSe 即只有一个花色，且花色必需是全红或全黑
func IsCouYiSe(tree *Tree) bool {
	if len(tree._mapHuaListPoker) != 2 {
		return false
	}
	var lastHua PokerHua = NoneHua
	for hua, _ := range tree._mapHuaListPoker {
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
func IsSiTaoSanTiao(tree *Tree) bool {
	if len(tree._listSanTiao) == 4 {
		return true
	}

	return false
}

//IsWuDuiSanTiao 是否五对三条,即AA、BB、CC、DD、EE、FFF
func IsWuDuiSanTiao(tree *Tree) bool {
	if len(tree._listDui) == 5 && len(tree._listSanTiao) == 1 {
		return true
	}

	return false
}

//IsLiuDuiBan 是否六对半,即AA、BB、CC、DD、EE、FF、G
func IsLiuDuiBan(tree *Tree) bool {
	if len(tree._listDui) == 6 {
		return true
	}

	return false
}

//IsSanTongHua 是否三同花
func IsSanTongHua(tree *Tree) bool {
	//只可能是2种或3种花色，一种花色是同花顺
	if len(tree._mapHuaListPoker) == 4 || len(tree._mapHuaListPoker) == 1 {
		return false
	}

	//如果是2种花色，其中一花为8张或10，另一个花为5张或3张；如果是3种花色，同一个花下只能是3张或5张
	for _, pokers := range tree._mapHuaListPoker {
		count := len(pokers)
		if count != 3 && count != 5 && count != 8 && count != 10 {
			return false
		}
	}

	return true
}

//IsSanSunZi 是否三顺子
func IsSanSunZi(fatherTree *Tree) bool {
	splitSunZi(fatherTree)
	for _, node1 := range fatherTree.Nodes {
		if node1.normalType == SHUN_ZI {
			sonTree := NewTree(node1.rest)
			splitSunZi(sonTree)
			for _, node2 := range sonTree.Nodes {
				if node2.normalType == SHUN_ZI {
					node3 := node2.rest
					SortPoker(node3)
					if node3[0].Score+1 == node3[1].Score && node3[1].Score+1 == node3[2].Score {
						//case1:分值连续
						return true
					} else if node3[0].Point == Poker2 && node3[1].Point == Poker3 && node3[2].Point == PokerA {
						//case2:A,2,3 算
						return true
					}
				}
			}
		}
	}
	return false
}
