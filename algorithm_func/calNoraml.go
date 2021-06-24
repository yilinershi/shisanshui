package algorithm_func

import (
	"fmt"
	"sort"
)

type ResultNormal struct {
	BestScore int //好牌值
	Head      *Node
	Middle    *Node
	Tail      *Node
}

func (this *ResultNormal) String() string {
	tailPokerDesc := ""
	for _, poker := range this.Tail.pokers {
		tailPokerDesc += poker.Desc
	}
	middlePokerDesc := ""
	for _, poker := range this.Middle.pokers {
		middlePokerDesc += poker.Desc
	}
	headPokerDesc := ""
	for _, poker := range this.Head.pokers {
		headPokerDesc += poker.Desc
	}

	return fmt.Sprintf("{上墩:【%s】= {%s},中墩:【%s】= {%s},下墩:【%s】= {%s},好牌值=【%d】}",
		this.Head.normalType, headPokerDesc, this.Middle.normalType, middlePokerDesc, this.Tail.normalType, tailPokerDesc, this.BestScore)
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
			if node1.CompareExternal(node2) != Worse {
				bestScore := 0
				node3 := CalCardType(node2.rest)
				if node2.CompareExternal(node3) != Worse {

					switch node1.normalType {
					case TONG_HUA_SHUN:
						bestScore += 5 + int(node1.normalType)
					case TIE_ZHI:
						bestScore += 4 + int(node1.normalType)
					default:
						bestScore += 1 + int(node1.normalType)
					}
					switch node2.normalType {
					case TONG_HUA_SHUN:
						bestScore += 10 + int(node2.normalType)
					case TIE_ZHI:
						bestScore += 8 + int(node2.normalType)
					case HU_LU:
						bestScore += 2 + int(node2.normalType)
					default:
						bestScore += 1 + int(node2.normalType)
					}
					switch node3.normalType {
					case SAN_TIAO:
						bestScore += 3 + int(node3.normalType)
					case DUI_ZI:
						bestScore += 2 + int(node3.normalType)
					default:
						bestScore += 1 + int(node3.normalType)
					}

					result := &ResultNormal{
						Tail:      node1,
						Middle:    node2,
						Head:      node3,
						BestScore: bestScore,
					}
					normalInfo = append(normalInfo, result)
				}
			}
		}
	}

	return normalInfo
}

//SortFilterResult 计算出值得推荐的普通牌型，step1:优先保证分值，其次保证上墩，再次保证中墩，最后保证下墩，step2:过滤掉重复的类型
func SortFilterResult(resultList []*ResultNormal) []*ResultNormal {
	sort.Slice(resultList, func(i, j int) bool {
		if resultList[i].BestScore > resultList[j].BestScore {
			return true
		} else if resultList[i].BestScore < resultList[j].BestScore {
			return false
		}
		if resultList[i].Head.CompareExternal(resultList[j].Head) == Better {
			return true
		} else if resultList[i].Head.CompareExternal(resultList[j].Head) == Worse {
			return false
		}
		if resultList[i].Tail.CompareExternal(resultList[j].Tail) == Better {
			return true
		} else if resultList[i].Tail.CompareExternal(resultList[j].Tail) == Worse {
			return false
		}
		if resultList[i].Middle.CompareExternal(resultList[j].Middle) == Better {
			return true
		} else if resultList[i].Middle.CompareExternal(resultList[j].Middle) == Worse {
			return false
		}
		return false
	})

	var filterRes = make([]*ResultNormal, 0)
	var last *ResultNormal = nil
	for _, result := range resultList {
		if last == nil {
			last = result
			filterRes = append(filterRes, last)
		} else if result.Head.normalType != last.Head.normalType || result.Middle.normalType != last.Middle.normalType || result.Tail.normalType != last.Tail.normalType {
			last = result
			filterRes = append(filterRes, result)
		}
	}
	return filterRes
}

//split 将树节点按牌型拆分
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

	for _, pokers := range tree.mapHuaListPoker {
		count := len(pokers)
		if count < 5 {
			continue
		}

		for i1 := 0; i1 < count-4; i1++ {
			for i2 := i1 + 1; i2 < count-3; i2++ {
				if pokers[i1].Score+1 == pokers[i2].Score {
					for i3 := i2 + 1; i3 < count-2; i3++ {
						if pokers[i2].Score+1 == pokers[i3].Score {
							for i4 := i3 + 1; i4 < count-1; i4++ {
								if pokers[i3].Score+1 == pokers[i4].Score {
									for i5 := i4 + 1; i5 < count; i5++ {
										if pokers[i4].Score+1 == pokers[i5].Score || (pokers[i4].Point == Poker5 && pokers[i5].Point == PokerA) {
											n := NewNode()
											n.normalType = TONG_HUA_SHUN
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
		}
	}
}

//splitTieZhi 拆分出铁支
func splitTieZhi(tree *Tree) {
	count := len(tree.listTieZhi)
	if count <= 0 {
		return
	}

	maxTieZhiScore := tree.listTieZhi[count-1]

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

//splitHuLu 拆分出葫芦,葫芦是最大的三条和任意对组合，有可能最小的对能拆成同花或顺子
func splitHuLu(tree *Tree) {
	countDui := len(tree.listDui)
	countSanTiao := len(tree.listSanTiao)
	if countDui < 1 || countSanTiao <= 0 {
		return
	}

	huLuScore := tree.listSanTiao[countSanTiao-1] //取最大的三条

	for _, duiScore := range tree.listDui {

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
	for _, pokers := range tree.mapHuaListPoker {
		count := len(pokers)
		if count < 5 {
			continue
		}

		for i1 := 0; i1 < count-4; i1++ {
			for i2 := i1 + 1; i2 < count-3; i2++ {
				for i3 := i2 + 1; i3 < count-2; i3++ {
					for i4 := i3 + 1; i4 < count-1; i4++ {
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

//splitSunZi 顺子
func splitSunZi(tree *Tree) {

	count := len(tree.listShunZi)
	if count <= 0 {
		return
	}

	for _, shunZi := range tree.listShunZi {
		poker1s := tree.mapScoreListPoker[shunZi[0]]
		poker2s := tree.mapScoreListPoker[shunZi[0]+1]
		poker3s := tree.mapScoreListPoker[shunZi[0]+2]
		poker4s := tree.mapScoreListPoker[shunZi[0]+3]
		poker5s := tree.mapScoreListPoker[shunZi[1]]

		for i1 := 0; i1 < len(poker1s); i1++ {
			for i2 := 0; i2 < len(poker2s); i2++ {
				for i3 := 0; i3 < len(poker3s); i3++ {
					for i4 := 0; i4 < len(poker4s); i4++ {
						for i5 := 0; i5 < len(poker5s); i5++ {
							n := NewNode()
							n.normalType = SHUN_ZI
							for _, poker := range tree.pokers {
								if poker == poker1s[i1] || poker == poker2s[i2] || poker == poker3s[i3] || poker == poker4s[i4] || poker == poker5s[i5] {
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

//splitSanTiao 拆分出三条
func splitSanTiao(tree *Tree) {

	sanTiaoCount := len(tree.listSanTiao)
	danPaiCount := len(tree.listDanPai)
	if sanTiaoCount <= 0 || danPaiCount < 2 {
		return
	}

	maxSanTiaoScore := tree.listSanTiao[sanTiaoCount-1] //最大的三条分数
	smallDanPaiScore := tree.listDanPai[0]              //最小的单牌
	bigDanPaiScore := tree.listDanPai[1]                //第二小的单牌
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
	danPaiCount := len(tree.listDanPai)
	duiCount := len(tree.listDui)
	if duiCount < 2 {
		return
	}

	var listLiangDui = make([]*LiangDui, 0)

	if duiCount == 2 || duiCount == 3 {
		//2对或3对没有单牌时，说明其它3张能和对里凑出最起码是同花或是顺子，因为顺子和同花比较大，完全可以不考虑二对
		if danPaiCount == 0 {
			return
		}
		liangDui := &LiangDui{
			Dui1Score: tree.listDui[0],
			Dui2Score: tree.listDui[1],
			DanScore:  tree.listDanPai[0],
		}
		listLiangDui = append(listLiangDui, liangDui)
	} else if duiCount == 4 {
		if danPaiCount == 0 { //四对无单牌时，拆最小的对，对为第2小的对和第3小的对
			liangDui := &LiangDui{
				Dui1Score: tree.listDui[2],
				Dui2Score: tree.listDui[1],
				DanScore:  tree.listDui[0],
			}
			listLiangDui = append(listLiangDui, liangDui)
		} else { //四对有单牌时，对为第1小的对及第2小的对 或最小及最大的对
			liangDui1 := &LiangDui{
				Dui1Score: tree.listDui[1],
				Dui2Score: tree.listDui[0],
				DanScore:  tree.listDanPai[0],
			}
			liangDui2 := &LiangDui{
				Dui1Score: tree.listDui[3],
				Dui2Score: tree.listDui[0],
				DanScore:  tree.listDanPai[0],
			}
			listLiangDui = append(listLiangDui, liangDui1)
			listLiangDui = append(listLiangDui, liangDui2)
		}
	} else if duiCount == 5 {
		//5对没有单牌时，说明其它3张能和对里凑出最起码是同花或是顺子，因为顺子和同花比较大，完全可以不考虑二对
		if danPaiCount == 0 {
			return
		} else {
			liangDui1 := &LiangDui{
				Dui1Score: tree.listDui[3],
				Dui2Score: tree.listDui[0],
				DanScore:  tree.listDanPai[0],
			}
			listLiangDui = append(listLiangDui, liangDui1)
		}
	}

	for _, liangDui := range listLiangDui {

		n := NewNode()
		n.normalType = LIANG_DUI
		n.liangDui = liangDui
		isAddDanPai := false
		for _, poker := range tree.pokers {
			if poker.Score == liangDui.Dui1Score || poker.Score == liangDui.Dui2Score {
				n.pokers = append(n.pokers, poker)
			} else if poker.Score == liangDui.DanScore && isAddDanPai == false {
				n.pokers = append(n.pokers, poker)
				isAddDanPai = true
			} else {
				n.rest = append(n.rest, poker)
			}
		}
		tree.Nodes = append(tree.Nodes, n)
	}
}

//splitDui 一对+3张单=对子，这个单张不可能是顺子或是同花中的牌，因为顺子和同花都比一对大，这样对子就没必要出现
func splitDui(tree *Tree) {
	duiCount := len(tree.listDui)
	danPaiCount := len(tree.listDanPai)
	if (duiCount == 1 && danPaiCount >= 3) || (duiCount == 2 && danPaiCount >= 3) || (duiCount == 3 && danPaiCount >= 3) {

		duiScore := 0
		if duiCount == 1 {
			duiScore = tree.listDui[0]
		} else if duiCount == 2 { //有两对时，取大的对
			duiScore = tree.listDui[1]
		} else if duiCount == 3 {
			duiScore = tree.listDui[2]
		}

		dan1Score := tree.listDanPai[0]
		dan2Score := tree.listDanPai[1]
		dan3Score := tree.listDanPai[2]
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
	duiCount := len(tree.listDui)
	danPaiCount := len(tree.listDanPai)
	if duiCount > 0 || danPaiCount < 5 {
		return
	}
	dan1Score := tree.listDanPai[danPaiCount-1] //最大的单牌
	dan2Score := tree.listDanPai[0]             //第1小的单牌
	dan3Score := tree.listDanPai[1]             //第2小的单牌
	dan4Score := tree.listDanPai[2]             //第3小的单牌
	dan5Score := tree.listDanPai[3]             //第4小的单牌
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
