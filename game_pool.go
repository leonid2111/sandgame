package main

import (
	"fmt"
	"strconv"
	"encoding/json"
)

type ServerMessage struct {
	Name      string     `json:"name"`
    Grid      [][]int    `json:"grid"`
    Scores    []string   `json:"scores"`
	Active    bool       `json:"active"`
    Comment   string     `json:"comment"`
}

type GamePool struct {
	counter    int
    register   chan *Player
    unregister chan *Player
    move       chan *PlayerMessage

    //players    map[*Player] bool
	first      *Player
	active     *Player
	
	update     chan bool
	grid       [][]int 
}

func NewGame(size int) *GamePool {
    return &GamePool{
		counter:    0,
        register:   make(chan *Player),
        unregister: make(chan *Player),
        move:       make(chan *PlayerMessage),
        //players:    make(map[*Player] bool),  
		first:      nil, 
		active:     nil, 
		update:     make(chan bool),
		grid:       initialize(size),
    }
}

func (pool *GamePool) get_players_scores() []string {
	var lines []string
	p := pool.first
	for {
		line := p.id + " : " + strconv.Itoa(p.score)
		fmt.Printf("p: %s  pn: %s  pp:%s \n", p.id, p.next.id, p.prev.id )
		
		lines = append(lines, line)
		if p.next == pool.first {
			break;
		} else {
			p = p.next
		}
	}
	return lines
}

func (pool *GamePool) update_all(msg ServerMessage) {
	updateJson, _ := json.Marshal(msg)
	p := pool.first
	for {
		fmt.Printf("player %s msg: %s\n", p.id, msg.Comment )
		p.Conn.WriteJSON(string(updateJson)) 
		if p.next == pool.first {
			break;
		} else {
			p = p.next
		}
	}	
}



func (pool *GamePool) Start() {
    for {
        select {
        case client := <-pool.register:
			pool.counter++
			fmt.Printf("player %d registering\n", pool.counter )
			client.id = "Player "+strconv.Itoa(pool.counter)
			
			if pool.active == nil {  // first client 
				client.next = client
				client.prev = client
				pool.active = client
				pool.first = client
			} else {
				client.next = pool.active
				client.prev = pool.active.prev
				pool.active.prev.next = client
				pool.active.prev = client
			}
			scoreboard := pool.get_players_scores()
			comm := client.id + " joined"
			msg := ServerMessage{
				Name:client.id,
				Grid:pool.grid,
				Scores:scoreboard,
				Active:(client==pool.active),
				Comment:comm}
			updateJson, _ := json.Marshal(msg)
			client.Conn.WriteJSON(string(updateJson))

			msg = ServerMessage{Scores:scoreboard, Comment:comm}
			pool.update_all(msg)
            break
			
		case client := <-pool.unregister:
            fmt.Printf("%s left\n", client.id)
			if client == client.next {
				pool.active = nil
				pool.first = nil
			} else {
				client.next.prev = client.prev
				client.prev.next = client.next
				if client == pool.active {
					pool.active = client.next
					// message to client.next here that it is now active
					// or message this to eveyone
				}
				if client == pool.first {
					pool.first = client.next
				}
				scoreboard := pool.get_players_scores()
				comm := client.id + " left"
				msg := ServerMessage{Scores:scoreboard, Comment:comm}
				pool.update_all(msg)
			}
            break

		case move := <-pool.move:
			fmt.Printf("client move")
			var xy [2]int
			json.Unmarshal(move.Message, &xy)
			go add_sand(xy, pool.grid, pool.update)

			pool.active = pool.active.next
			// message to everyone that pool.active is now active
			break

		case <-pool.update:
			fmt.Printf("updating the grid\n")			
			msg := ServerMessage{ Grid:pool.grid}
			pool.update_all(msg)
			break			
        }
    }
}



