package tradition_func

import (
	"fmt"
	"sort"
)

type ResultNormal struct {
	BestScore int //好牌值
	Head      *Dun
	Middle    *Dun
	Tail      *Dun
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
		this.Head.Type, headPokerDesc, this.Middle.Type, middlePokerDesc, this.Tail.Type, tailPokerDesc, this.BestScore)
}

func CalNormal(pokers []*Poker) []*ResultNormal {
	listResultNormal := make([]*ResultNormal, 0)

	node1s := ForEachPoker(pokers)
	for _, node1 := range node1s {

		//fmt.Println("---------------------------------------------------------------------")
		//fmt.Println("下墩=", node1.String())

		node2s := ForEachPoker(node1.rest)
		for _, node2 := range node2s {
			if node1.dun.Compare(node2.dun) == Worse {
				continue
			}
			//fmt.Println("中墩=", node2.String())

			dun1 := node1.dun
			dun2 := node2.dun
			dun3 := NewDun(node2.rest)

			//特殊牌型 三顺子
			if dun1.Type == SHUN_ZI && dun2.Type == SHUN_ZI {
				if (dun3.pokers[0].Score+1 == dun3.pokers[1].Score && dun3.pokers[1].Score+1 == dun3.pokers[2].Score) || (dun3.pokers[0].Point == Poker2 && dun3.pokers[1].Point == Poker3 && dun3.pokers[2].Point == PokerA) {
					dun3.Type = SHUN_ZI
					listResultNormal = make([]*ResultNormal, 0)
					result := &ResultNormal{
						Tail:      dun1,
						Middle:    dun2,
						Head:      dun3,
						BestScore: 0,
					}
					listResultNormal = append(listResultNormal, result)
					return listResultNormal
				}
			}

			if dun2.Compare(dun3) != Worse {
				bestScore := 0
				switch dun1.Type {
				case TONG_HUA_SHUN:
					bestScore += 5 + int(dun1.Type)
				case TIE_ZHI:
					bestScore += 4 + int(dun1.Type)
				default:
					bestScore += 1 + int(dun1.Type)
				}
				switch dun2.Type {
				case TONG_HUA_SHUN:
					bestScore += 10 + int(dun2.Type)
				case TIE_ZHI:
					bestScore += 8 + int(dun2.Type)
				case HU_LU:
					bestScore += 2 + int(dun2.Type)
				default:
					bestScore += 1 + int(dun2.Type)
				}
				switch dun3.Type {
				case SAN_TIAO:
					bestScore += 3 + int(dun3.Type)
				case DUI_ZI:
					bestScore += 2 + int(dun3.Type)
				default:
					bestScore += 1 + int(dun3.Type)
				}
				result := &ResultNormal{
					Tail:      dun1,
					Middle:    dun2,
					Head:      dun3,
					BestScore: bestScore,
				}

				listResultNormal = append(listResultNormal, result)
			}
		}
	}

	return listResultNormal
}

func ForEachPoker(pokers []*Poker) []*Node {
	nodes := make([]*Node, 0)
	treePokerCount := len(pokers)
	SortPoker(pokers)
	for i1 := 0; i1 < treePokerCount-4; i1++ {
		for i2 := i1 + 1; i2 < treePokerCount-3; i2++ {
			if pokers[i2].Score != pokers[i1].Score && pokers[i2].Hua != pokers[i1].Hua && pokers[i2].Score != pokers[i1].Score+1 {
				continue
			}
			for i3 := i2 + 1; i3 < treePokerCount-2; i3++ {
				if pokers[i3].Score != pokers[i2].Score && pokers[i3].Hua != pokers[i3].Hua && pokers[i3].Score != pokers[i2].Score+1 {
					continue
				}
				for i4 := i3 + 1; i4 < treePokerCount-1; i4++ {
					if pokers[i4].Score != pokers[i4].Score && pokers[i4].Hua != pokers[i3].Hua && pokers[i4].Score != pokers[i3].Score+1 {
						continue
					}
					for i5 := i4 + 1; i5 < treePokerCount; i5++ {
						n := &Node{}
						n.rest = make([]*Poker, 0)
						for i, poker := range pokers {
							if i != i1 && i != i2 && i != i3 && i != i4 && i != i5 {
								n.rest = append(n.rest, poker)
							}
						}

						tempPokers := []*Poker{
							pokers[i1],
							pokers[i2],
							pokers[i3],
							pokers[i4],
							pokers[i5],
						}

						n.dun = NewDun(tempPokers)
						if n.dun.Type != WU_LONG {
							nodes = append(nodes, n)
						}
					}
				}
			}
		}
	}

	return nodes
}

//SortFilterResult 计算出最好的普通牌型,优先保证分值，其次保证上墩，再次保证中墩，最后保证下墩
func SortFilterResult(resultList []*ResultNormal) []*ResultNormal {
	sort.Slice(resultList, func(i, j int) bool {
		if resultList[i].BestScore > resultList[j].BestScore {
			return true
		} else if resultList[i].BestScore < resultList[j].BestScore {
			return false
		}
		if resultList[i].Head.Compare(resultList[j].Head) == Better {
			return true
		} else if resultList[i].Head.Compare(resultList[j].Head) == Worse {
			return false
		}
		if resultList[i].Tail.Compare(resultList[j].Tail) == Better {
			return true
		} else if resultList[i].Tail.Compare(resultList[j].Tail) == Worse {
			return false
		}
		if resultList[i].Middle.Compare(resultList[j].Middle) == Better {
			return true
		} else if resultList[i].Middle.Compare(resultList[j].Middle) == Worse {
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
		} else if result.Head.Type != last.Head.Type || result.Middle.Type != last.Middle.Type || result.Tail.Type != last.Tail.Type {
			last = result
			filterRes = append(filterRes, result)
		}
	}

	return filterRes
}
