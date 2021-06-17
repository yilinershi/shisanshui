package shi_san_shui

import (
	"fmt"
	"log"
)

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

	desc := fmt.Sprintf("结点：{左:【%s】= {%s},中:【%s】= {%s},右:【%s】= {%s},好牌值=【%d】}",
		this.Left.normalType, leftPokerDesc, this.Middle.normalType, middlePokerDesc, this.Right.normalType, rightPokerDesc, this.BestScore)
	return desc
}

func CalResult(pokers []*Poker) {

	mainTree := NewTree(pokers)
	if isSpecial, specialType := mainTree.CalSpecial(); isSpecial {
		log.Printf("该牌型为特殊牌型：%s", specialType.String())
		return
	}

	res := CalNormalResults(mainTree)
	for i, result := range res {
		log.Printf("所有牌型=>[%d]=%s", i, result.String())
	}

	best := CalBest(res)
	if best != nil {

		log.Println("最好牌型=> ", best.String())
	}
}

func CalNormalResults(mainTree *Tree) []*ResultNormal {
	mainTree.CalNormal()
	normalInfo := make([]*ResultNormal, 0)
	for _, node1 := range mainTree.Nodes {
		fmt.Println("left=", node1.String())
		tree2 := NewTree(node1.rest)
		tree2.CalNormal()
		for _, node2 := range tree2.Nodes {
			fmt.Println("middle=", node2.String())
			if node1.CompareExternal(node2)!=Worse {
				bestScore := 0
				switch node1.normalType {
				case TONG_HUA_SHUN:
					bestScore = bestScore + 5 + int(TONG_HUA_SHUN)
				case TIE_ZHI:
					bestScore = bestScore + 5 + int(TIE_ZHI)
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
						sanTiaoScore: node3.pokers[0].Score,
					}
					bestScore = bestScore + 3 + int(SAN_TIAO)
				} else if node3.pokers[0].Score == node3.pokers[1].Score || node3.pokers[1].Score == node3.pokers[2].Score || node3.pokers[0].Score == node3.pokers[2].Score {
					node3.normalType = DUI_ZI
					node3.dui = &Dui{
						duiScore: node3.pokers[0].Score,
					}
					bestScore = bestScore + 2 + int(DUI_ZI)
				} else {
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
			if result.BestScore > best.BestScore {
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
