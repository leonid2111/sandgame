package main

import (
	"fmt"
	"time"
    "math/rand"
)

var Spread = [4][2]int{{0,1},{0,-1},{1,0},{-1,0}}
const DELAY = 1


//source := rand.NewSource(time.Now().UnixNano())
var source = rand.NewSource(42)
var randm = rand.New(source)

func initialize(size int) [][]int {
	var grid = make([][]int, size)
	for i := range grid {
		grid[i] = make([]int, size)
		for j := range grid[i]{
			grid[i][j] = randm.Intn(4)
		}
	}	
    return grid
}

func add_sand(xy [2]int, grid [][]int, update chan<- int, next chan<- bool) {
	time.Sleep(DELAY * time.Second)
	grid[xy[0]][xy[1]]++
	if grid[xy[0]][xy[1]] > 3 {
		full := make([][2]int,1)
		full[0] = xy
		distribute(full, grid, update)
	} else {
		update<-0
	}
	next<-true
}

func distribute(full [][2]int, grid [][]int, update chan<- int) {

	m := len(full)
	i := 0
	if m > 1 {
		i = randm.Intn(m)
	}
	fmt.Printf("%d full cells, distributing cell %d, full: %+v\n", m, i, full)
	
	x := full[i]
	grid[x[0]][x[1]] -= 4
	switch {
	case i==0:
		full = full[1:]
	case i==m-1:
		full = append(full[:i], full[i+1:]...)
	default:
		full = full[:m-1]
	}	

	score := 0
	for _, s := range Spread {
		y := fall(x,s)
		//fmt.Printf(" adding to %+v\n", y)
		if is_inside(y, len(grid)){
			//fmt.Printf("  %+v is inside \n", y)
			grid[y[0]][y[1]] += 1
			//fmt.Printf(" grid : %+v\n", grid)
			if grid[y[0]][y[1]] > 3 {
				full = append(full, y)
			}
		} else {
			score++
		}
	}

	//fmt.Printf(" full after spread: %+v\n", full)
	//fmt.Printf(" grid after spread: %+v\n", grid)

	update<-score
	time.Sleep(DELAY * time.Second)
	if len(full) > 0 {
		distribute(full, grid, update)
	}	
} 


func fall(x [2]int, s [2]int) [2]int {
	var y = [2]int{x[0]+s[0], x[1]+s[1]}
	return y
}

func is_inside(x [2]int, size int) bool {
	if (x[0]<0 || x[1]<0 || x[0]>=size || x[1]>=size) {
		return false } else {
		return true }
}
