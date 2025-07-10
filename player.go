package pokertexas

import (
	"fmt"
	"github.com/karolWazny/pokergo"
	"strconv"
)

type Player struct {
	name            string
	money           int64
	hand            []pokergo.Card
	bestCombination []pokergo.Card
	bestHand        *pokergo.Hand
	hasFolded       bool
	hasPlayed       bool
	currentPot      int64
}

func NewPlayer(name string, money int64) Player {
	return Player{name: name, money: money}
}

func (p *Player) GetPublicInfo() TexasPlayerPublicInfo {
	playerInfo := TexasPlayerPublicInfo{
		Name:       p.name,
		Money:      p.money,
		HasFolded:  p.hasFolded,
		CurrentPot: p.currentPot,
	}
	if !p.hasFolded && p.bestHand != nil {
		playerInfo.Cards = p.hand
		playerInfo.BestCards = p.bestCombination
		playerInfo.Hand = p.bestHand
	}
	return playerInfo
}

func (p *Player) Reset() {
	p.hand = make([]pokergo.Card, 0)
	p.bestCombination = make([]pokergo.Card, 0)
	p.bestHand = nil
	p.hasFolded = false
	p.hasPlayed = false
	p.currentPot = 0
}

func (p *Player) Money() int64 {
	return p.money
}

func (p *Player) Name() string {
	return p.name
}

func (p *Player) String() string {
	return p.name + ", " + strconv.FormatInt(p.money, 10) + " " + fmt.Sprintf("%v", p.hand) + " " + strconv.FormatInt(p.currentPot, 10)
}
