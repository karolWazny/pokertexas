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
	if state.Table.DealerIndex != 0 {
		t.Errorf("state.Table.DealerIndex should be 0, was %v", state.Table.DealerIndex)
	}
}
