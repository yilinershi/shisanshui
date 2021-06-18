package shi_san_shui

import "fmt"

type ResultNormal struct {
	BestScore int //好牌值
	Left      *Node
	Middle    *Node
	Right     *Node
}

func (this *ResultNormal) String() string {
	leftPokerDesc := ""
	for _, poker := range this.Left.pokers {
		leftPokerDesc += poker.Desc
	}
	middlePokerDesc := ""
	for _, poker := range this.Middle.pokers {
		middlePokerDesc += poker.Desc
	}
	rightPokerDesc := ""
	for _, poker := range this.Right.pokers {
		rightPokerDesc += poker.Desc
	}

	desc := fmt.Sprintf("结果：{左:【%s】= {%s},中:【%s】= {%s},右:【%s】= {%s},好牌值=【%d】}",
		this.Left.normalType, leftPokerDesc, this.Middle.normalType, middlePokerDesc, this.Right.normalType, rightPokerDesc, this.BestScore)
	return desc
}

//CalNormalResults 计算出所有普通牌型
func CalNormalResults(fatherTree *Tree) []*ResultNormal {
	normalInfo := make([]*ResultNormal, 0)
	split(fatherTree)
	for _, node1 := range fatherTree.Nodes {
		//fmt.Println("left=", node1.String())
		sonTree := NewTree(node1.rest)
		split(sonTree)
		for _, node2 := range sonTree.Nodes {
			//fmt.Println("middle=", node2.String())
			if node1.CompareExternal(node2)!=Worse {
				bestScore := 0
				switch node1.normalType {
				case TONG_HUA_SHUN:
					bestScore = bestScore + 5 + int(TONG_HUA_SHUN)
				case TIE_ZHI:
					bestScore = bestScore + 4 + int(TIE_ZHI)
				default:
					bestScore = 1 + int(node1.normalType)
				}

				switch node2.normalType {
				case TONG_HUA_SHUN:
					bestScore = bestScore + 10 + int(TONG_HUA_SHUN)
				case TIE_ZHI:
					bestScore = bestScore + 8 + int(TIE_ZHI)
				case HU_LU:
					bestScore = bestScore + 2 + int(HU_LU)
				default:
					bestScore = 1 + int(node2.normalType)
				}

				node3 := NewNode()
				node3.pokers = node2.rest
				if node3.pokers[0].Score == node3.pokers[1].Score && node3.pokers[1].Score == node3.pokers[2].Score {
					node3.normalType = SAN_TIAO
					node3.sanTiao = &SanTiao{
						SanTiaoScore: node3.pokers[0].Score,
					}
					bestScore = bestScore + 3 + int(node3.normalType)
				} else if node3.pokers[0].Score == node3.pokers[1].Score || node3.pokers[1].Score == node3.pokers[2].Score || node3.pokers[0].Score == node3.pokers[2].Score {
					node3.normalType = DUI_ZI
					node3.dui = &Dui{
						DuiScore: node3.pokers[0].Score,
					}
					bestScore = bestScore + 2 + int(DUI_ZI)
				} else {
					node3.normalType=WU_LONG
					bestScore = bestScore + 1 + int(WU_LONG)
				}

				if node2.CompareExternal(node3)!=Worse {
					result := &ResultNormal{
						Left:      node1,
						Middle:    node2,
						Right:     node3,
						BestScore: bestScore,
					}

					normalInfo = append(normalInfo, result)
				}
			}
		}
	}

	return normalInfo
}


//CalBest 计算出最好的普通牌型
func CalBest(resultList []*ResultNormal) *ResultNormal {
	if len(resultList) <= 0 {
		return nil
	}

	if len(resultList) == 1 {
		return resultList[0]
	}

	var best *ResultNormal = nil
	for _, result := range resultList {
		if best == nil {
			best = result
		} else {
			//case1:按分值比较
			if result.BestScore >= best.BestScore {
				best = result
				//case2:左最优
				if result.Left.CompareInter(best.Left)==Better {
					best = result
					//case3:中最优
					if result.Middle.CompareInter(best.Middle)==Better {
						best = result
						//case4:右最优
						if result.Right.CompareInter(best.Right)==Better {
							best = result
						}
					}
				}
			}
		}
	}
	return best
}

//split 计算普通牌型
func split(tree *Tree) {
	splitTongHuaSun(tree)
	splitTieZhi(tree)
	splitHuLu(tree)

	splitTongHua(tree)
	splitSunZi(tree)
	splitSanTiao(tree)

	splitLiangDui(tree)
	splitDui(tree)
	splitWuLong(tree)
}

//splitTongHuaSun 拆分出同花顺
func splitTongHuaSun(tree *Tree) {

	type tongHuaSunStruct struct {
		hua        PokerHua
		startScore int
		endScore   int
	}
	allTongHuaSun := make([]*tongHuaSunStruct, 0)
	for _, sunZi := range tree._listShunZi {
		for hua, pokers := range tree._mapHuaListPoker {
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
		for _, poker := range tree.pokers {
			if poker.Hua == tongHuaSun.hua && poker.Score >= tongHuaSun.startScore && poker.Score <= tongHuaSun.endScore {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)
	}
}

//splitTieZhi 拆分出铁支
func splitTieZhi(tree *Tree) {
	count := len(tree._listTieZhi)
	if count <= 0 {
		return
	}

	maxTieZhiScore := tree._listTieZhi[count-1]

	for _, poker1 := range tree.pokers {
		if poker1.Score != maxTieZhiScore {
			n := NewNode()
			n.normalType = TIE_ZHI
			n.pokers = append(n.pokers, poker1)
			n.tieZhi = &TieZhi{
				TieZhiScore: maxTieZhiScore,
				DanScore:    poker1.Score,
			}
			for _, poker2 := range tree.pokers {
				if poker1 != poker2 {
					if poker2.Score == maxTieZhiScore {
						n.pokers = append(n.pokers, poker2)
					} else {
						n.rest = append(n.rest, poker2)
					}
				}
			}
			tree.Nodes = append(tree.Nodes, n)
		}
	}

}

//splitHuLu 拆分出葫芦
func splitHuLu(tree *Tree) {
	countDui := len(tree._listDui)
	countSanTiao := len(tree._listSanTiao)
	if countDui < 1 || countSanTiao <= 0 {
		return
	}

	huLuScore := tree._listSanTiao[countSanTiao-1] //取最大的三条

	for _, duiScore := range tree._listDui {

		n := NewNode()
		n.normalType = HU_LU
		n.huLu = &HuLu{
			HuLuScore: huLuScore,
		}
		for _, poker := range tree.pokers {
			if poker.Score == huLuScore || poker.Score == duiScore {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)
	}
}

//splitTongHua 拆分出同花
func splitTongHua(tree *Tree) {

	for _, pokers := range tree._mapHuaListPoker {
		count := len(pokers)
		if count >= 5 {

			for i1 := 0; i1 < count; i1++ {
				for i2 := i1 + 1; i2 < count; i2++ {
					for i3 := i2 + 1; i3 < count; i3++ {
						for i4 := i3 + 1; i4 < count; i4++ {
							for i5 := i4 + 1; i5 < count; i5++ {
								n := NewNode()
								n.normalType = TONG_HUA
								for _, poker := range tree.pokers {
									if poker == pokers[i1] || poker == pokers[i2] || poker == pokers[i3] || poker == pokers[i4] || poker == pokers[i5] {
										n.pokers = append(n.pokers, poker)
									} else {
										n.rest = append(n.rest, poker)
									}
								}
								tree.Nodes = append(tree.Nodes, n)
							}
						}
					}
				}
			}
		}
	}

}

//splitSunZi 顺子
func splitSunZi(tree *Tree) {

	count := len(tree._listShunZi)
	if count <= 0 {
		return
	}

	for _, shunZi := range tree._listShunZi {
		n := NewNode()
		n.normalType = SHUN_ZI

		//顺子中，同样点数的牌可能有多张，要避免这张牌反复加入左节点
		sameScore := 0
		for _, poker := range tree.pokers {
			if poker.Score >= shunZi[0] && poker.Score <= shunZi[1] && poker.Score != sameScore {
				n.pokers = append(n.pokers, poker)
				sameScorePokerCount := len(tree._mapScoreListPoker[poker.Score])
				if sameScorePokerCount > 1 {
					sameScore = poker.Score
				}
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)

	}

	if tree._isHaveSpecialSunZi {
		n := NewNode()
		n.normalType = SHUN_ZI

		//顺子中，同样点数的牌可能有多张，要避免这张牌反复加入左节点
		sameScore := 0
		for _, poker := range tree.pokers {
			if (poker.Point == PokerA || poker.Point == Poker2 || poker.Point == Poker3 || poker.Point == Poker4 || poker.Point == Poker5) && poker.Score != sameScore {
				n.pokers = append(n.pokers, poker)
				sameScorePokerCount := len(tree._mapScoreListPoker[poker.Score])
				if sameScorePokerCount > 1 {
					sameScore = poker.Score
				}
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)
	}
}

//splitSanTiao 拆分出三条
func splitSanTiao(tree *Tree) {

	sanTiaoCount := len(tree._listSanTiao)
	danPaiCount := len(tree._listDanPai)
	if sanTiaoCount <= 0  || danPaiCount < 2 {
		return
	}

	maxSanTiaoScore := tree._listSanTiao[sanTiaoCount-1] //最大的三条分数
	smallDanPaiScore := tree._listDanPai[0]              //最小的单牌
	bigDanPaiScore := tree._listDanPai[1]                //第二小的单牌
	n := NewNode()
	n.normalType = SAN_TIAO
	n.sanTiao = &SanTiao{
		SanTiaoScore: maxSanTiaoScore,
		Dan1Score:    bigDanPaiScore,
		Dan2Score:    smallDanPaiScore,
	}
	for _, poker := range tree.pokers {
		if poker.Score == smallDanPaiScore || poker.Score == bigDanPaiScore || poker.Score == maxSanTiaoScore {
			n.pokers = append(n.pokers, poker)
		} else {
			n.rest = append(n.rest, poker)
		}
	}
	tree.Nodes = append(tree.Nodes, n)
}

//splitLiangDui 两对+单张=两对，这个单张不可能是顺子或是同花中的牌，因为顺子和同花都比两对大，这样两对就没必要出现
//有2对，直接取2对；有3对或4对时，取最小的两对；有5对，取第2大和最小的对
func splitLiangDui(tree *Tree) {
	danPaiCount := len(tree._listDanPai)
	duiCount := len(tree._listDui)
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
		dui1Score = tree._listDui[0]
		dui2Score = tree._listDui[1]
		danPaiScore = tree._listDanPai[0]
	} else if duiCount == 4 {
		if danPaiCount == 0 {  //四对无单牌时，拆最小的对，对为第2小的对和第3小的对
			dui1Score = tree._listDui[2]
			dui2Score = tree._listDui[1]
			danPaiScore = tree._listDui[0]
		} else {//四对有单牌时，对为第1小的对及第2小的对
			dui1Score = tree._listDui[1]
			dui2Score = tree._listDui[0]
			danPaiScore = tree._listDanPai[0]
		}
	}else if duiCount == 5 {
		//5对没有单牌时，说明其它3张能和对里凑出最起码是同花或是顺子，因为顺子和同花比较大，完全可以不考虑二对
		if danPaiCount == 0 {
			return
		} else {
			dui1Score = tree._listDui[3]
			dui2Score = tree._listDui[0]
			danPaiScore = tree._listDanPai[0]
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
	for _, poker := range tree.pokers {
		if poker.Score == dui2Score || poker.Score == dui1Score {
			n.pokers = append(n.pokers, poker)
		} else if poker.Score == danPaiScore && isAddDanPai == false {
			n.pokers = append(n.pokers, poker)
			isAddDanPai = true
		} else {
			n.rest = append(n.rest, poker)
		}
	}
	tree.Nodes = append(tree.Nodes, n)
}

//splitDui 一对+3*单张=两对，这个单张不可能是顺子或是同花中的牌，因为顺子和同花都比一对大，这样一对就没必要出现
func splitDui(tree *Tree) {
	duiCount := len(tree._listDui)
	danPaiCount := len(tree._listDanPai)
	if (duiCount == 1 && danPaiCount >= 3) || (duiCount == 2 && danPaiCount >= 3) || (duiCount == 3 && danPaiCount >= 3) {

		duiScore := 0
		if duiCount == 1 {
			duiScore = tree._listDui[0]
		} else if duiCount == 2 { //有两对时，取大的对
			duiScore = tree._listDui[1]
		} else if duiCount == 3 {
			duiScore = tree._listDui[2]
		}

		dan1Score := tree._listDanPai[0]
		dan2Score := tree._listDanPai[1]
		dan3Score := tree._listDanPai[2]
		n := NewNode()
		n.normalType = DUI_ZI
		n.dui = &Dui{
			DuiScore:  duiScore,
			Dan1Score: dan1Score,
			Dan2Score: dan2Score,
			Dan3Score: dan3Score,
		}
		for _, poker := range tree.pokers {
			if poker.Score == duiScore || poker.Score == dan1Score || poker.Score == dan2Score || poker.Score == dan3Score {
				n.pokers = append(n.pokers, poker)
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)
	}
}

//splitWuLong 拆出一个乌龙
func splitWuLong(tree *Tree) {
	duiCount := len(tree._listDui)
	danPaiCount := len(tree._listDanPai)
	if duiCount > 0 || danPaiCount < 5 {
		return
	}
	dan1Score := tree._listDanPai[danPaiCount-1] //最大的单牌
	dan2Score := tree._listDanPai[0]             //第1小的单牌
	dan3Score := tree._listDanPai[1]             //第2小的单牌
	dan4Score := tree._listDanPai[2]             //第3小的单牌
	dan5Score := tree._listDanPai[3]             //第4小的单牌
	n := NewNode()
	n.normalType = WU_LONG
	for _, poker := range tree.pokers {
		if poker.Score == dan4Score || poker.Score == dan1Score || poker.Score == dan2Score || poker.Score == dan3Score || poker.Score == dan5Score {
			n.pokers = append(n.pokers, poker)

		} else {
			n.rest = append(n.rest, poker)
		}
	}
	tree.Nodes = append(tree.Nodes, n)
}
