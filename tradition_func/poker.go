package tradition_func

type Poker struct {
	Point PokerPoint
	Hua   PokerHua
	Score int
	Desc  string
}

func NewPoker(point PokerPoint, hua PokerHua) *Poker {
	return &Poker{
		Point: point,
		Hua:   hua,
		Score: point.Score(),
		Desc:  hua.String() + point.String(),
	}
}

func NewPokerById(id byte) *Poker {
	p := &Poker{
		Point: PokerPoint(id & 0x0F),
		Hua:   PokerHua(id >> 4),
	}
	p.Score = p.Point.Score()
	p.Desc = p.Hua.String() + p.Point.String()
	return p
}
