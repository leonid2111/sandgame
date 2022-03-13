package main

import "fmt"

type Game struct {
    Register   chan *Client
    Unregister chan *Client
    Clients    map[*Client]bool
    Broadcast  chan Message
	//GameEngine [][]int
}

func NewGame() *Game {
    return &Game{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Clients:    make(map[*Client]bool),
        Broadcast:  make(chan Message),
		//GameEngine:  nil,
    }
}

func (pool *Game) Start() {
    for {
        select {
        case client := <-pool.Register:
            pool.Clients[client] = true
            fmt.Println("Connected players: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                //fmt.Println(client)
                client.Conn.WriteJSON(Message{Type: 1, Body: "New Player Joined..."})
            }
            break
        case client := <-pool.Unregister:
            delete(pool.Clients, client)
            fmt.Println("Players left: ", len(pool.Clients))
            for client, _ := range pool.Clients {
                client.Conn.WriteJSON(Message{Type: 1, Body: "Player Disconnected..."})
            }
            break
			
        case message := <-pool.Broadcast:
            fmt.Println("Sending message to all clients in Pool")
            for client, _ := range pool.Clients {
				fmt.Print("trying to write  %+v\n", message)
                if err := client.Conn.WriteJSON(message); err != nil {
                    fmt.Println(err)
                    return
                }
            }

			
        }
    }
}
