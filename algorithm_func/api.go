package algorithm_func

import (
	"log"
)

//CalResult 计算牌的api
func CalResult(pokers []*Poker) {
	fatherTree := NewTree(pokers)
	if isSpecial, specialType := CalSpecial(fatherTree); isSpecial {
		log.Printf("该牌型为特殊牌型：%s", specialType.String())
		return
	}
	res := CalNormalResults(fatherTree)
	filterRes := SortFilterResult(res)
	for i, result := range filterRes {
		log.Printf("智能排序：牌型[%d] = %s", i, result.String())
	}
}

func CalCardType(pokers []*Poker) *Node {
	n := NewNode()
	if len(pokers) == 5 {
		tree := NewTree(pokers)
		split(tree)
		n = tree.Nodes[0]
	}
	if len(pokers) == 3 {
		n.pokers = pokers
		if pokers[0].Score == pokers[1].Score && pokers[1].Score == pokers[2].Score {
			n.normalType = SAN_TIAO
			n.sanTiao = &SanTiao{
				SanTiaoScore: pokers[0].Score,
			}
		} else if pokers[0].Score == pokers[1].Score {
			n.normalType = DUI_ZI
			n.dui = &Dui{
				DuiScore:  pokers[1].Score,
				Dan3Score: pokers[2].Score,
			}
		} else if pokers[0].Score == pokers[2].Score {
			n.normalType = DUI_ZI
			n.dui = &Dui{
				DuiScore:  pokers[0].Score,
				Dan3Score: pokers[1].Score,
			}
		} else if pokers[1].Score == pokers[2].Score {
			n.normalType = DUI_ZI
			n.dui = &Dui{
				DuiScore:  pokers[1].Score,
				Dan3Score: pokers[0].Score,
			}
		} else {
			n.normalType = WU_LONG
		}
	}
	return n
}
