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
}


func (player *Player) update(msg ServerMessage) {
	fmt.Println("update client player")
	if player.Conn != nil {
		updateJson, _ := json.Marshal(msg)
		player.Conn.WriteJSON(string(updateJson))
	} else if msg.Activate {
		xy := [2]int{3,3}
		go add_sand(xy, player.gamePool.grid, player.gamePool.update, player.gamePool.next)
	}
}



func (player *Player) simulate() {
	defer func() {
        player.gamePool.unregister <- player
    }()
	
	time.Sleep(10000 * time.Millisecond)
	fmt.Println("New sim registeres")
    player.gamePool.register <- player

	for{}

}





type PlayerMessage struct {
    Player    *Player    //`json:"type"`
    Message   []byte     //`json:"body"`
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
        msg := PlayerMessage{Player: player, Message: p}
        fmt.Printf("%s moves:  %+v \n", player.id, string(p))
        player.gamePool.move <- &msg
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
    }	
    pool.register <- client
    client.Listen()
}



