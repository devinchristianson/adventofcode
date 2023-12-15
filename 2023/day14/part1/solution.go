package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type Point struct {
	y int
	x int
}

func printGrid(grid [][]rune) {
	for y := range grid {
		for x := range grid[y] {
			fmt.Print(string(grid[y][x]))
		}
		fmt.Println()
	}
	fmt.Println()
}

func moveRocks(d Point, rocks []*Point, grid [][]rune) {
	for _, rock := range rocks {
		move := &Point{x: 0, y: 0}
		for canMove := true; canMove; {
			if rock.y+move.y+d.y < 0 || rock.x+move.x+d.x < 0 || rock.y+move.y+d.y >= len(grid) || rock.x+move.x+d.x >= len(grid[0]) {
				canMove = false
			} else {
				if grid[rock.y+move.y+d.y][rock.x+move.x+d.x] == '.' {
					move.x += d.x
					move.y += d.y
				} else {
					canMove = false
				}
			}
		}
		grid[rock.y][rock.x] = '.'
		grid[rock.y+move.y][rock.x+move.x] = 'O'
		rock.y += move.y
		rock.x += move.x
	}
}

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part1sum := 0
	grid := [][]rune{}
	rocks := []*Point{}
	for y := 0; scanner.Scan(); y++ {
		grid = append(grid, make([]rune, len(scanner.Text())))
		for x, r := range scanner.Text() {
			grid[y][x] = r
			if r == 'O' {
				rocks = append(rocks, &Point{y: y, x: x})
			}
		}
	}
	printGrid(grid)
	moveRocks(Point{x: 0, y: -1}, rocks, grid) //north
	printGrid(grid)
	for _, rock := range rocks {
		part1sum += len(grid) - rock.y
	}
	fmt.Println(part1sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
