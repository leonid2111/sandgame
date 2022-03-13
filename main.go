package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
    ReadBufferSize:  1024,
    WriteBufferSize: 1024,
}


func connectWs(pool *Game, w http.ResponseWriter, r *http.Request) {
    fmt.Println("New player hits wsendpoint")
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        fmt.Fprintf(w, "%+v\n", err)
    }
    client := &Client{
        Conn: conn,
        Game: pool,
    }	
    pool.Register <- client
    client.Listen()
}



func setupRoutes() {
	pool := NewGame()
    go pool.Start()
	
    http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
        connectWs(pool, w, r)
    })	
}

func main() {
	fmt.Println("test ws starts")
    setupRoutes()
    log.Fatal(http.ListenAndServe(":8080", nil))
}
