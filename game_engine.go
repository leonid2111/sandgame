package main

import (
	"fmt"
	"time"
    //"math/rand" - does not have binomial sampler
	"golang.org/x/exp/rand"
	"gonum.org/v1/gonum/stat/distuv"
)

var Spread = [4][2]int{{0,1},{0,-1},{1,0},{-1,0}}

var randm *rand.Rand

func initialize(size int, saturation float64, rseed uint64) [][]int {
	if rseed < 0 {
		rseed = uint64(time.Now().UnixNano())
	}
	source := rand.NewSource(rseed)	
	randm = rand.New(source)

	
	binom := distuv.Binomial{N:3, P:saturation, Src:source}	
	var grid = make([][]int, size)
	for i := range grid {
		grid[i] = make([]int, size)
		for j := range grid[i]{
			grid[i][j] = int(binom.Rand())
		}
	}
    return grid
}


func add_sand(xy [2]int, grid [][]int, update chan<- int, next chan<- bool) {
	time.Sleep(100 * time.Millisecond)
	grid[xy[0]][xy[1]]++
	if grid[xy[0]][xy[1]] > 3 {
		full := make( map[[2]int] bool)
		full[xy] = true
		distribute(full, grid, update)
	} else {
		update<-0
	}
	next<-true
}


func distribute(full map[[2]int]bool, grid [][]int, update chan<- int) {
	// Randomly select a full cell
	i := randm.Intn(len(full))
	j := 0
	var x [2]int
	for c := range full {
		x = c
		j++
		if i==j { break }
	}
		
	fmt.Printf("Init weight = %d\n", total_sand(grid))
	print_fulls(full, grid)
	fmt.Printf("Distributing cell %d,  xy = %+v\n", i, x)

	grid[x[0]][x[1]] -= 4
	if grid[x[0]][x[1]] < 4 {
		delete(full,x)
	}
	
	score := 0
	for _, s := range Spread {
		y := [2]int{x[0]+s[0], x[1]+s[1]}
		if (y[0]>=0 && y[1]>=0 && y[0] < *size && y[1] < *size) {	
			grid[y[0]][y[1]] += 1
			if grid[y[0]][y[1]] > 3 {
				full[y] = true
			}
		} else {
			score++
		}
	}

	fmt.Printf("weight = %d  score=%d\n", total_sand(grid), score)
	print_fulls(full, grid)

	update<-score
	time.Sleep(DELAY)
	if len(full) > 0 {
		distribute(full, grid, update)
	}	
} 



func total_sand(grid [][]int) int {
	total := 0
	for i := range grid {
		for j := range grid[i]{
			total += grid[i][j]
		}
	}
	return total
}


func print_fulls(full map[[2]int]bool, grid [][]int){
	for x := range full {
		fmt.Printf(" [%d %d)  %d]  ", x[0], x[1], grid[x[0]][x[1]])
	}
	fmt.Printf("\n")
}
