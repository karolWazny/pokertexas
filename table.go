package pokertexas

import (
	"errors"
	"github.com/karolWazny/pokergo"
	"strings"
)

type Table struct {
	players     []*Player
	currentGame *Game
	smallBlind  int64
	bigBlind    int64
	dealerIndex int
}

func (table *Table) Players() []*Player {
	return table.players
}

func NewTable(smallBlind int64, bigBlind int64) Table {
	return Table{
		players:     make([]*Player, 0),
		smallBlind:  smallBlind,
		bigBlind:    bigBlind,
		dealerIndex: -1,
	}
}

func (table *Table) AddPlayer(player *Player) error {
	for _, existingPlayer := range table.players {
		if strings.ToUpper(existingPlayer.Name()) == strings.ToUpper(player.Name()) {
			return errors.New("player already exists")
		}
	}
	table.players = append(table.players, player)
	return nil
}

func (table *Table) StartGame() Game {
	table.dealerIndex = (table.dealerIndex + 1) % len(table.players)
	orderedPlayers := append(table.players[table.dealerIndex+1:], table.players[:table.dealerIndex+1]...)
	texasPlayers := make([]*TexasPlayer, len(orderedPlayers))
	deck := pokergo.CreateDeck().Shuffled()
	for i, player := range orderedPlayers {
		hand, smallerDeck := deck.Deal(2)
		deck = smallerDeck
		texasPlayers[i] = &TexasPlayer{
			player:     player,
			hand:       hand.Cards,
			hasFolded:  false,
			currentPot: 0,
		}
	}
	texasPlayers[0].currentPot = table.smallBlind
	texasPlayers[0].player.money -= table.smallBlind
	texasPlayers[1].currentPot = table.bigBlind
	texasPlayers[1].player.money -= table.bigBlind
	table.currentGame = &Game{
		players:           texasPlayers,
		lastBet:           table.bigBlind,
		deck:              deck,
		activePlayerIndex: 2,
		community:         make([]pokergo.Card, 0),
		round:             PREFLOP,
	}
	return *table.currentGame
}

func (table *Table) DumpState() TableState {
	players := make(map[string]PlayerDto, len(table.players))
	playerNames := make([]string, len(table.players))
	for i, player := range table.players {
		playerNames[i] = player.Name()
		players[player.Name()] = PlayerDto{
			Name:  player.Name(),
			Money: player.Money(),
		}
	}
	state := TableState{
		Table: TableDto{
			SmallBlind:  table.smallBlind,
			BigBlind:    table.bigBlind,
			Players:     playerNames,
			DealerIndex: table.dealerIndex,
		},
		Players: players,
	}
	if table.GetCurrentGame() != nil {
		state.Game = &GameDto{
			Players: playerNames,
		}
	}
	return state
}

func (table *Table) GetCurrentGame() *Game {
	return table.currentGame
}
