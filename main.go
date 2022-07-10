package main

import (
	"fmt"
	"log"
	"flag"
	"time"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}

type Player struct {
    id   string
    Conn *websocket.Conn
    gamePool *GamePool
	next *Player
	prev *Player
	score int
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


var DELAY time.Duration

func main() {
	var port = flag.String("p", "8080", "game port")
	var size = flag.Int("s", 12, "game grid size")
	var delay = flag.Int("d", 1, "delay")
	flag.Parse()
	fmt.Printf("Starting sandgame on port %s with grid size %d\n", *port, *size)
	DELAY = time.Duration(*delay);
	pool := NewGame(*size)
    go pool.Start()
	//pool.next <- true   // to remove?
	
    http.HandleFunc("/sandgame", func(w http.ResponseWriter, r *http.Request) {
        connectWs(pool, w, r)
    })	
    log.Fatal(http.ListenAndServe(":"+*port, nil))
}
