package main

import (
    "fmt"
    "log"

    "github.com/gorilla/websocket"
)

type Client struct {
    ID   string
    Conn *websocket.Conn
    Game *Game
}

type Message struct {
    Type int    `json:"type"`
    Body string `json:"body"`
}

func (c *Client) Listen() {
    defer func() {
        c.Game.Unregister <- c
        c.Conn.Close()
    }()

    for {
        messageType, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        message := Message{Type: messageType, Body: string(p)}
        c.Game.Broadcast <- message
        fmt.Printf("Message Received: %+v\n", message)
    }
}
