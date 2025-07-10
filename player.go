package pokertexas

import (
	"fmt"
	"github.com/karolWazny/pokergo"
	"strconv"
)

type Player struct {
	Name            string         `json:"name"`
	Money           int64          `json:"money"`
	Hand            []pokergo.Card `json:"hand"`
	BestCombination []pokergo.Card `json:"best_combination"`
	BestHand        *pokergo.Hand  `json:"best_hand"`
	HasFolded       bool           `json:"has_folded"`
	HasPlayed       bool           `json:"has_played"`
	CurrentPot      int64          `json:"current_pot"`
}

func NewPlayer(name string, money int64) Player {
	return Player{Name: name, Money: money}
}

func (p *Player) GetPublicInfo() TexasPlayerPublicInfo {
	playerInfo := TexasPlayerPublicInfo{
		Name:       p.Name,
		Money:      p.Money,
		HasFolded:  p.HasFolded,
		CurrentPot: p.CurrentPot,
	}
	if !p.HasFolded && p.BestHand != nil {
		playerInfo.Cards = p.Hand
		playerInfo.BestCards = p.BestCombination
		playerInfo.Hand = p.BestHand
	}
	return playerInfo
}

func (p *Player) Reset() {
	p.Hand = make([]pokergo.Card, 0)
	p.BestCombination = make([]pokergo.Card, 0)
	p.BestHand = nil
	p.HasFolded = false
	p.HasPlayed = false
	p.CurrentPot = 0
}

func (p *Player) GetMoney() int64 {
	return p.Money
}

func (p *Player) GetName() string {
	return p.Name
}

func (p *Player) String() string {
	return p.Name + ", " + strconv.FormatInt(p.Money, 10) + " " + fmt.Sprintf("%v", p.Hand) + " " + strconv.FormatInt(p.CurrentPot, 10)
}
