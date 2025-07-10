package pokertexas

import (
	"fmt"
	"github.com/karolWazny/pokergo"
	"strconv"
)

type PlayerState struct {
	Name            string         `json:"name"`
	Money           int64          `json:"money"`
	Hand            []pokergo.Card `json:"hand"`
	BestCombination []pokergo.Card `json:"best_combination"`
	BestHand        *pokergo.Hand  `json:"best_hand"`
	HasFolded       bool           `json:"has_folded"`
	HasPlayed       bool           `json:"has_played"`
	CurrentPot      int64          `json:"current_pot"`
}

type Player struct {
	s PlayerState
}

func NewPlayer(name string, money int64) Player {
	return Player{s: PlayerState{Name: name, Money: money}}
}

func (p *Player) GetPublicInfo() TexasPlayerPublicInfo {
	playerInfo := TexasPlayerPublicInfo{
		Name:       p.s.Name,
		Money:      p.s.Money,
		HasFolded:  p.s.HasFolded,
		CurrentPot: p.s.CurrentPot,
	}
	if !p.s.HasFolded && p.s.BestHand != nil {
		playerInfo.Cards = p.s.Hand
		playerInfo.BestCards = p.s.BestCombination
		playerInfo.Hand = p.s.BestHand
	}
	return playerInfo
}

func (p *Player) Reset() {
	p.s.Hand = make([]pokergo.Card, 0)
	p.s.BestCombination = make([]pokergo.Card, 0)
	p.s.BestHand = nil
	p.s.HasFolded = false
	p.s.HasPlayed = false
	p.s.CurrentPot = 0
}

func (p *Player) Money() int64 {
	return p.s.Money
}

func (p *Player) Name() string {
	return p.s.Name
}

func (p *Player) String() string {
	return p.s.Name + ", " + strconv.FormatInt(p.s.Money, 10) + " " + fmt.Sprintf("%v", p.s.Hand) + " " + strconv.FormatInt(p.s.CurrentPot, 10)
}

func (p *Player) CurrentPot() int64 {
	return p.s.CurrentPot
}

func (p *Player) HasFolded() bool {
	return p.s.HasFolded
}

func (p *Player) HasPlayed() bool {
	return p.s.HasPlayed
}
