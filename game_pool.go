package main

import (
	"fmt"
	"encoding/json"
)

const MAX_PLAYERS = 10

type GamePool struct {
    Register   chan *Client
    Unregister chan *Client
    Move       chan *ClientMessage
    players    map[*Client] int
	active     []bool
	queue      chan int
	update     chan bool
	next       chan bool
	scores     []int
	grid       [][]int 
}

func NewGame(size int) *GamePool {
    return &GamePool{
        Register:   make(chan *Client),
        Unregister: make(chan *Client),
        Move:       make(chan *ClientMessage),
        players:    make(map[*Client] int),  
		active:     make([]bool, 0), 
		queue:      make(chan int, MAX_PLAYERS),
		update:       make(chan bool),
		next:       make(chan bool),
		scores:      make([]int, 0),
		grid:        initialize(size),
    }
}

func (pool *GamePool) Start() {
    for {
        select {
        case client := <-pool.Register:
			n := len(pool.players)
			pool.players[client] = n
			pool.active = append(pool.active, true)
			pool.scores = append(pool.scores, 0)
			if n==MAX_PLAYERS {
				fmt.Println("Max number of players reached")
			}
			pool.queue <- n
			fmt.Printf("Player %d joined\n", n)
			updateJson, _ := json.Marshal(pool.grid)
			client.Conn.WriteJSON(string(updateJson))
            break
			
		case client := <-pool.Unregister:
			n := pool.players[client]
			pool.active[n] = false
            fmt.Printf("Player %d left\n", n)
            break

		case move := <-pool.Move:
			n := pool.players[move.Player]
			var xy [2]int
			json.Unmarshal(move.Message, &xy)
			go add_sand(n, xy, pool.grid, &pool.scores[n], pool.update, pool.next)
			break

		case <-pool.update:
			fmt.Printf("updating the grid\n")
			updateJson, _ := json.Marshal(pool.grid)
			for client, _ := range pool.players {
				client.Conn.WriteJSON(string(updateJson))
			}
			break
			
		case <-pool.next:
			for n := range pool.queue {
				if pool.active[n] {
					fmt.Printf("next player: %d\n", n)
					pool.queue <- n
					break
				}
			}
			break
			
		
        }
    }
}



