package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	part1sum := 0
	hasMoreLines := true
	for hasMoreLines {
		maxX := 0
		maxY := 0
		grid := map[int]map[int]bool{}
		scanner.Scan()
		for len(strings.TrimSpace(scanner.Text())) > 0 {
			grid[maxY] = make(map[int]bool)
			line := scanner.Text()
			maxX = len(line)
			for x, c := range line {
				if c == '#' {
					grid[maxY][x] = true
				}
			}
			maxY++
			hasMoreLines = scanner.Scan()
		}
		for x := 0; x < maxX-1; x++ {
			reflection := 0
			distance := 0
			for reflection < 2 && x-distance >= 0 && x+distance+1 < maxX {
				for y := 0; y < maxY; y++ {
					if grid[y][x-distance] != grid[y][x+distance+1] {
						//fmt.Printf("failed  (%d,%d) | (%d,%d)\n", y, x-distance, y, x+distance+1)
						reflection++
					} else {
						//fmt.Printf("success (%d,%d) | (%d,%d)\n", y, x-distance, y, x+distance+1)
					}
				}
				distance++
			}
			if reflection == 1 {
				part1sum += x + 1
			}
		}

		for y := 0; y < maxY-1; y++ {
			reflection := 0
			distance := 0
			for reflection < 2 && y-distance >= 0 && y+distance+1 < maxY {
				for x := 0; x < maxX; x++ {
					if grid[y-distance][x] != grid[y+distance+1][x] {
						//fmt.Printf("failed (%d,%d) | (%d,%d)\n", y-distance, x, y+distance+1, x)
						reflection++
					} else {
						//fmt.Printf("succeeded (%d,%d) | (%d,%d)\n", y-distance, x, y+distance+1, x)
					}
				}
				distance++
			}
			if reflection == 1 {
				part1sum += 100 * (y + 1)
			}
		}
	}
	fmt.Println(part1sum)
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}
