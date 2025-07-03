package pokertexas

type TexasHoldEmAction string

const (
	check TexasHoldEmAction = "check"
	fold                    = "fold"
	call                    = "call"
	raise                   = "raise"
)

type TexasHoldEmRound int8

const (
	PREFLOP TexasHoldEmRound = iota
	FLOP
	TURN
	RIVER
	FINISHED
)

func (round TexasHoldEmRound) String() string {
	switch round {
	case PREFLOP:
		return "preflop"
	case FLOP:
		return "flop"
	case TURN:
		return "turn"
	case RIVER:
		return "river"
	case FINISHED:
		return "finished"
	default:
		return "unknown"
	}
}
