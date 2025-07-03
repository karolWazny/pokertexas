package pokertexas

import "github.com/karolWazny/pokergo"

type TableState struct {
	Game    *GameDto
	Table   TableDto
	Players map[string]PlayerDto
}

type GameDto struct {
	Players []string
}

type PlayerDto struct {
	Name  string
	Money int64
}

type TexasPlayerDto struct {
	Name            string
	HasFolded       bool
	HasPlayed       bool
	CurrentPot      int64
	BestHand        *pokergo.Hand
	BestCombination []pokergo.Card
	Hand            []pokergo.Card
}

type TableDto struct {
	SmallBlind  int64
	BigBlind    int64
	DealerIndex int
	Players     []string
}
