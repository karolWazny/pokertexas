package pokertexas

import (
	"errors"
	"github.com/karolWazny/pokergo"
	"gonum.org/v1/gonum/stat/combin"
	"slices"
)

type Game struct {
	table             *Table
	players           []*Player
	winner            *Player
	lastBet           int64
	deck              pokergo.Deck
	activePlayerIndex int
	community         []pokergo.Card
	round             TexasHoldEmRound
}

func (game *Game) Call() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, call) {
		return errors.New("call action not available")
	}
	currentPlayer := game.unsafeGetCurrentPlayer()
	pot := game.getPreviousPlayerPot()
	difference := pot - currentPlayer.currentPot
	currentPlayer.currentPot = pot
	currentPlayer.money -= difference
	game.nextPlayer()
	return nil
}

func (game *Game) Fold() error {
	availableActions := game.AvailableActions()
	if !slices.Contains(availableActions, fold) {
		return errors.New("fold action not available")
	}
	game.unsafeGetCurrentPlayer().hasFolded = true
	game.unsafeGetCurrentPlayer().hand = make([]pokergo.Card, 0)
	if game.playersInGame() == 1 {
		// finish game
		lastActivePlayerIndex := slices.IndexFunc(game.players, func(player *Player) bool {
			return !player.hasFolded
		})
		game.winner = game.players[lastActivePlayerIndex]
		game.activePlayerIndex = -1
		game.round = FINISHED
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
	if amount < game.lastBet {
		return errors.New("amount must be greater than last bet")
	}
	game.lastBet = amount
	realAmount := game.getPreviousPlayerPot() - game.getCurrentPlayerPot() + amount
	currentPlayer, e := game.CurrentPlayer()
	if e == nil {
		currentPlayer.currentPot += realAmount
		currentPlayer.money -= realAmount
	}
	game.nextPlayer()
	return nil
}

func (game *Game) Winner() (*Player, error) {
	if game.round != FINISHED {
		return nil, errors.New("there is no winner before game end")
	}
	return game.winner, nil
}

func (game *Game) CurrentPlayer() (*Player, error) {
	if game.round != FINISHED {
		return game.unsafeGetCurrentPlayer(), nil
	}
	return nil, errors.New("in finished game there is no current player")
}

func (game *Game) AvailableActions() []TexasHoldEmAction {
	if game.round == FINISHED {
		return []TexasHoldEmAction{}
	}
	currentPlayer := game.unsafeGetCurrentPlayer()
	previousPlayerPot := game.getPreviousPlayerPot()
	availableActions := []TexasHoldEmAction{fold, raise}
	if previousPlayerPot == currentPlayer.currentPot {
		availableActions = append(availableActions, check)
	} else {
		availableActions = append(availableActions, call)
	}
	return availableActions
}

func (game *Game) CommunityCards() []pokergo.Card {
	return game.community
}

func (game *Game) GetVisibleGameState() VisibleGameState {
	visibleGameState := VisibleGameState{}
	players := make([]TexasPlayerPublicInfo, len(game.players))
	for i, player := range game.players {
		players[i] = player.GetPublicInfo()
	}
	visibleGameState.Players = players
	activePlayer, e := game.CurrentPlayer()
	if e == nil {
		activePlayerInfo := activePlayer.GetPublicInfo()
		visibleGameState.ActivePlayer = &activePlayerInfo
	}
	visibleGameState.Round = game.round
	visibleGameState.Dealer = game.players[len(game.players)-1].GetPublicInfo()
	visibleGameState.Community = game.CommunityCards()
	if game.winner != nil {
		visibleGameState.Winner = game.winner.name
	}
	return visibleGameState
}

func (game *Game) transferPotToWinner() {
	for _, player := range game.players {
		game.winner.money += player.currentPot
	}
}

func (game *Game) unsafeGetCurrentPlayer() *Player {
	return game.players[game.activePlayerIndex]
}

func (game *Game) playersInGame() int {
	playersInGame := 0
	for _, player := range game.players {
		if !player.hasFolded {
			playersInGame++
		}
	}
	return playersInGame
}

func (game *Game) getPreviousPlayerPot() int64 {
	for i := 1; i < len(game.players); i++ {
		previousPlayerIndex := (game.activePlayerIndex - i + len(game.players)) % len(game.players)
		if !game.players[previousPlayerIndex].hasFolded {
			return game.players[previousPlayerIndex].currentPot
		}
	}
	panic("There should be at least two active players!")
}

func (game *Game) getCurrentPlayerPot() int64 {
	return game.unsafeGetCurrentPlayer().currentPot
}

func (game *Game) nextPlayer() {
	game.unsafeGetCurrentPlayer().hasPlayed = true
	game.changeActivePlayerToFirstNonFolded()
	isRoundFinished := game.isCurrentRoundFinished()
	if isRoundFinished {
		game.finishRound()
	}
}

func (game *Game) finishRound() {
	if game.round == RIVER {
		// trigger showdown
		game.showdown()
		return
	}
	game.activePlayerIndex = len(game.players) - 1
	game.changeActivePlayerToFirstNonFolded()
	for _, player := range game.players {
		player.hasPlayed = false
	}
	game.showCommunityCards()
	game.round++
}

func (game *Game) showdown() {
	game.determineBestHandForEachPlayer()
	game.findWinner()
	game.transferPotToWinner()
	game.round = FINISHED
}

func (game *Game) showCommunityCards() {
	_, game.deck = game.deck.Deal(1)
	cardsToShow := game.numberOfCardsToShow()
	game.dealCardsToCommunity(cardsToShow)
}

func (game *Game) numberOfCardsToShow() int {
	cardsToShow := 1
	isFlop := game.round == PREFLOP
	if isFlop {
		cardsToShow = 3
	}
	return cardsToShow
}

func (game *Game) dealCardsToCommunity(cardsToShow int) {
	newCards, deck := game.deck.Deal(cardsToShow)
	game.deck = deck
	game.community = append(game.CommunityCards(), newCards.Cards...)
}

func (game *Game) findWinner() {
	bestPlayer := &Player{}
	bestHand := pokergo.CreateLowGuardian()
	for _, player := range game.players {
		if !player.hasFolded {
			comparisonResult := pokergo.CompareHands(bestHand, *player.bestHand)
			if comparisonResult != pokergo.FirstWins {
				bestPlayer = player
				bestHand = *player.bestHand
			}
		}
	}
	game.winner = bestPlayer
}

func (game *Game) determineBestHandForEachPlayer() {
	for _, player := range game.players {
		game.determineBestHandForPlayer(player)
	}
}

func (game *Game) determineBestHandForPlayer(player *Player) {
	if !player.hasFolded {
		allCards := append(game.community, player.hand...)
		bestHand, bestCombination := game.findBestHand(allCards)
		player.bestHand = &bestHand
		player.bestCombination = bestCombination
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
	for game.unsafeGetCurrentPlayer().hasFolded {
		game.incrementActivePlayerIndex()
	}
}

func (game *Game) incrementActivePlayerIndex() {
	game.activePlayerIndex = (game.activePlayerIndex + 1) % len(game.players)
}

func (game *Game) isCurrentRoundFinished() bool {
	uniquePots := make(map[int64]bool)
	for _, player := range game.players {
		if !player.hasFolded {
			if !player.hasPlayed {
				return false
			}
			uniquePots[player.currentPot] = true
		}
	}
	return len(uniquePots) == 1
}
