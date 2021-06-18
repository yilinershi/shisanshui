package shi_san_shui

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
	for i, result := range res {
		log.Printf("普通牌型[%d] = %s", i, result.String())
	}
	best := CalBest(res)
	if best != nil {
		log.Println("最好牌型 = ", best.String())
	}
}
