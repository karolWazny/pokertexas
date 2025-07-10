package pokertexas

import (
	"errors"
	"github.com/karolWazny/pokergo"
	"strings"
)

type TableState struct {
	PlayerNames []string           `json:"player_names"`
	Players     map[string]*Player `json:"Players"`
	SmallBlind  int64              `json:"small_blind"`
	BigBlind    int64              `json:"big_blind"`
	DealerIndex int                `json:"dealer_index"`
}

type Table struct {
	currentGame *Game
	s           *TableState
}

func (t *Table) PlayersList() []*Player {
	result := make([]*Player, len(t.s.PlayerNames))
	for i, name := range t.s.PlayerNames {
		result[i] = t.s.Players[name]
	}
	return result
}

func NewTable(smallBlind int64, bigBlind int64) Table {
	return Table{
		s: &TableState{
			Players:     make(map[string]*Player),
			SmallBlind:  smallBlind,
			BigBlind:    bigBlind,
			DealerIndex: -1,
		},
	}
}

func (t *Table) AddPlayer(player *Player) error {
	for _, existingPlayer := range t.s.PlayerNames {
		if strings.ToUpper(existingPlayer) == strings.ToUpper(player.Name()) {
			return errors.New("player already exists")
		}
	}
	t.s.PlayerNames = append(t.s.PlayerNames, player.Name())
	t.s.Players[player.Name()] = player
	return nil
}

func (t *Table) StartGame() Game {
	t.s.DealerIndex = (t.s.DealerIndex + 1) % len(t.s.Players)
	orderedPlayerNames := append(t.s.PlayerNames[t.s.DealerIndex+1:], t.s.PlayerNames[:t.s.DealerIndex+1]...)
	orderedPlayers := make([]*Player, len(orderedPlayerNames))
	for i, name := range orderedPlayerNames {
		orderedPlayers[i] = t.s.Players[name]
	}
	deck := pokergo.CreateDeck().Shuffled()
	for _, player := range orderedPlayers {
		player.Reset()
		hand, smallerDeck := deck.Deal(2)
		deck = smallerDeck
		player.s.Hand = hand.Cards
	}
	orderedPlayers[0].s.CurrentPot = t.s.SmallBlind
	orderedPlayers[0].s.Money -= t.s.SmallBlind
	orderedPlayers[1].s.CurrentPot = t.s.BigBlind
	orderedPlayers[1].s.Money -= t.s.BigBlind
	t.currentGame = &Game{
		s: &GameState{
			PlayerNames:       orderedPlayerNames,
			LastBet:           t.s.BigBlind,
			Deck:              deck,
			ActivePlayerIndex: 2,
			Community:         make([]pokergo.Card, 0),
			Round:             PREFLOP,
		},
		table: t,
	}
	return *t.currentGame
}

func (t *Table) GetCurrentGame() *Game {
	return t.currentGame
}
