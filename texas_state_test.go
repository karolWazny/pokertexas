package pokertexas

import (
	"testing"
)

func TestDumpTableWithoutPlayers(t *testing.T) {
	table := NewTable(20, 50)
	_ = table.DumpState()
}
