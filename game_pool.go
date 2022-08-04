package main

import (
	"fmt"
	"strconv"
)


type ServerMessage struct {
	Header      string   `json:"header"`
    Grid      [][]int    `json:"grid"`
    Scores    []string   `json:"scores"`
	Activate    bool       `json:"activate"`
    Comment   string     `json:"comment"`
}

type GamePool struct {
	counter    int
    register   chan *Player
    unregister chan *Player
    move       chan [2]int

	first      *Player
	active     *Player
	
	update     chan int
	next       chan bool
	grid       [][]int 
}

func NewGame(size int, saturation float64, rseed uint64) *GamePool {
    return &GamePool{
		counter:    0,
        register:   make(chan *Player),
        unregister: make(chan *Player),
        move:       make(chan [2]int),
		first:      nil, 
		active:     nil, 
		update:     make(chan int),
		next:       make(chan bool),
		grid:       initialize(size, saturation, rseed),
    }
}

func (pool *GamePool) get_players_scores() []string {
	var lines []string
	for p:= pool.first;; p = p.next {
		line := p.id + " : " + strconv.Itoa(p.score)
		lines = append(lines, line)
		if p.next == pool.first {
			break;
		}
	}
	return lines
}

func (pool *GamePool) update_all(activate bool, comm string) {
	msg := ServerMessage{
		Header:pool.active.id+" - your move",
		Grid:pool.grid,
		Scores:pool.get_players_scores(),
		Activate:activate,
		Comment:comm,
	}
	//fmt.Printf("Updating active player, %s\n", pool.active.id )
	pool.active.update(msg)
	
	//fmt.Printf("Updating waiting players\n")
	for p:= pool.active.next; p != pool.active; p = p.next {
		msg.Header = p.id+" - waiting for "+pool.active.id
		msg.Activate = false
		p.update(msg)
	}
}



func (pool *GamePool) Start() {
    for {
        select {
        case client := <-pool.register:
			pool.counter++
			fmt.Printf("player %d registering\n", pool.counter )
			client.id += strconv.Itoa(pool.counter)
			var activate bool
			if pool.active == nil {  // first client 
				client.next = client
				client.prev = client
				pool.active = client
				pool.first = client
				activate = true
			} else {                // put new client last in the queue
				client.next = pool.active
				client.prev = pool.active.prev
				pool.active.prev.next = client
				pool.active.prev = client
			}

			comm := client.id + " joined"
			pool.update_all(activate, comm)
            break
			
		case client := <-pool.unregister:
            fmt.Printf("%s left\n", client.id)
			if client == client.next {
				pool.active = nil
				pool.first = nil
				fmt.Printf("No more players, waiting for someone to join.\n")
			} else {
				client.next.prev = client.prev
				client.prev.next = client.next
				
				var activate bool
				if client == pool.active {
					pool.active = client.next
					activate = true
				}
				if client == pool.first {
					pool.first = client.next
				}
				comm := client.id + " left"
				pool.update_all(activate, comm)
			}
            break

		case move := <-pool.move:
			//fmt.Printf("client moves\n")
			fmt.Printf("%s moves:  %+v \n", pool.active.id, move)
			go add_sand(move, pool.grid, pool.update, pool.next)
			break

		case n := <-pool.update:
			if pool.active != nil { // make sure last player didn't leave while the grid is still updating
				//fmt.Printf("updating the grid, adding %d to %s\n", n, pool.active.id)
				pool.active.score += n
				pool.update_all(false, pool.active.id+" adding sand")
			} 
			break			
        		
		case <-pool.next:
			if pool.active != nil { // make sure last player didn't leave while the grid was updating
				pool.active = pool.active.next
				pool.update_all(true, pool.active.id+" is next")
			}
			break			
		}
	}
}
