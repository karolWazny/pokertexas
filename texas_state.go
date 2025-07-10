package pokertexas

import "github.com/karolWazny/pokergo"

type TableState struct {
	Table   TableDto             `json:"table"`
	Players map[string]PlayerDto `json:"players"`
}

type GameDto struct {
	Players           []TexasPlayerDto `json:"players"`
	Winner            string           `json:"winner"`
	LastBet           int64            `json:"lastBet"`
	Deck              []pokergo.Card   `json:"deck"`
	ActivePlayerIndex int              `json:"activePlayerIndex"`
	Community         []pokergo.Card   `json:"community"`
	Round             TexasHoldEmRound `json:"round"`
}

type PlayerDto struct {
	Name  string `json:"name"`
	Money int64  `json:"money"`
}

type TexasPlayerDto struct {
	Name            string         `json:"name"`
	HasFolded       bool           `json:"has_folded"`
	HasPlayed       bool           `json:"has_played"`
	CurrentPot      int64          `json:"current_pot"`
	BestHand        *pokergo.Hand  `json:"best_hand"`
	BestCombination []pokergo.Card `json:"best_combination"`
	Hand            []pokergo.Card `json:"hand"`
}

type TableDto struct {
	Game        *GameDto `json:"game"`
	SmallBlind  int64    `json:"small_blind"`
	BigBlind    int64    `json:"big_blind"`
	DealerIndex int      `json:"dealer_index"`
	Players     []string `json:"players"`
}
