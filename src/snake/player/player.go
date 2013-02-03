package player

import (
	"time"
)

type Player struct {
	FromClient        chan *Message    `json:"-"`
	ToClient          chan interface{} `json:"-"`
	Heading           string           `json:"-"`
	HeadingChanges    []string         `json:"-"`
	Position          [][2]int         `json:"position"`
	LostGame          bool             `json:"hasLost"`
	JustAte           bool             `json:"-"`
	Disconnected      bool             `json:"-"`
	PingSent          time.Time        `json:"-"`
}

type Message struct {
	Heading string `json:"heading"`
	Ping    string `json:"ping"`
}

var (
	reversals = map[string]string{
		"up":    "down",
		"down":  "up",
		"right": "left",
		"left":  "right",
	}
)

func (player *Player) UpdateHeading(update *Message) {
	pending := len(player.HeadingChanges)
	if pending == 2 {
		return
	}
	var nextHeading string
	if pending > 0 {
		nextHeading = player.HeadingChanges[0]
	} else {
		nextHeading = player.Heading
	}
	reversal, validHeading := reversals[nextHeading]
	h := update.Heading
	if !validHeading || h == nextHeading || h == reversal {
		return
	}
	player.HeadingChanges = append(player.HeadingChanges, h);
}

func (player *Player) CollidedInto(otherPlayer *Player) bool {
	for _, pos := range otherPlayer.Position {
		if player.Position[0] == pos {
			return true
		}
	}
	return false
}

func (player *Player) HitSelf() bool {
	for i, pos := range player.Position {
		if player.Position[0] == pos && i > 0 {
			return true
		}
	}
	return false
}

func (player *Player) ExceededBounds(width, height int) bool {
	frontPosition := player.Position[0]
	x, y := frontPosition[0], frontPosition[1]
	return x < 0 || x >= width || y < 0 || y >= height
}

func (player *Player) AdvancePosition() {
	frontPosition := player.Position[0]
	x, y := frontPosition[0], frontPosition[1]
	var newFrontPosition [2]int
	pending := len(player.HeadingChanges)
	if pending > 0 {
		player.Heading = player.HeadingChanges[0]
		player.HeadingChanges = player.HeadingChanges[1:pending]
	}
	switch player.Heading {
	case "left":
		newFrontPosition = [2]int{x - 1, y}
	case "right":
		newFrontPosition = [2]int{x + 1, y}
	case "up":
		newFrontPosition = [2]int{x, y - 1}
	case "down":
		newFrontPosition = [2]int{x, y + 1}
	}
	nextPosition := [][2]int{newFrontPosition}
	lastPosition := len(player.Position) - 1
	if player.JustAte {
		lastPosition += 1
		player.JustAte = false
	}
	player.Position = append(nextPosition, player.Position[:lastPosition]...)
}
