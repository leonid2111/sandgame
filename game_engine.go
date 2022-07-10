package main

import (
	"fmt"
	//"log"
	"time"
    "math/rand"
)

var Spread = [4][2]int{{0,1},{0,-1},{1,0},{-1,0}}


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
		//full := make([][2]int,1)
		//full[0] = xy
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
		/*		switch {
		case i==0:
			full = full[1:]
		case i<m-1:
			full = append(full[:i], full[i+1:]...)
		case i==m-1:
			full = full[:m-1]
		default:
			log.Fatalf("This cannot happen: i=%d m=%d\n", i, m)
		}	*/
	}
	
	score := 0
	for _, s := range Spread {
		y := fall(x,s)
		//fmt.Printf(" adding to %+v\n", y)
		if is_inside(y, len(grid)){
			grid[y[0]][y[1]] += 1
			if grid[y[0]][y[1]] > 3 {
				//full = append(full, y)
				full[y] = true
			}
		} else {
			score++
		}
	}

	//fmt.Printf(" full after spread: %+v\n", full)
	//print_fulls(full, grid)
	//fmt.Printf(" grid after spread: %+v\n", grid)
	fmt.Printf("weight = %d  score=%d\n", total_sand(grid), score)
	print_fulls(full, grid)
	//fmt.Printf("Distributing cell %d,  xy = %+v\n", i, x)

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
