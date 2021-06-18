package shi_san_shui

//CalNormal 计算普通牌型
func (this *Tree) CalNormal() {
	this.splitTongHuaSun()
	this.splitTieZhi()
	this.splitHuLu()

	this.splitTongHua()
	this.splitSunZi()
	this.splitSanTiao()

	this.splitLiangDui()
	this.splitDui()
	this.splitWuLong()
}

//splitTongHuaSun 拆分出同花顺
func (this *Tree) splitTongHuaSun() {

	type tongHuaSunStruct struct {
		hua        PokerHua
		startScore int
		endScore   int
	}
	allTongHuaSun := make([]*tongHuaSunStruct, 0)
	for _, sunZi := range this._listShunZi {
		for hua, pokers := range this._mapHuaListPoker {
			if len(pokers) > 5 {
				temp := make([]*Poker, 0)
				for _, poker := range pokers {
					if poker.Hua == hua && poker.Score <= sunZi[1] && poker.Score >= sunZi[0] {
						temp = append(temp, poker)
					}
				}
				if len(temp) == 5 {
					ths := &tongHuaSunStruct{
						hua:        hua,
						startScore: sunZi[0],
						endScore:   sunZi[1],
					}
					allTongHuaSun = append(allTongHuaSun, ths)
				}
			}
		}
	}

	for _, tongHuaSun := range allTongHuaSun {
		n := NewNode()
		n.normalType = TONG_HUA_SHUN
		for _, poker := range this.pokers {
			if poker.Hua == tongHuaSun.hua && poker.Score >= tongHuaSun.startScore && poker.Score <= tongHuaSun.endScore {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		this.Nodes = append(this.Nodes, n)
	}
}

//splitTieZhi 拆分出铁支
func (this *Tree) splitTieZhi() {
	count := len(this._listTieZhi)
	if count <= 0 {
		return
	}

	maxTieZhiScore := this._listTieZhi[count-1]

	for _, poker1 := range this.pokers {
		if poker1.Score != maxTieZhiScore {
			n := NewNode()
			n.normalType = TIE_ZHI
			n.pokers = append(n.pokers, poker1)
			n.tieZhi = &TieZhi{
				TieZhiScore: maxTieZhiScore,
				DanScore:    poker1.Score,
			}
			for _, poker2 := range this.pokers {
				if poker1 != poker2 {
					if poker2.Score == maxTieZhiScore {
						n.pokers = append(n.pokers, poker2)
					} else {
						n.rest = append(n.rest, poker2)
					}
				}
			}
			this.Nodes = append(this.Nodes, n)
		}
	}

}

//splitHuLu 拆分出葫芦
func (this *Tree) splitHuLu() {
	countDui := len(this._listDui)
	countSanTiao := len(this._listSanTiao)
	if countDui < 1 || countSanTiao <= 0 {
		return
	}

	huLuScore := this._listSanTiao[countSanTiao-1] //取最大的三条

	for _, duiScore := range this._listDui {

		n := NewNode()
		n.normalType = HU_LU
		n.huLu = &HuLu{
			HuLuScore: huLuScore,
		}
		for _, poker := range this.pokers {
			if poker.Score == huLuScore || poker.Score == duiScore {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		this.Nodes = append(this.Nodes, n)
	}
}

//splitTongHua 拆分出同花
func (this *Tree) splitTongHua() {

	for _, pokers := range this._mapHuaListPoker {
		count := len(pokers)
		if count >= 5 {

			for i1 := 0; i1 < count; i1++ {
				for i2 := i1 + 1; i2 < count; i2++ {
					for i3 := i2 + 1; i3 < count; i3++ {
						for i4 := i3 + 1; i4 < count; i4++ {
							for i5 := i4 + 1; i5 < count; i5++ {
								n := NewNode()
								n.normalType = TONG_HUA
								for _, poker := range this.pokers {
									if poker == pokers[i1] || poker == pokers[i2] || poker == pokers[i3] || poker == pokers[i4] || poker == pokers[i5] {
										n.pokers = append(n.pokers, poker)
									} else {
										n.rest = append(n.rest, poker)
									}
								}
								this.Nodes = append(this.Nodes, n)
							}
						}
					}
				}
			}
		}
	}

}

//splitSunZi 顺子
func (this *Tree) splitSunZi() {

	count := len(this._listShunZi)
	if count <= 0 {
		return
	}

	for _, shunZi := range this._listShunZi {
		n := NewNode()
		n.normalType = SHUN_ZI

		//顺子中，同样点数的牌可能有多张，要避免这张牌反复加入左节点
		sameScore := 0
		for _, poker := range this.pokers {
			if poker.Score >= shunZi[0] && poker.Score <= shunZi[1] && poker.Score != sameScore {
				n.pokers = append(n.pokers, poker)
				sameScorePokerCount := len(this._mapScoreListPoker[poker.Score])
				if sameScorePokerCount > 1 {
					sameScore = poker.Score
				}
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		this.Nodes = append(this.Nodes, n)

	}

	if this._isHaveSpecialSunZi {
		n := NewNode()
		n.normalType = SHUN_ZI

		//顺子中，同样点数的牌可能有多张，要避免这张牌反复加入左节点
		sameScore := 0
		for _, poker := range this.pokers {
			if (poker.Point == PokerA || poker.Point == Poker2 || poker.Point == Poker3 || poker.Point == Poker4 || poker.Point == Poker5) && poker.Score != sameScore {
				n.pokers = append(n.pokers, poker)
				sameScorePokerCount := len(this._mapScoreListPoker[poker.Score])
				if sameScorePokerCount > 1 {
					sameScore = poker.Score
				}
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		this.Nodes = append(this.Nodes, n)
	}
}

//splitSanTiao 拆分出三条
func (this *Tree) splitSanTiao() {

	sanTiaoCount := len(this._listSanTiao)
	danPaiCount := len(this._listDanPai)
	if sanTiaoCount <= 0  || danPaiCount < 2 {
		return
	}

	maxSanTiaoScore := this._listSanTiao[sanTiaoCount-1] //最大的三条分数
	smallDanPaiScore := this._listDanPai[0]              //最小的单牌
	bigDanPaiScore := this._listDanPai[1]                //第二小的单牌
	n := NewNode()
	n.normalType = SAN_TIAO
	n.sanTiao = &SanTiao{
		SanTiaoScore: maxSanTiaoScore,
		Dan1Score:    bigDanPaiScore,
		Dan2Score:    smallDanPaiScore,
	}
	for _, poker := range this.pokers {
		if poker.Score == smallDanPaiScore || poker.Score == bigDanPaiScore || poker.Score == maxSanTiaoScore {
			n.pokers = append(n.pokers, poker)
		} else {
			n.rest = append(n.rest, poker)
		}
	}
	this.Nodes = append(this.Nodes, n)
}

//splitLiangDui 两对+单张=两对，这个单张不可能是顺子或是同花中的牌，因为顺子和同花都比两对大，这样两对就没必要出现
//有2对，直接取2对；有3对或4对时，取最小的两对；有5对，取第2大和最小的对
func (this *Tree) splitLiangDui() {
	danPaiCount := len(this._listDanPai)
	duiCount := len(this._listDui)
	if duiCount < 2 {
		return
	}
	danPaiScore := 0
	dui1Score := 0
	dui2Score := 0

	if duiCount == 2 || duiCount == 3 {
		//2对或3对没有单牌时，说明其它3张能和对里凑出最起码是同花或是顺子，因为顺子和同花比较大，完全可以不考虑二对
		if danPaiCount == 0 {
			return
		}
		dui1Score = this._listDui[0]
		dui2Score = this._listDui[1]
		danPaiScore = this._listDanPai[0]
	} else if duiCount == 4 {
		if danPaiCount == 0 {  //四对无单牌时，拆最小的对，对为第2小的对和第3小的对
			dui1Score = this._listDui[2]
			dui2Score = this._listDui[1]
			danPaiScore = this._listDui[0]
		} else {//四对有单牌时，对为第1小的对及第2小的对
			dui1Score = this._listDui[1]
			dui2Score = this._listDui[0]
			danPaiScore = this._listDanPai[0]
		}
	}else if duiCount == 5 {
		//5对没有单牌时，说明其它3张能和对里凑出最起码是同花或是顺子，因为顺子和同花比较大，完全可以不考虑二对
		if danPaiCount == 0 {
			return
		} else {
			dui1Score = this._listDui[3]
			dui2Score = this._listDui[0]
			danPaiScore = this._listDanPai[0]
		}
	}

	n := NewNode()
	n.normalType = LIANG_DUI
	n.liangDui = &LiangDui{
		Dui1Score: dui1Score,
		Dui2Score: dui2Score,
		DanScore:  danPaiScore,
	}
	isAddDanPai := false
	for _, poker := range this.pokers {
		if poker.Score == dui2Score || poker.Score == dui1Score {
			n.pokers = append(n.pokers, poker)
		} else if poker.Score == danPaiScore && isAddDanPai == false {
			n.pokers = append(n.pokers, poker)
			isAddDanPai = true
		} else {
			n.rest = append(n.rest, poker)
		}
	}
	this.Nodes = append(this.Nodes, n)
}

//splitDui 一对+3*单张=两对，这个单张不可能是顺子或是同花中的牌，因为顺子和同花都比一对大，这样一对就没必要出现
func (this *Tree) splitDui() {
	duiCount := len(this._listDui)
	danPaiCount := len(this._listDanPai)
	if (duiCount == 1 && danPaiCount >= 3) || (duiCount == 2 && danPaiCount >= 3) || (duiCount == 3 && danPaiCount >= 3) {

		duiScore := 0
		if duiCount == 1 {
			duiScore = this._listDui[0]
		} else if duiCount == 2 { //有两对时，取大的对
			duiScore = this._listDui[1]
		} else if duiCount == 3 {
			duiScore = this._listDui[2]
		}

		dan1Score := this._listDanPai[0]
		dan2Score := this._listDanPai[1]
		dan3Score := this._listDanPai[2]
		n := NewNode()
		n.normalType = DUI_ZI
		n.dui = &Dui{
			DuiScore:  duiScore,
			Dan1Score: dan1Score,
			Dan2Score: dan2Score,
			Dan3Score: dan3Score,
		}
		for _, poker := range this.pokers {
			if poker.Score == duiScore || poker.Score == dan1Score || poker.Score == dan2Score || poker.Score == dan3Score {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		this.Nodes = append(this.Nodes, n)
	}
}

//splitWuLong 拆出一个乌龙
func (this *Tree) splitWuLong() {
	duiCount := len(this._listDui)
	danPaiCount := len(this._listDanPai)
	if duiCount > 0 || danPaiCount < 5 {
		return
	}
	dan1Score := this._listDanPai[danPaiCount-1] //最大的单牌
	dan2Score := this._listDanPai[0]             //第1小的单牌
	dan3Score := this._listDanPai[1]             //第2小的单牌
	dan4Score := this._listDanPai[2]             //第3小的单牌
	dan5Score := this._listDanPai[3]             //第4小的单牌
	n := NewNode()
	n.normalType = WU_LONG
	for _, poker := range this.pokers {
		if poker.Score == dan4Score || poker.Score == dan1Score || poker.Score == dan2Score || poker.Score == dan3Score || poker.Score == dan5Score {
			n.pokers = append(n.pokers, poker)

		} else {
			n.rest = append(n.rest, poker)
		}
	}
	this.Nodes = append(this.Nodes, n)
}
