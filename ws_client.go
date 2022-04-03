package main

import (
    "fmt"
    "log"
    "github.com/gorilla/websocket"
)

type Client struct {
    ID   string
    Conn *websocket.Conn
    GamePool *GamePool
}

type ClientMessage struct {
    Player    *Client    //`json:"type"`
    Message   []byte      //`json:"body"`
}

func (c *Client) Listen() {
    defer func() {
        c.GamePool.Unregister <- c
        c.Conn.Close()
    }()

    for {
        _, p, err := c.Conn.ReadMessage()
        if err != nil {
            log.Println(err)
            return
        }
        msg := ClientMessage{Player: c, Message: p}//Message{Type: messageType, Body: string(p)}
        c.GamePool.Move <- &msg
        fmt.Printf("Message Received: %+v \n", string(p))
    }
}
