package pokertexas

type SerializableStateDto struct {
	Table   *TableState             `json:"table"`
	Game    *GameState              `json:"game"`
	Players map[string]*PlayerState `json:"players"`
}
