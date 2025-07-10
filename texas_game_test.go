package pokertexas

import (
	"fmt"
	"slices"
	"testing"
)

func TestThreePlayersCanStartAGame(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	visibleGameState := game.GetVisibleGameState()
	if len(visibleGameState.Players) != 3 {
		t.Errorf("There should be 3 Players")
	}
}

func TestPlayerCannotCallIfThereWasNoRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	availableActions := game.AvailableActions()
	if slices.Contains(availableActions, call) {
		t.Errorf("Player cannot call if there was no raise")
	}
}

func TestPlayerCannotCheckIfThereWasRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	// flop
	game.Raise(50)
	availableActions := game.AvailableActions()
	if slices.Contains(availableActions, check) {
		t.Errorf("Player cannot check if there was raise")
	}
}

func TestRoundIsNotFinishedIfThereWasRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Raise(50)
	round := game.GetVisibleGameState().Round
	if round != PREFLOP {
		t.Errorf("Round should be PREFLOP (is %s)", round)
	}
}

func TestSecondRaiseCausesAReRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	fmt.Println(game.CurrentPlayer())
	game.Raise(50)
	fmt.Println(game.CurrentPlayer())
	player, _ := game.CurrentPlayer()
	fmt.Println(game.CurrentPlayer())
	currentMoney := player.Money
	game.Raise(50)
	fmt.Println(game.CurrentPlayer())
	moneyAfterRaise := player.Money
	difference := currentMoney - moneyAfterRaise
	if difference != 100 {
		t.Errorf("Raising 50 after raise of 50 should cause re-raise (100$ total) (was %d)", difference)
	}
}

func TestCannotRaiseLessThanBigBlind(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	e := game.Raise(25)
	if e == nil {
		t.Errorf("Raising 25 with big blind of 50 should cause an error")
	}
}

func TestCannotRaiseLessThanLastRaise(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	game.Call()
	game.Call()
	game.Check()
	game.Raise(100)
	e := game.Raise(50)
	if e == nil {
		t.Errorf("Raising 50 after a raise of 100 should cause an error")
	}
}

func TestWhenEverybodyFoldsRemainingPlayerWins(t *testing.T) {
	table := prepareThreePlayerTable()
	game := table.StartGame()
	// master calls
	fmt.Println(game.CurrentPlayer())
	game.Call()
	// badmann folds
	fmt.Println(game.CurrentPlayer())
	game.Fold()
	// hanku folds
	fmt.Println(game.CurrentPlayer())
	game.Fold()
	// master is the last player standing
	visibleGameState := game.GetVisibleGameState()
	fmt.Println(visibleGameState)
	if visibleGameState.Round != FINISHED {
		t.Errorf("When only one player remains, the game should be finished (is %s)", visibleGameState.Round)
	}
	winner, e := game.Winner()
	if e != nil {
		t.Errorf("There should be no error fetching winner after game is finished (was %v)", e)
	}
	if winner.Name != "MasterOfDisaster" {
		t.Errorf("Winner should be MasterOfDisaster (was %s)", winner.Name)
	}
	if winner.Money != 1570 {
		t.Errorf("Winner money should be 1570 (initial money + blinds) (was %d)", winner.Money)
	}
}

func TestErrorIsRaisedWhenDuplicatePlayerIsAdded(t *testing.T) {
	table := NewTable(20, 50)
	first := NewPlayer("firstplayer", 1500)
	err := table.AddPlayer(&first)
	if err != nil {
		t.Errorf("There should be no error adding non-duplicate player to table")
	}
	duplicate := NewPlayer("FIRSTPLAYER", 1500)
	err = table.AddPlayer(&duplicate)
	if err == nil {
		t.Errorf("There should be an error adding duplicate player to table")
	}
}

func prepareThreePlayerTable() Table {
	table := NewTable(20, 50)
	master := NewPlayer("MasterOfDisaster", 1500)
	_ = table.AddPlayer(&master)
	badman := NewPlayer("BadMannTM", 1500)
	_ = table.AddPlayer(&badman)
	hanku := NewPlayer("hank.prostokat", 1500)
	_ = table.AddPlayer(&hanku)
	return table
}
