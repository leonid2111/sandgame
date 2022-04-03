package main

import (
	"fmt"
	"log"
	"flag"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func connectWs(pool *GamePool, w http.ResponseWriter, r *http.Request) {
    fmt.Println("New player hits wsendpoint")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Fprintf(w, "%+v\n", err)
    }
    client := &Client{
        Conn: conn,
        GamePool: pool,
    }	
    pool.Register <- client
    client.Listen()
}

func main() {
	var port = flag.String("p", "8080", "game port")
	var size = flag.Int("s", 12, "game grid size")
	flag.Parse()
	fmt.Printf("Starting sandgame on port %s with grid size %d\n", *port, *size)

	pool := NewGame(*size)
    go pool.Start()
	pool.next <- true
	
    http.HandleFunc("/sandgame", func(w http.ResponseWriter, r *http.Request) {
        connectWs(pool, w, r)
    })	
    log.Fatal(http.ListenAndServe(":"+*port, nil))
}
