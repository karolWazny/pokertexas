package pokertexas

import (
	"fmt"
	"github.com/karolWazny/pokergo"
)

type VisibleGameState struct {
	Players      []TexasPlayerPublicInfo
	Round        TexasHoldEmRound
	ActivePlayer *TexasPlayerPublicInfo
	Winner       string
	Dealer       TexasPlayerPublicInfo
	Community    []pokergo.Card
}

func (gameState VisibleGameState) Print() {
	fmt.Printf("Little Friendly Game of Poker, stage: %s\n", gameState.Round)
	fmt.Printf("Dealer: %s\n", gameState.Dealer.Name)
	fmt.Printf("Community Cards:\n")
	for _, card := range gameState.Community {
		fmt.Printf("- %s\n", card)
	}
	fmt.Printf("PlayersList:\n")
	for _, player := range gameState.Players {
		fmt.Printf("- %s\n", player)
	}
	if gameState.ActivePlayer != nil {
		fmt.Printf("Now playing: %s\n", gameState.ActivePlayer.Name)
	}
	if gameState.Winner != "" {
		fmt.Printf("Winner: %s\n", gameState.Winner)
	}
}

type TexasPlayerPublicInfo struct {
	Name       string
	Money      int64
	HasFolded  bool
	CurrentPot int64
	Cards      []pokergo.Card
	Hand       *pokergo.Hand
	BestCards  []pokergo.Card
}

func (playerPublicInfo TexasPlayerPublicInfo) String() string {
	foldedString := "in game"
	if playerPublicInfo.HasFolded {
		foldedString = "has folded"
	}
	playerString := fmt.Sprintf("%s, pot: %d$, %s, total: %d$",
		playerPublicInfo.Name,
		playerPublicInfo.CurrentPot,
		foldedString,
		playerPublicInfo.Money)
	if playerPublicInfo.Hand != nil {
		playerString += fmt.Sprintf(", Cards: %s, Hand: %s, Combination: %s",
			playerPublicInfo.Cards,
			playerPublicInfo.Hand,
			playerPublicInfo.BestCards,
		)
	}
	return playerString
}
