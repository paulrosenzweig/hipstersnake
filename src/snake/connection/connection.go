package connection

import (
	"code.google.com/p/go.net/websocket"
	"fmt"
	"snake/game"
	"snake/player"
	"time"
)

var (
	ConnectionHandler = websocket.Handler(socketHandler)
)

func socketHandler(ws *websocket.Conn) {
	p := &player.Player{
		FromClient:     make(chan *player.Message, 0),
		ToClient:       make(chan interface{}, 0),
		HeadingChanges: make([]string, 0),
		Disconnected:   false,
	}

	quitPinger := make(chan int, 0)
	defer func() { quitPinger <- 0 }()
	go pinger(p, quitPinger)

	go sender(ws, p)

	quit := receiver(ws, p)

	myName, theirName := game.Pair(p)

	p.ToClient <- map[string]string{"myName": myName, "theirName": theirName}

	<-quit
}

func receiver(ws *websocket.Conn, p *player.Player) chan int {
	quit := make(chan int, 0)
	go func() {
		for {
			m := &player.Message{}
			err := websocket.JSON.Receive(ws, m)
			if err != nil {
				break
			}
			if m.Ping != "" {
				fmt.Printf(
					"%s %s %s\n",
					time.Now().UTC(),
					time.Since(p.PingSent),
					ws.Request().RemoteAddr,
				)
			} else {
				p.FromClient <- m
			}
		}
		p.Disconnected = true
		quit <- 1
	}()
	return quit
}

func sender(ws *websocket.Conn, p *player.Player) {
	var err error
	err = nil
	for m := range p.ToClient {
		if err == nil {
			err = websocket.JSON.Send(ws, m)
		}
	}
	ws.Close()
}

func pinger(p *player.Player, quit chan int) {
	p.PingSent = time.Now()
	p.ToClient <- map[string]string{"ping": "ping"}
	ticker := time.Tick(30e9)
	for {
		select {
		case t := <-ticker:
			p.PingSent = t
			p.ToClient <- map[string]string{"ping": "ping"}
		case <-quit:
			return
		}
	}
}
