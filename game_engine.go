package main

import (
	"fmt"
	"time"
    "math/rand"
)

var Spread = [4][2]int{{0,1},{0,-1},{1,0},{-1,0}}
const DELAY = 1

func initialize(size int) [][]int {
	var grid = make([][]int, size)
	//source := rand.NewSource(time.Now().UnixNano())
	source := rand.NewSource(42)
    randm := rand.New(source)
	for i := range grid {
		grid[i] = make([]int, size)
		for j := range grid[i]{
			grid[i][j] = randm.Intn(4)
		}
	}	
    return grid
}

func add_sand(xy [2]int, grid [][]int, update chan<- bool, next chan<- bool) {
	time.Sleep(DELAY * time.Second)
	h := grid[xy[0]][xy[1]] + 1
	if h==4 {
		full := make([][2]int,1)
		full[0] = xy
		distribute(full, grid, update)
	} else {
		grid[xy[0]][xy[1]]++
		update<-true
	}
	next<-true
}

func distribute(full [][2]int, grid [][]int, update chan<- bool) {
	fmt.Printf("full: %+v\n", full)
	var new_full [][2]int
	for _, x := range full {
		grid[x[0]][x[1]] -= 4
		for _, s := range Spread {
			y := fall(x,s)
			if is_inside(y, len(grid)){
				grid[y[0]][y[1]]++
				if grid[y[0]][y[1]]==4 {
					new_full = append(new_full, y)
				}
			} else {
				//*score++
			}
		}
	}
	update<-true
	time.Sleep(DELAY * time.Second)
	if len(new_full)>0 {
		distribute(new_full, grid, update)
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
