package pokertexas

import (
	"errors"
	"github.com/karolWazny/pokergo"
	"gonum.org/v1/gonum/stat/combin"
	"slices"
)

type Game struct {
	table             *Table
	PlayerNames       []string         `json:"player_names"`
	WinnerName        *string          `json:"winner_name"`
	LastBet           int64            `json:"last_bet"`
	Deck              pokergo.Deck     `json:"deck"`
	ActivePlayerIndex int              `json:"active_player_index"`
	Community         []pokergo.Card   `json:"community"`
	Round             TexasHoldEmRound `json:"round"`
}

func (game *Game) Call() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, call) {
		return errors.New("call action not available")
	}
	currentPlayer := game.unsafeGetCurrentPlayer()
	pot := game.getPreviousPlayerPot()
	difference := pot - currentPlayer.CurrentPot
	currentPlayer.CurrentPot = pot
	currentPlayer.Money -= difference
	game.nextPlayer()
	return nil
}

func (game *Game) Fold() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, fold) {
		return errors.New("fold action not available")
	}
	game.unsafeGetCurrentPlayer().HasFolded = true
	game.unsafeGetCurrentPlayer().Hand = make([]pokergo.Card, 0)
	if game.playersInGame() == 1 {
		// finish game
		lastActivePlayerIndex := slices.IndexFunc(game.PlayerNames, func(player string) bool {
			return !game.table.Players[player].HasFolded
		})
		game.WinnerName = &game.PlayerNames[lastActivePlayerIndex]
		game.ActivePlayerIndex = -1
		game.Round = FINISHED
		game.transferPotToWinner()
	} else {
		game.nextPlayer()
	}
	return nil
}

func (game *Game) Check() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, check) {
		return errors.New("check action not available")
	}
	game.nextPlayer()
	return nil
}

func (game *Game) Raise(amount int64) error {
	if amount < game.LastBet {
		return errors.New("amount must be greater than last bet")
	}
	game.LastBet = amount
	realAmount := game.getPreviousPlayerPot() - game.getCurrentPlayerPot() + amount
	currentPlayer, e := game.CurrentPlayer()
	if e == nil {
		currentPlayer.CurrentPot += realAmount
		currentPlayer.Money -= realAmount
	}
	game.nextPlayer()
	return nil
}

func (game *Game) Winner() (*Player, error) {
	if game.Round != FINISHED {
		return nil, errors.New("there is no winner before game end")
	}
	return game.table.Players[*game.WinnerName], nil
}

func (game *Game) CurrentPlayer() (*Player, error) {
	if game.Round != FINISHED {
		return game.unsafeGetCurrentPlayer(), nil
	}
	return nil, errors.New("in finished game there is no current player")
}

func (game *Game) AvailableActions() []TexasHoldEmAction {
	if game.Round == FINISHED {
		return []TexasHoldEmAction{}
	}
	currentPlayer := game.unsafeGetCurrentPlayer()
	previousPlayerPot := game.getPreviousPlayerPot()
	availableActions := []TexasHoldEmAction{fold, raise}
	if previousPlayerPot == currentPlayer.CurrentPot {
		availableActions = append(availableActions, check)
	} else {
		availableActions = append(availableActions, call)
	}
	return availableActions
}

func (game *Game) CommunityCards() []pokergo.Card {
	return game.Community
}

func (game *Game) GetVisibleGameState() VisibleGameState {
	visibleGameState := VisibleGameState{}
	players := make([]TexasPlayerPublicInfo, len(game.PlayerNames))
	for i, player := range game.PlayerNames {
		players[i] = game.table.Players[player].GetPublicInfo()
	}
	visibleGameState.Players = players
	activePlayer, e := game.CurrentPlayer()
	if e == nil {
		activePlayerInfo := activePlayer.GetPublicInfo()
		visibleGameState.ActivePlayer = &activePlayerInfo
	}
	visibleGameState.Round = game.Round
	visibleGameState.Dealer = game.table.Players[game.PlayerNames[len(game.PlayerNames)-1]].GetPublicInfo()
	visibleGameState.Community = game.CommunityCards()
	if game.WinnerName != nil {
		visibleGameState.Winner = *game.WinnerName
	}
	return visibleGameState
}

func (game *Game) transferPotToWinner() {
	for _, player := range game.table.Players {
		game.table.Players[*game.WinnerName].Money += player.CurrentPot
	}
}

func (game *Game) unsafeGetCurrentPlayer() *Player {
	return game.table.Players[game.PlayerNames[game.ActivePlayerIndex]]
}

func (game *Game) playersInGame() int {
	playersInGame := 0
	for _, player := range game.PlayerNames {
		if !game.table.Players[player].HasFolded {
			playersInGame++
		}
	}
	return playersInGame
}

func (game *Game) getPreviousPlayerPot() int64 {
	for i := 1; i < len(game.PlayerNames); i++ {
		previousPlayerIndex := (game.ActivePlayerIndex - i + len(game.PlayerNames)) % len(game.PlayerNames)
		if !game.table.Players[game.PlayerNames[previousPlayerIndex]].HasFolded {
			return game.table.Players[game.PlayerNames[previousPlayerIndex]].CurrentPot
		}
	}
	panic("There should be at least two active Players!")
}

func (game *Game) getCurrentPlayerPot() int64 {
	return game.unsafeGetCurrentPlayer().CurrentPot
}

func (game *Game) nextPlayer() {
	game.unsafeGetCurrentPlayer().HasPlayed = true
	game.changeActivePlayerToFirstNonFolded()
	isRoundFinished := game.isCurrentRoundFinished()
	if isRoundFinished {
		game.finishRound()
	}
}

func (game *Game) finishRound() {
	if game.Round == RIVER {
		// trigger showdown
		game.showdown()
		return
	}
	game.ActivePlayerIndex = len(game.PlayerNames) - 1
	game.changeActivePlayerToFirstNonFolded()
	for _, player := range game.table.PlayersList() {
		player.HasPlayed = false
	}
	game.showCommunityCards()
	game.Round++
}

func (game *Game) showdown() {
	game.determineBestHandForEachPlayer()
	game.findWinner()
	game.transferPotToWinner()
	game.Round = FINISHED
}

func (game *Game) showCommunityCards() {
	_, game.Deck = game.Deck.Deal(1)
	cardsToShow := game.numberOfCardsToShow()
	game.dealCardsToCommunity(cardsToShow)
}

func (game *Game) numberOfCardsToShow() int {
	cardsToShow := 1
	isFlop := game.Round == PREFLOP
	if isFlop {
		cardsToShow = 3
	}
	return cardsToShow
}

func (game *Game) dealCardsToCommunity(cardsToShow int) {
	newCards, deck := game.Deck.Deal(cardsToShow)
	game.Deck = deck
	game.Community = append(game.CommunityCards(), newCards.Cards...)
}

func (game *Game) findWinner() {
	bestPlayer := &Player{}
	bestHand := pokergo.CreateLowGuardian()
	for _, player := range game.table.Players {
		if !player.HasFolded {
			comparisonResult := pokergo.CompareHands(bestHand, *player.BestHand)
			if comparisonResult != pokergo.FirstWins {
				bestPlayer = player
				bestHand = *player.BestHand
			}
		}
	}
	game.WinnerName = &bestPlayer.Name
}

func (game *Game) determineBestHandForEachPlayer() {
	for _, player := range game.table.Players {
		game.determineBestHandForPlayer(player)
	}
}

func (game *Game) determineBestHandForPlayer(player *Player) {
	if !player.HasFolded {
		allCards := append(game.Community, player.Hand...)
		bestHand, bestCombination := game.findBestHand(allCards)
		player.BestHand = &bestHand
		player.BestCombination = bestCombination
	}
}

func (game *Game) findBestHand(allCards []pokergo.Card) (pokergo.Hand, []pokergo.Card) {
	combinations := combin.NewCombinationGenerator(7, 5)
	combinationMapping := make([]int, 5)
	checkedHand := make([]pokergo.Card, 5)
	bestHand := pokergo.CreateLowGuardian()
	bestCombination := make([]pokergo.Card, 5)
	for combinations.Next() {
		combinations.Combination(combinationMapping)
		for index, value := range combinationMapping {
			checkedHand[index] = allCards[value]
		}
		recognisedHand, e := pokergo.RecogniseHand(pokergo.DeckOf(checkedHand...))
		if e != nil {
			continue
		}
		result := pokergo.CompareHands(bestHand, recognisedHand)
		if result != pokergo.FirstWins {
			bestHand = recognisedHand
			copy(bestCombination, checkedHand)
		}
	}
	return bestHand, bestCombination
}

func (game *Game) changeActivePlayerToFirstNonFolded() {
	game.incrementActivePlayerIndex()
	for game.unsafeGetCurrentPlayer().HasFolded {
		game.incrementActivePlayerIndex()
	}
}

func (game *Game) incrementActivePlayerIndex() {
	game.ActivePlayerIndex = (game.ActivePlayerIndex + 1) % len(game.table.Players)
}

func (game *Game) isCurrentRoundFinished() bool {
	uniquePots := make(map[int64]bool)
	for _, player := range game.table.Players {
		if !player.HasFolded {
			if !player.HasPlayed {
				return false
			}
			uniquePots[player.CurrentPot] = true
		}
	}
	return len(uniquePots) == 1
}
