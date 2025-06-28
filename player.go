package pokertexas

import "strconv"

type Player struct {
	name  string
	money int64
}

func (p *Player) Money() int64 {
	return p.money
}

func (p *Player) Name() string {
	return p.name
}

func NewPlayer(name string, money int64) Player {
	return Player{name: name, money: money}
}

func (p *Player) String() string {
	return p.name + ", " + strconv.FormatInt(p.money, 10)
}
