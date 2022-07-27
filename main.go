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
var size *int

//Usage:  ./sandgame -t 100 -n 12 -r 42 -s 0.75

func main() {
	var port = flag.String("p", "8080", "game port")
	size = flag.Int("n", 12, "game grid size")
	var delay = flag.Int("t", 200, "dealy time between updates, millis")
	var saturation = flag.Float64("s", 0.8, "initial grid saturation, in [0,1]")
	var rseed = flag.Uint64("r", 42, "seed for rand source, taken from timer if < 0")
	flag.Parse()
	
	DELAY  = time.Duration(*delay) * time.Millisecond
	fmt.Printf("Starting sandgame on port %s with grid size %d, dt=%+v, saturation=%.2f, rseed=%d\n",
		*port, *size, DELAY, *saturation, *rseed)
	
	pool := NewGame(*size, *saturation, *rseed)
    go pool.Start()
	
    http.HandleFunc("/sandgame", func(w http.ResponseWriter, r *http.Request) {
        connectWs(pool, w, r)
    })	
    log.Fatal(http.ListenAndServe(":"+*port, nil))
}

