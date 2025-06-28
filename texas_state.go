package pokertexas

type TableState struct {
	Game    *GameDto
	Table   TableDto
	Players map[string]PlayerDto
}

type GameDto struct {
}

type PlayerDto struct {
}

type TableDto struct {
	SmallBlind  int64
	BigBlind    int64
	DealerIndex int
	Players     []string
}
