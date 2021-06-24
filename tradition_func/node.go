package tradition_func

import "fmt"

type Node struct {
	dun  *Dun
	rest []*Poker
}

func (this *Node) String() string {
	pokerDesc := ""
	for _, poker := range this.dun.pokers {
		pokerDesc += poker.Desc
	}
	restDesc := ""
	for _, poker := range this.rest {
		restDesc += poker.Desc
	}
	desc := fmt.Sprintf("结点：{类型=【%s】,墩=【%s】,余=【%s】}", this.dun.Type, pokerDesc, restDesc)
	return desc
}
