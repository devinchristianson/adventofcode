package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

type Point struct {
	y int
	x int
}

type direction int

const (
	north direction = iota
	west
	south
	east
)

func printGrid(grid [][]rune) {
	for y := range grid {
		for x := range grid[y] {
			fmt.Print(string(grid[y][x]))
		}
		fmt.Println()
	}
	fmt.Println()
}

func makeSequence(r rune, count int) []rune {
	a := make([]rune, 0, count)
	for i := 0; i < count; i++ {
		a = append(a, r)
	}
	return a
}

func moveRocksEastWest(grid [][]rune, d direction) [][]rune {
	newGrid := make([][]rune, 0, len(grid))
	for y := range grid {
		newGrid = append(newGrid, make([]rune, 0, len(grid[y])))
		rockCount := 0
		spaceCount := 0
		for x := range grid[y] {
			switch grid[y][x] {
			case 'O':
				rockCount++
			case '#':
				if d == east {
					newGrid[y] = append(newGrid[y], makeSequence('.', spaceCount)...)
					newGrid[y] = append(newGrid[y], makeSequence('O', rockCount)...)
				} else {
					newGrid[y] = append(newGrid[y], makeSequence('O', rockCount)...)
					newGrid[y] = append(newGrid[y], makeSequence('.', spaceCount)...)
				}
				newGrid[y] = append(newGrid[y], '#')
				rockCount = 0
				spaceCount = 0
			case '.':
				spaceCount++
			}
		}
		if d == east {
			newGrid[y] = append(newGrid[y], makeSequence('.', spaceCount)...)
			newGrid[y] = append(newGrid[y], makeSequence('O', rockCount)...)
		} else {
			newGrid[y] = append(newGrid[y], makeSequence('O', rockCount)...)
			newGrid[y] = append(newGrid[y], makeSequence('.', spaceCount)...)
		}
	}
	return newGrid
}

func moveRocksNorthSouth(grid [][]rune, d direction) [][]rune {
	for x := range grid[0] {
		buffer := make([]rune, 0, len(grid[0]))
		rockCount := 0
		spaceCount := 0
		for y := range grid {
			switch grid[y][x] {
			case 'O':
				rockCount++
			case '#':
				if d == south {
					buffer = append(buffer, makeSequence('.', spaceCount)...)
					buffer = append(buffer, makeSequence('O', rockCount)...)
				} else {
					buffer = append(buffer, makeSequence('O', rockCount)...)
					buffer = append(buffer, makeSequence('.', spaceCount)...)
				}
				buffer = append(buffer, '#')
				rockCount = 0
				spaceCount = 0
			case '.':
				spaceCount++
			}
		}
		if d == south {
			buffer = append(buffer, makeSequence('.', spaceCount)...)
			buffer = append(buffer, makeSequence('O', rockCount)...)
		} else {
			buffer = append(buffer, makeSequence('O', rockCount)...)
			buffer = append(buffer, makeSequence('.', spaceCount)...)
		}
		for y := range grid {
			grid[y][x] = buffer[y]
		}
	}
	return grid
}

func spin(grid [][]rune) [][]rune {
	grid = moveRocksNorthSouth(grid, north) // north
	grid = moveRocksEastWest(grid, west)    // west
	grid = moveRocksNorthSouth(grid, south) // south
	grid = moveRocksEastWest(grid, east)    // east
	return grid
}

func findPoints(grid [][]rune, count int) []Point {
	points := make([]Point, 0, count)
	for y := range grid {
		for x := range grid[y] {
			if grid[y][x] == 'O' {
				points = append(points, Point{y: y, x: x})
			}
		}
	}
	return points
}

func calculateLoad(p []Point, length int) int {
	sum := 0
	for _, r := range p {
		sum += length - r.y
	}
	return sum
}

func contains(s [][]int, c []int) []int {
	var indexes []int
	for i := range s {
		if slices.Equal(s[i], c) {
			indexes = append(indexes, i)
		}
	}
	if len(indexes) == 0 {
		return []int{-1}
	} else {
		return indexes
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	grid := [][]rune{}
	numRocks := 0
	for y := 0; scanner.Scan(); y++ {
		grid = append(grid, make([]rune, len(scanner.Text())))
		for x, r := range scanner.Text() {
			grid[y][x] = r
			if r == 'O' {
				numRocks++
			}
		}
	}
	grid = moveRocksNorthSouth(grid, north)
	var states [][]Point
	total_cycles := 1_000_000_000
	period := 0
	transient := 0
	num_cycles := 0
	for {
		grid = spin(grid)
		points := findPoints(grid, numRocks)
		//fmt.Println(calculateLoad(points, len(grid)))
		if index := slices.IndexFunc(states, func(p []Point) bool {
			for i, r := range p {
				if points[i] != r {
					return false
				}
			}
			return true
		}); index != -1 {
			transient = index
			period = num_cycles - transient
			break
		}
		states = append(states, points)
		num_cycles++
	}
	final_state := states[(total_cycles-transient)%period+transient-1]
	sum := calculateLoad(final_state, len(grid))

	fmt.Println(sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
