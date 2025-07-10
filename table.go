package pokertexas

import (
	"errors"
	"github.com/karolWazny/pokergo"
	"strings"
)

type Table struct {
	PlayerNames []string           `json:"player_names"`
	Players     map[string]*Player `json:"Players"`
	CurrentGame *Game              `json:"current_game"`
	SmallBlind  int64              `json:"small_blind"`
	BigBlind    int64              `json:"big_blind"`
	DealerIndex int                `json:"dealer_index"`
}

func (table *Table) PlayersList() []*Player {
	result := make([]*Player, len(table.PlayerNames))
	for i, name := range table.PlayerNames {
		result[i] = table.Players[name]
	}
	return result
}

func NewTable(smallBlind int64, bigBlind int64) Table {
	return Table{
		Players:     make(map[string]*Player),
		SmallBlind:  smallBlind,
		BigBlind:    bigBlind,
		DealerIndex: -1,
	}
}

func (table *Table) AddPlayer(player *Player) error {
	for _, existingPlayer := range table.PlayerNames {
		if strings.ToUpper(existingPlayer) == strings.ToUpper(player.GetName()) {
			return errors.New("player already exists")
		}
	}
	table.PlayerNames = append(table.PlayerNames, player.GetName())
	table.Players[player.GetName()] = player
	return nil
}

func (table *Table) StartGame() Game {
	table.DealerIndex = (table.DealerIndex + 1) % len(table.Players)
	orderedPlayerNames := append(table.PlayerNames[table.DealerIndex+1:], table.PlayerNames[:table.DealerIndex+1]...)
	orderedPlayers := make([]*Player, len(orderedPlayerNames))
	for i, name := range orderedPlayerNames {
		orderedPlayers[i] = table.Players[name]
	}
	deck := pokergo.CreateDeck().Shuffled()
	for _, player := range orderedPlayers {
		player.Reset()
		hand, smallerDeck := deck.Deal(2)
		deck = smallerDeck
		player.Hand = hand.Cards
	}
	orderedPlayers[0].CurrentPot = table.SmallBlind
	orderedPlayers[0].Money -= table.SmallBlind
	orderedPlayers[1].CurrentPot = table.BigBlind
	orderedPlayers[1].Money -= table.BigBlind
	table.CurrentGame = &Game{
		PlayerNames:       orderedPlayerNames,
		LastBet:           table.BigBlind,
		Deck:              deck,
		ActivePlayerIndex: 2,
		Community:         make([]pokergo.Card, 0),
		Round:             PREFLOP,
		table:             table,
	}
	return *table.CurrentGame
}

func (table *Table) GetCurrentGame() *Game {
	return table.CurrentGame
}
