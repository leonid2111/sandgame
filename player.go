package main

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"encoding/json"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

	
type Player struct { 
    id   string
    gamePool *GamePool
	next *Player
	prev *Player
	score int
    Conn *websocket.Conn
	sim chan bool
}


func (player *Player) update(msg ServerMessage) {
	if player.Conn != nil {                          // real player
		updateJson, _ := json.Marshal(msg)
		player.Conn.WriteJSON(string(updateJson))
	} else if msg.Activate {                         // sim player
		player.sim <- true
	}
}



func (p *Player) simulate() {
	defer func() {
        p.gamePool.unregister <- p
    }()
	
	time.Sleep(10000 * time.Millisecond)
	fmt.Println("New sim registeres")
    p.gamePool.register <- p

	for{
		select {
		case <- p.sim:
			// do something else here
			xm := randm.Intn(*size)
			ym := randm.Intn(*size)
			xy := [2]int{xm, ym}   

			/*for k = 0; k < *size/2; k++ {
				for j=0; 
			}*/

			
			p.gamePool.move <- xy			
		}
	}
}


func (player *Player) Listen() {
    defer func() {
        player.gamePool.unregister <- player
        player.Conn.Close()
    }()

    for {
        _, p, err := player.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
		
        //fmt.Printf("%s moves:  %+v \n", player.id, string(p))
		var xy [2]int
		json.Unmarshal(p, &xy)
        player.gamePool.move <- xy
    }
}



func connectWs(pool *GamePool, w http.ResponseWriter, r *http.Request) {
    fmt.Println("New player hits wsendpoint")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Fprintf(w, "%+v\n", err)
    }

	client := &Player{
        Conn: conn,
        gamePool: pool,
		id: "Player ",
    }	
    pool.register <- client
    client.Listen()
}



