package game

import (
	"math/rand"
	"snake/player"
	"time"
)

var (
	pairingChannel = make(chan *player.Player, 0)
)

type Game struct {
	Width     int            `json:"width"`
	Height    int            `json:"height"`
	PlayerOne *player.Player `json:"playerOne"`
	PlayerTwo *player.Player `json:"playerTwo"`
	HasEnded  bool           `json:"hasEnded"`
	Food      [][2]int       `json:"food"`
}

func Pair(player *player.Player) (myName, theirName string) {
	select {
	case pairingChannel <- player:
		myName = "playerTwo"
		theirName = "playerOne"
	case otherPlayer := <-pairingChannel:
		go create(player, otherPlayer)
		myName = "playerOne"
		theirName = "playerTwo"
	}
	return myName, theirName
}

func create(PlayerOne *player.Player, PlayerTwo *player.Player) {
	if PlayerOne.Disconnected || PlayerTwo.Disconnected {
		if !PlayerTwo.Disconnected {
			Pair(PlayerTwo)
		}
		if !PlayerOne.Disconnected {
			Pair(PlayerOne)
		}
		return
	}
	PlayerOne.Position = [][2]int{
		[2]int{3, 7},
		[2]int{3, 6},
		[2]int{3, 5},
		[2]int{3, 4},
		[2]int{3, 3},
	}
	PlayerOne.Heading = "down"
	PlayerTwo.Position = [][2]int{
		[2]int{46, 42},
		[2]int{46, 43},
		[2]int{46, 44},
		[2]int{46, 45},
		[2]int{46, 46},
	}
	PlayerTwo.Heading = "up"
	game := &Game{
		Width:     50,
		Height:    50,
		PlayerOne: PlayerOne,
		PlayerTwo: PlayerTwo,
		Food:      [][2]int{},
	}
	game.run()
}

func (game *Game) run() {
	timeInterval := 2e8
	moveTicker := time.Tick(time.Duration(timeInterval))
	foodTicker := time.Tick(1e9)
	for {
		select {
		case <-moveTicker:
			game.PlayerOne.AdvancePosition()
			game.PlayerTwo.AdvancePosition()
			game.checkForLoser()
			game.PlayerOne.ToClient <- game
			game.PlayerTwo.ToClient <- game
			if game.HasEnded {
				close(game.PlayerOne.ToClient)
				close(game.PlayerTwo.ToClient)
				return
			} else {
				game.eatFood()
				if game.PlayerOne.JustAte {
					timeInterval *= 0.97
				}
				if game.PlayerTwo.JustAte {
					timeInterval *= 0.97
				}
				if game.PlayerOne.JustAte || game.PlayerTwo.JustAte {
					moveTicker = time.Tick(time.Duration(timeInterval))
				}
			}
		case <-foodTicker:
			x := rand.Int() % game.Width
			y := rand.Int() % game.Height
			game.Food = append(game.Food, [2]int{x, y})
		case update := <-game.PlayerOne.FromClient:
			game.PlayerOne.UpdateHeading(update)
		case update := <-game.PlayerTwo.FromClient:
			game.PlayerTwo.UpdateHeading(update)
		}
	}
}

func (game *Game) eatFood() {
	remainingFood := [][2]int{}
	for _, location := range game.Food {
		if game.PlayerTwo.Position[0] == location {
			game.PlayerTwo.JustAte = true
		} else if game.PlayerOne.Position[0] == location {
			game.PlayerOne.JustAte = true
		} else {
			remainingFood = append(remainingFood, location)
		}
	}
	game.Food = remainingFood
}

func (game *Game) checkForLoser() {
	game.PlayerOne.LostGame = game.PlayerOne.ExceededBounds(game.Width, game.Height) || game.PlayerOne.CollidedInto(game.PlayerTwo) || game.PlayerOne.HitSelf()
	game.PlayerTwo.LostGame = game.PlayerTwo.ExceededBounds(game.Width, game.Height) || game.PlayerTwo.CollidedInto(game.PlayerOne) || game.PlayerTwo.HitSelf()
	game.HasEnded = game.PlayerTwo.LostGame || game.PlayerOne.LostGame
}
