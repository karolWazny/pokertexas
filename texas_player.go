package pokertexas

import (
	"fmt"
	"github.com/karolWazny/pokergo"
	"strconv"
)

type TexasPlayer struct {
	player          *Player
	hand            []pokergo.Card
	bestCombination []pokergo.Card
	bestHand        *pokergo.Hand
	hasFolded       bool
	hasPlayed       bool
	currentPot      int64
}

func (texasPlayer TexasPlayer) GetPublicInfo() TexasPlayerPublicInfo {
	playerInfo := TexasPlayerPublicInfo{
		Name:       texasPlayer.player.name,
		Money:      texasPlayer.player.money,
		HasFolded:  texasPlayer.hasFolded,
		CurrentPot: texasPlayer.currentPot,
	}
	if !texasPlayer.hasFolded && texasPlayer.bestHand != nil {
		playerInfo.Cards = texasPlayer.hand
		playerInfo.BestCards = texasPlayer.bestCombination
		playerInfo.Hand = texasPlayer.bestHand
	}
	return playerInfo
}

func (texasPlayer TexasPlayer) String() string {
	return texasPlayer.player.String() + " " + fmt.Sprintf("%v", texasPlayer.hand) + " " + strconv.FormatInt(texasPlayer.currentPot, 10)
}
