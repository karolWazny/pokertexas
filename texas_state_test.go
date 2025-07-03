package pokertexas

import (
	"testing"
)

func TestDumpTableWithoutPlayersNorGame(t *testing.T) {
	table := NewTable(20, 50)
	state := table.DumpState()
	if state.Game != nil {
		t.Errorf("Game should be nil when no game was started.")
	}
	if len(state.Players) != 0 {
		t.Errorf("state.Players should be empty when no players were added.")
	}
	if len(state.Table.Players) != 0 {
		t.Errorf("state.Table.Players should be empty when no players were added.")
	}
	if state.Table.SmallBlind != 20 {
		t.Errorf("state.Table.SmallBlind should be 20, was %v", state.Table.SmallBlind)
	}
	if state.Table.BigBlind != 50 {
		t.Errorf("state.Table.BigBlind should be 50, was %v", state.Table.BigBlind)
	}
	if state.Table.DealerIndex != -1 {
		t.Errorf("state.Table.DealerIndex should be -1, was %v", state.Table.DealerIndex)
	}
}

func TestDumpStateWithOnePlayer(t *testing.T) {
	table := NewTable(20, 50)
	_ = table.AddPlayer(&Player{
		name:  "Player1",
		money: 5000,
	})
	state := table.DumpState()
	if len(state.Players) != 1 {
		t.Errorf("len(state.Players) should be 1, was %v", len(state.Players))
	}
	player, exists := state.Players["Player1"]
	if !exists {
		t.Errorf("state.Players['Player1'] should be present")
	}
	if player.Name != "Player1" {
		t.Errorf("Player1's name shoud be Player1, was %v", player.Name)
	}
	if player.Money != 5000 {
		t.Errorf("Player1's money should be 5000, was %v", player.Money)
	}
	if len(state.Table.Players) != 1 {
		t.Errorf("len(state.Table.Players) should be 1, was %v", len(state.Table.Players))
	}
}

func TestDumpStateWithTwoPlayersAndStartedGame(t *testing.T) {
	table := NewTable(20, 50)
	_ = table.AddPlayer(&Player{
		name:  "Player1",
		money: 5000,
	})
	_ = table.AddPlayer(&Player{
		name:  "Player2",
		money: 5000,
	})
	_ = table.StartGame()
	state := table.DumpState()
	if state.Game == nil {
		t.Errorf("state.Game should not be nil, was %v", state.Game)
	}
	if len(state.Game.Players) != 2 {
		t.Errorf("len(state.Game) should be 2, was %v", len(state.Game.Players))
	}
}
