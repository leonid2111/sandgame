package main

import (
	"fmt"
	"log"
	"flag"
	//"strconv"
	"time"
	"net/http"
)




var DELAY time.Duration
var size *int

//Usage:  ./sandgame -t 100 -n 12 -r 42 -s 0.75 -a 2

func main() {
	var port = flag.String("p", "8080", "game port")
	size = flag.Int("n", 12, "game grid size")
	var delay = flag.Int("t", 200, "dealy time between updates, millis")
	var saturation = flag.Float64("s", 0.8, "initial grid saturation, in [0,1]")
	var rseed = flag.Uint64("r", 42, "seed for rand source, taken from timer if < 0")
	var autoplayers = flag.Int("a", 0, "number of auto-players")
	flag.Parse()
	
	DELAY  = time.Duration(*delay) * time.Millisecond
	fmt.Printf("Starting sandgame on port %s with grid size %d, dt=%+v, saturation=%.2f, rseed=%d\n",
		*port, *size, DELAY, *saturation, *rseed)
	fmt.Printf("Number of simulated players: =%d\n", *autoplayers)
	
	pool := NewGame(*size, *saturation, *rseed)
    go pool.Start()
	
    http.HandleFunc("/sandgame", func(w http.ResponseWriter, r *http.Request) {
        connectWs(pool, w, r)
    })
	
	for m := 0; m < *autoplayers; m++ {
		p :=  &Player{gamePool: pool}	
		go p.simulate()
	}
	
    log.Fatal(http.ListenAndServe(":"+*port, nil))
}

